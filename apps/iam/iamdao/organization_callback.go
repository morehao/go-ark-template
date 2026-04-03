package iamdao

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/pkg/contextkeys"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	organizationCreateCallbackName = "iam:organization:create"
	organizationQueryCallbackName  = "iam:organization:query"
	organizationUpdateCallbackName = "iam:organization:update"
	organizationDeleteCallbackName = "iam:organization:delete"
)

type organizationScopeKind string

const (
	organizationScopeNone   organizationScopeKind = ""
	organizationScopeTenant organizationScopeKind = "tenant"
	organizationScopeOrg    organizationScopeKind = "organization"
)

type Scope struct {
	OrganizationID uint
	TenantID       uint
	UserType       string
}

func GetScope(ctx context.Context) (Scope, bool) {
	if ctx == nil {
		return Scope{}, false
	}
	var scope Scope
	ok := false

	if v := ctx.Value(contextkeys.KeyTenantID); v != nil {
		if id, found := readUintValue(v); found {
			scope.TenantID = id
			ok = true
		}
	}
	if v := ctx.Value(contextkeys.KeyOrgID); v != nil {
		if id, found := readUintValue(v); found {
			scope.OrganizationID = id
			ok = true
		}
	}
	if v := ctx.Value(contextkeys.KeyUserType); v != nil {
		if ut, found := readStringValue(v); found {
			scope.UserType = ut
		}
	}
	return scope, ok
}

func (s Scope) IsPlatformAdmin() bool {
	return s.UserType == "platform_admin"
}

func readUintValue(v any) (uint, bool) {
	switch value := v.(type) {
	case uint:
		return value, true
	case uint64:
		return uint(value), true
	case uint32:
		return uint(value), true
	case int:
		if value < 0 {
			return 0, false
		}
		return uint(value), true
	case int64:
		if value < 0 {
			return 0, false
		}
		return uint(value), true
	case string:
		if value == "" {
			return 0, false
		}
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(parsed), true
	default:
		return 0, false
	}
}

