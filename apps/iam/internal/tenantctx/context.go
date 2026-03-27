package tenantctx

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

type stdCtxKey string

const (
	stdCtxOrganizationIDKey stdCtxKey = "iam_std_organization_id"
	stdCtxTenantIDKey       stdCtxKey = "iam_std_tenant_id"
	stdCtxUserTypeKey       stdCtxKey = "iam_std_user_type"
)

type Scope struct {
	OrganizationID uint
	TenantID       uint
	UserType       string
}

const (
	HeaderOrganizationID = "X-Organization-ID"
	HeaderTenantID       = "X-Tenant-ID"
	HeaderUserType       = "X-User-Type"

	CtxOrganizationID = "iam_organization_id"
	CtxTenantID       = "iam_tenant_id"
	CtxUserType       = "iam_user_type"

	UserTypePlatformAdmin = "platform_admin"
)

func SetFromHeader(ctx *gin.Context) {
	if organizationID := parseUintHeader(ctx.GetHeader(HeaderOrganizationID)); organizationID > 0 {
		ctx.Set(CtxOrganizationID, organizationID)
	}
	if tenantID := parseUintHeader(ctx.GetHeader(HeaderTenantID)); tenantID > 0 {
		ctx.Set(CtxTenantID, tenantID)
	}
	if userType := ctx.GetHeader(HeaderUserType); userType != "" {
		ctx.Set(CtxUserType, userType)
	}
}

func GetOrganizationID(ctx *gin.Context) uint {
	v, ok := ctx.Get(CtxOrganizationID)
	if !ok {
		return 0
	}
	if value, ok := v.(uint); ok {
		return value
	}
	return 0
}

func GetTenantID(ctx *gin.Context) uint {
	v, ok := ctx.Get(CtxTenantID)
	if !ok {
		return 0
	}
	if value, ok := v.(uint); ok {
		return value
	}
	return 0
}

func GetUserType(ctx *gin.Context) string {
	return ctx.GetString(CtxUserType)
}

func IsPlatformAdmin(ctx *gin.Context) bool {
	return GetUserType(ctx) == UserTypePlatformAdmin
}

func InjectToRequestContext(ctx *gin.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	stdCtx := context.Background()
	if ctx.Request != nil {
		stdCtx = ctx.Request.Context()
	}

	if organizationID := GetOrganizationID(ctx); organizationID > 0 {
		stdCtx = context.WithValue(stdCtx, stdCtxOrganizationIDKey, organizationID)
	}
	if tenantID := GetTenantID(ctx); tenantID > 0 {
		stdCtx = context.WithValue(stdCtx, stdCtxTenantIDKey, tenantID)
	}
	if userType := GetUserType(ctx); userType != "" {
		stdCtx = context.WithValue(stdCtx, stdCtxUserTypeKey, userType)
	}
	return stdCtx
}

func FromStdContext(ctx context.Context) (Scope, bool) {
	if ctx == nil {
		return Scope{}, false
	}

	var (
		scope Scope
		ok    bool
	)

	if v, found := readUintValue(ctx.Value(stdCtxOrganizationIDKey)); found {
		scope.OrganizationID = v
		ok = true
	} else if v, found := readUintValue(ctx.Value(CtxOrganizationID)); found {
		scope.OrganizationID = v
		ok = true
	}

	if v, found := readUintValue(ctx.Value(stdCtxTenantIDKey)); found {
		scope.TenantID = v
		ok = true
	} else if v, found := readUintValue(ctx.Value(CtxTenantID)); found {
		scope.TenantID = v
		ok = true
	}

	if v, found := readStringValue(ctx.Value(stdCtxUserTypeKey)); found {
		scope.UserType = v
		ok = true
	} else if v, found := readStringValue(ctx.Value(CtxUserType)); found {
		scope.UserType = v
		ok = true
	}

	return scope, ok
}

func (s Scope) IsPlatformAdmin() bool {
	return s.UserType == UserTypePlatformAdmin
}

func parseUintHeader(raw string) uint {
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0
	}
	return uint(v)
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
