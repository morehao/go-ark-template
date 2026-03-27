package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/tenantctx"
	"github.com/morehao/goark/pkg/code"
)

func RequireTenantID(ctx *gin.Context) (uint, error) {
	tenantID := tenantctx.GetTenantID(ctx)
	if tenantID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return tenantID, nil
}

func CheckTenantAccess(ctx *gin.Context, targetTenantID uint) error {
	if targetTenantID == 0 {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return nil
	}
	tenantID, err := RequireTenantID(ctx)
	if err != nil {
		return err
	}
	if tenantID != targetTenantID {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	return nil
}

func RequireTargetTenantID(ctx *gin.Context, requestTenantID uint) (uint, error) {
	tenantID, err := NormalizeTargetTenantID(ctx, requestTenantID)
	if err != nil {
		return 0, err
	}
	if tenantID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return tenantID, nil
}

func NormalizeTargetTenantID(ctx *gin.Context, requestTenantID uint) (uint, error) {
	if requestTenantID > 0 {
		if err := CheckTenantAccess(ctx, requestTenantID); err != nil {
			return 0, err
		}
		return requestTenantID, nil
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return 0, nil
	}
	tenantID, err := RequireTenantID(ctx)
	if err != nil {
		return 0, err
	}
	return tenantID, nil
}

func RequireOrganizationID(ctx *gin.Context) (uint, error) {
	organizationID := tenantctx.GetOrganizationID(ctx)
	if organizationID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return organizationID, nil
}

func CheckOrganizationAccess(ctx *gin.Context, targetOrganizationID uint) error {
	if targetOrganizationID == 0 {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return nil
	}
	organizationID, err := RequireOrganizationID(ctx)
	if err != nil {
		return err
	}
	if organizationID != targetOrganizationID {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	return nil
}

func RequireTargetOrganizationID(ctx *gin.Context, requestOrganizationID uint) (uint, error) {
	organizationID, err := NormalizeTargetOrganizationID(ctx, requestOrganizationID)
	if err != nil {
		return 0, err
	}
	if organizationID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return organizationID, nil
}

func NormalizeTargetOrganizationID(ctx *gin.Context, requestOrganizationID uint) (uint, error) {
	if requestOrganizationID > 0 {
		if err := CheckOrganizationAccess(ctx, requestOrganizationID); err != nil {
			return 0, err
		}
		return requestOrganizationID, nil
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return 0, nil
	}
	return RequireOrganizationID(ctx)
}