func readStringValue(v any) (string, bool) {
	value, ok := v.(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

var (
	registerOrganizationCallbackOnce sync.Once

	tableScopeMap = map[string]organizationScopeKind{
		"iam_tenant":              organizationScopeOrg,
		"iam_organization_config": organizationScopeOrg,
		"iam_user":                organizationScopeTenant,
		"iam_role":                organizationScopeTenant,
		"iam_menu":                organizationScopeTenant,
		"iam_department":          organizationScopeTenant,
		"iam_user_role":           organizationScopeTenant,
		"iam_user_department":     organizationScopeTenant,
		"iam_role_menu":           organizationScopeTenant,
	}
)

func RegisterOrganizationCallbacks() error {
	var callbackErr error
	registerOrganizationCallbackOnce.Do(func() {
		db := dbclient.IamDB(context.Background())
		if db == nil {
			callbackErr = fmt.Errorf("iam db is nil")
			return
		}

		if err := db.Callback().Query().Before("gorm:query").Register(organizationQueryCallbackName, organizationScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Create().Before("gorm:create").Register(organizationCreateCallbackName, organizationCreateScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Update().Before("gorm:update").Register(organizationUpdateCallbackName, organizationUpdateScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Delete().Before("gorm:delete").Register(organizationDeleteCallbackName, organizationScopeCallback); err != nil {
			callbackErr = err
			return
		}
	})
	return callbackErr
}

func organizationCreateScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := GetScope(tx.Statement.Context)

	mainTableName, mainQualifier := resolveMainTable(tx, false, clause.From{})
	if mainTableName == "" {
		return
	}
	mainScopeKind := resolveMainScopeKind(tx.Statement.Schema, mainTableName)
	if mainScopeKind == organizationScopeNone {
		return
	}

	if scope.IsPlatformAdmin() {
		if err := validatePlatformAdminCreateScope(tx, mainQualifier, mainScopeKind, scope); err != nil {
			tx.AddError(err)
		}
		return
	}

	if err := fillCreateScopeValue(tx, mainQualifier, mainScopeKind, scope); err != nil {
		tx.AddError(err)
	}
}

func organizationUpdateScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := GetScope(tx.Statement.Context)
	if !scope.IsPlatformAdmin() {
		mainTableName, _ := resolveMainTable(tx, false, clause.From{})
		if mainTableName != "" {
			mainScopeKind := resolveMainScopeKind(tx.Statement.Schema, mainTableName)
			if err := validateScopeMutation(tx.Statement.Dest, mainScopeKind); err != nil {
				tx.AddError(err)
				return
			}
		}
	}

	organizationScopeCallback(tx)
}

func organizationScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := GetScope(tx.Statement.Context)
	if scope.IsPlatformAdmin() {
		return
	}

	fromClause, hasFrom := getFromClause(tx)
	if hasFrom {
		if err := applyJoinOrganizationScope(tx, scope, &fromClause); err != nil {
			tx.AddError(err)
			return
		}
		tx.Statement.AddClause(fromClause)
	}

	mainTableName, mainQualifier := resolveMainTable(tx, hasFrom, fromClause)
	if mainTableName == "" {
		return
	}

	mainScopeKind := resolveMainScopeKind(tx.Statement.Schema, mainTableName)
	if mainScopeKind == organizationScopeNone {
		return
	}

	if err := appendScopeExpression(tx, mainQualifier, mainScopeKind, scope); err != nil {
		tx.AddError(err)
	}
}

func getFromClause(tx *gorm.DB) (clause.From, bool) {
	clauseValue, ok := tx.Statement.Clauses["FROM"]
	if !ok {
		return clause.From{}, false
	}
	fromExpr, ok := clauseValue.Expression.(clause.From)
	if !ok {
		return clause.From{}, false
	}
	return fromExpr, true
}

func applyJoinOrganizationScope(tx *gorm.DB, scope Scope, from *clause.From) error {
	if from == nil || len(from.Joins) == 0 {
		return nil
	}

	for idx := range from.Joins {
		joinItem := from.Joins[idx]
		if joinItem.Expression != nil {
			return code.GetError(code.TenantJoinUnsafeError)
		}

		tableName := normalizeTableName(joinItem.Table.Name)
		scopeKind := resolveScopeByTableName(tableName)
		if scopeKind == organizationScopeNone {
			continue
		}

		qualifier := tableQualifier(joinItem.Table)
		if qualifier == "" {
			return code.GetError(code.TenantJoinUnsafeError)
		}

		expr, err := buildScopeExpression(qualifier, scopeKind, scope)
		if err != nil {
			return err
		}
		joinItem.ON.Exprs = append(joinItem.ON.Exprs, expr)
		from.Joins[idx] = joinItem
	}

	return nil
}

func resolveMainTable(tx *gorm.DB, hasFrom bool, from clause.From) (tableName string, qualifier string) {
	if hasFrom && len(from.Tables) > 0 {
		table := from.Tables[0]
		name := normalizeTableName(table.Name)
		if name != "" {
			return name, tableQualifier(table)
		}
	}

	name := normalizeTableName(tx.Statement.Table)
	if name == "" && tx.Statement.Schema != nil {
		name = normalizeTableName(tx.Statement.Schema.Table)
	}
	return name, name
}

func resolveMainScopeKind(schemaRef *schema.Schema, tableName string) organizationScopeKind {
	if schemaRef != nil {
		if _, ok := schemaRef.FieldsByDBName["tenant_id"]; ok {
			return organizationScopeTenant
		}
		if _, ok := schemaRef.FieldsByDBName["organization_id"]; ok {
			return organizationScopeOrg
		}
	}
	return resolveScopeByTableName(tableName)
}

func resolveScopeByTableName(tableName string) organizationScopeKind {
	return tableScopeMap[normalizeTableName(tableName)]
}

func appendScopeExpression(tx *gorm.DB, qualifier string, scopeKind organizationScopeKind, scope Scope) error {
	expr, err := buildScopeExpression(qualifier, scopeKind, scope)
	if err != nil {
		return err
	}
	tx.Statement.AddClause(clause.Where{Exprs: []clause.Expression{expr}})
	return nil
}

func buildScopeExpression(qualifier string, scopeKind organizationScopeKind, scope Scope) (clause.Expression, error) {
	switch scopeKind {
	case organizationScopeTenant:
		if scope.TenantID == 0 {
			return nil, code.GetError(code.TenantContextMissingError)
		}
		return clause.Eq{Column: clause.Column{Table: qualifier, Name: "tenant_id"}, Value: scope.TenantID}, nil
	case organizationScopeOrg:
		if scope.OrganizationID == 0 {
			return nil, code.GetError(code.TenantContextMissingError)
		}
		return clause.Eq{Column: clause.Column{Table: qualifier, Name: "organization_id"}, Value: scope.OrganizationID}, nil
	default:
		return nil, nil
	}
}

func fillCreateScopeValue(tx *gorm.DB, qualifier string, scopeKind organizationScopeKind, scope Scope) error {
	switch scopeKind {
	case organizationScopeTenant:
		if scope.TenantID == 0 {
			return code.GetError(code.TenantContextMissingError)
		}
		return setScopeFieldValue(tx, qualifier, "tenant_id", scope.TenantID)
	case organizationScopeOrg:
		if scope.OrganizationID == 0 {
			return code.GetError(code.TenantContextMissingError)
		}
		return setScopeFieldValue(tx, qualifier, "organization_id", scope.OrganizationID)
	default:
		return nil
	}
}

func validatePlatformAdminCreateScope(tx *gorm.DB, qualifier string, scopeKind organizationScopeKind, scope Scope) error {
	fieldName := scopeFieldName(scopeKind)
	if fieldName == "" {
		return nil
	}
	v, ok := readScopeValueFromDest(tx.Statement.Dest, fieldName)
	if ok && v > 0 {
		return nil
	}

	scopeValue := uint(0)
	switch scopeKind {
	case organizationScopeTenant:
		scopeValue = scope.TenantID
	case organizationScopeOrg:
		scopeValue = scope.OrganizationID
	}
	if scopeValue == 0 {
		return code.GetError(code.TenantContextMissingError)
	}
	return setScopeFieldValue(tx, qualifier, fieldName, scopeValue)
}

func validateScopeMutation(dest any, scopeKind organizationScopeKind) error {
	fieldName := scopeFieldName(scopeKind)
	if fieldName == "" {
		return nil
	}
	if hasScopeFieldInMap(dest, fieldName) {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	return nil
}

func scopeFieldName(scopeKind organizationScopeKind) string {
	switch scopeKind {
	case organizationScopeTenant:
		return "tenant_id"
	case organizationScopeOrg:
		return "organization_id"
	default:
		return ""
	}
}

func setScopeFieldValue(tx *gorm.DB, qualifier, fieldName string, value uint) error {
	if tx == nil || tx.Statement == nil {
		return nil
	}

	tx.Statement.SetColumn(fieldName, value)
	if v, ok := readScopeValueFromDest(tx.Statement.Dest, fieldName); ok && v == value {
		return nil
	}
	tx.Statement.SetColumn(toCamelField(fieldName), value)
	if v, ok := readScopeValueFromDest(tx.Statement.Dest, fieldName); ok && v == value {
		return nil
	}

	if hasScopeFieldInMap(tx.Statement.Dest, fieldName) {
		setScopeFieldInMap(tx.Statement.Dest, fieldName, value)
		return nil
	}

	column := clause.Column{Table: qualifier, Name: fieldName}
	tx.Statement.AddClause(clause.Set{{Column: column, Value: value}})
	return nil
}

func hasScopeFieldInMap(dest any, fieldName string) bool {
	v := reflect.ValueOf(dest)
	if !v.IsValid() {
		return false
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return false
	}
	for _, key := range v.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		if normalizeColumnName(key.String()) == fieldName {
			return true
		}
	}
	return false
}

func setScopeFieldInMap(dest any, fieldName string, value uint) {
	v := reflect.ValueOf(dest)
	if !v.IsValid() {
		return
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return
	}
	for _, key := range v.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		if normalizeColumnName(key.String()) != fieldName {
			continue
		}
		v.SetMapIndex(key, reflect.ValueOf(value).Convert(v.Type().Elem()))
		return
	}
	v.SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(value).Convert(v.Type().Elem()))
}

