package iamdao

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/morehao/goark/apps/iam/internal/tenantctx"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	tenantCreateCallbackName = "iam:tenant:create"
	tenantQueryCallbackName  = "iam:tenant:query"
	tenantUpdateCallbackName = "iam:tenant:update"
	tenantDeleteCallbackName = "iam:tenant:delete"
)

type tenantScopeKind string

const (
	tenantScopeNone    tenantScopeKind = ""
	tenantScopeCompany tenantScopeKind = "company"
	tenantScopeTenant  tenantScopeKind = "tenant"
)

var (
	registerTenantCallbackOnce sync.Once

	tableScopeMap = map[string]tenantScopeKind{
		"iam_company":         tenantScopeTenant,
		"iam_tenant_config":   tenantScopeTenant,
		"iam_user":            tenantScopeCompany,
		"iam_role":            tenantScopeCompany,
		"iam_menu":            tenantScopeCompany,
		"iam_department":      tenantScopeCompany,
		"iam_user_role":       tenantScopeCompany,
		"iam_user_department": tenantScopeCompany,
		"iam_role_menu":       tenantScopeCompany,
	}
)

func RegisterTenantCallbacks() error {
	var callbackErr error
	registerTenantCallbackOnce.Do(func() {
		db := dbclient.IamDB(context.Background())
		if db == nil {
			callbackErr = fmt.Errorf("iam db is nil")
			return
		}

		if err := db.Callback().Query().Before("gorm:query").Register(tenantQueryCallbackName, tenantScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Create().Before("gorm:create").Register(tenantCreateCallbackName, tenantCreateScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Update().Before("gorm:update").Register(tenantUpdateCallbackName, tenantUpdateScopeCallback); err != nil {
			callbackErr = err
			return
		}
		if err := db.Callback().Delete().Before("gorm:delete").Register(tenantDeleteCallbackName, tenantScopeCallback); err != nil {
			callbackErr = err
			return
		}
	})
	return callbackErr
}

func tenantCreateScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := tenantctx.FromStdContext(tx.Statement.Context)

	mainTableName, mainQualifier := resolveMainTable(tx, false, clause.From{})
	if mainTableName == "" {
		return
	}
	mainScopeKind := resolveMainScopeKind(tx.Statement.Schema, mainTableName)
	if mainScopeKind == tenantScopeNone {
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

func tenantUpdateScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := tenantctx.FromStdContext(tx.Statement.Context)
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

	tenantScopeCallback(tx)
}

func tenantScopeCallback(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Unscoped {
		return
	}

	scope, _ := tenantctx.FromStdContext(tx.Statement.Context)
	if scope.IsPlatformAdmin() {
		return
	}

	fromClause, hasFrom := getFromClause(tx)
	if hasFrom {
		if err := applyJoinTenantScope(tx, scope, &fromClause); err != nil {
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
	if mainScopeKind == tenantScopeNone {
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

func applyJoinTenantScope(tx *gorm.DB, scope tenantctx.Scope, from *clause.From) error {
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
		if scopeKind == tenantScopeNone {
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

func resolveMainScopeKind(schemaRef *schema.Schema, tableName string) tenantScopeKind {
	if schemaRef != nil {
		if _, ok := schemaRef.FieldsByDBName["company_id"]; ok {
			return tenantScopeCompany
		}
		if _, ok := schemaRef.FieldsByDBName["tenant_id"]; ok {
			return tenantScopeTenant
		}
	}
	return resolveScopeByTableName(tableName)
}

func resolveScopeByTableName(tableName string) tenantScopeKind {
	return tableScopeMap[normalizeTableName(tableName)]
}

func appendScopeExpression(tx *gorm.DB, qualifier string, scopeKind tenantScopeKind, scope tenantctx.Scope) error {
	expr, err := buildScopeExpression(qualifier, scopeKind, scope)
	if err != nil {
		return err
	}
	tx.Statement.AddClause(clause.Where{Exprs: []clause.Expression{expr}})
	return nil
}

func buildScopeExpression(qualifier string, scopeKind tenantScopeKind, scope tenantctx.Scope) (clause.Expression, error) {
	switch scopeKind {
	case tenantScopeCompany:
		if scope.CompanyID == 0 {
			return nil, code.GetError(code.TenantContextMissingError)
		}
		return clause.Eq{Column: clause.Column{Table: qualifier, Name: "company_id"}, Value: scope.CompanyID}, nil
	case tenantScopeTenant:
		if scope.TenantID == 0 {
			return nil, code.GetError(code.TenantContextMissingError)
		}
		return clause.Eq{Column: clause.Column{Table: qualifier, Name: "tenant_id"}, Value: scope.TenantID}, nil
	default:
		return nil, nil
	}
}

func fillCreateScopeValue(tx *gorm.DB, qualifier string, scopeKind tenantScopeKind, scope tenantctx.Scope) error {
	switch scopeKind {
	case tenantScopeCompany:
		if scope.CompanyID == 0 {
			return code.GetError(code.TenantContextMissingError)
		}
		return setScopeFieldValue(tx, qualifier, "company_id", scope.CompanyID)
	case tenantScopeTenant:
		if scope.TenantID == 0 {
			return code.GetError(code.TenantContextMissingError)
		}
		return setScopeFieldValue(tx, qualifier, "tenant_id", scope.TenantID)
	default:
		return nil
	}
}

func validatePlatformAdminCreateScope(tx *gorm.DB, qualifier string, scopeKind tenantScopeKind, scope tenantctx.Scope) error {
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
	case tenantScopeCompany:
		scopeValue = scope.CompanyID
	case tenantScopeTenant:
		scopeValue = scope.TenantID
	}
	if scopeValue == 0 {
		return code.GetError(code.TenantContextMissingError)
	}
	return setScopeFieldValue(tx, qualifier, fieldName, scopeValue)
}

func validateScopeMutation(dest any, scopeKind tenantScopeKind) error {
	fieldName := scopeFieldName(scopeKind)
	if fieldName == "" {
		return nil
	}
	if hasScopeFieldInMap(dest, fieldName) {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	return nil
}

func scopeFieldName(scopeKind tenantScopeKind) string {
	switch scopeKind {
	case tenantScopeCompany:
		return "company_id"
	case tenantScopeTenant:
		return "tenant_id"
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
