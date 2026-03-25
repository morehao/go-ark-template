package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/tenantctx"
	"github.com/morehao/goark/pkg/code"
)

func RequireCompanyID(ctx *gin.Context) (uint, error) {
	companyID := tenantctx.GetCompanyID(ctx)
	if companyID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return companyID, nil
}

func CheckCompanyAccess(ctx *gin.Context, targetCompanyID uint) error {
	if targetCompanyID == 0 {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return nil
	}
	companyID, err := RequireCompanyID(ctx)
	if err != nil {
		return err
	}
	if companyID != targetCompanyID {
		return code.GetError(code.TenantScopeForbiddenError)
	}
	return nil
}

func RequireTargetCompanyID(ctx *gin.Context, requestCompanyID uint) (uint, error) {
	companyID, err := NormalizeTargetCompanyID(ctx, requestCompanyID)
	if err != nil {
		return 0, err
	}
	if companyID == 0 {
		return 0, code.GetError(code.TenantContextMissingError)
	}
	return companyID, nil
}

func NormalizeTargetCompanyID(ctx *gin.Context, requestCompanyID uint) (uint, error) {
	if requestCompanyID > 0 {
		if err := CheckCompanyAccess(ctx, requestCompanyID); err != nil {
			return 0, err
		}
		return requestCompanyID, nil
	}
	if tenantctx.IsPlatformAdmin(ctx) {
		return 0, nil
	}
	companyID, err := RequireCompanyID(ctx)
	if err != nil {
		return 0, err
	}
	return companyID, nil
}

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
	return RequireTenantID(ctx)
}