func readScopeValueFromDest(dest any, fieldName string) (uint, bool) {
	v := reflect.ValueOf(dest)
	if !v.IsValid() {
		return 0, false
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, false
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			if key.Kind() != reflect.String {
				continue
			}
			if normalizeColumnName(key.String()) != fieldName {
				continue
			}
			return castUintValue(v.MapIndex(key))
		}
		return 0, false
	}

	if v.Kind() == reflect.Struct {
		field := v.FieldByName(toCamelField(fieldName))
		if !field.IsValid() {
			return 0, false
		}
		return castUintValue(field)
	}

	return 0, false
}

func castUintValue(v reflect.Value) (uint, bool) {
	if !v.IsValid() {
		return 0, false
	}
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, false
		}
		return castUintValue(v.Elem())
	}

	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(v.Uint()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := v.Int()
		if iv < 0 {
			return 0, false
		}
		return uint(iv), true
	}
	return 0, false
}

func toCamelField(column string) string {
	parts := strings.Split(column, "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

func normalizeColumnName(name string) string {
	trimmed := strings.TrimSpace(name)
	trimmed = strings.Trim(trimmed, "`")
	trimmed = strings.ToLower(trimmed)
	if idx := strings.LastIndex(trimmed, "."); idx >= 0 {
		trimmed = trimmed[idx+1:]
	}
	return trimmed
}

func tableQualifier(table clause.Table) string {
	if table.Alias != "" {
		return table.Alias
	}
	return normalizeTableName(table.Name)
}

func normalizeTableName(name string) string {
	trimmed := strings.TrimSpace(name)
	trimmed = strings.Trim(trimmed, "`")
	if idx := strings.Index(trimmed, " "); idx > 0 {
		trimmed = trimmed[:idx]
	}
	return trimmed
}
