package iammiddleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/internal/tenantctx"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gcontext"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/gauth/jwtauth"
)

func Auth() gin.HandlerFunc {
	initJWT()
	return func(ctx *gin.Context) {
		if jwtInited {
			tokenStr := extractToken(ctx)
			if tokenStr == "" {
				gincontext.Abort(ctx, code.GetError(gconstant.UnauthorizedErr))
				ctx.Abort()
				return
			}

			claims, err := auth.Parse(tokenStr)
			if err != nil {
				gincontext.Abort(ctx, code.GetError(gconstant.UnauthorizedErr))
				ctx.Abort()
				return
			}

			ctx.Set(gcontext.KeyUserID, claims.CustomData.UserID)
			ctx.Set(tenantctx.CtxUserType, claims.CustomData.Username)
		}

		if err := setOrganizationContextFromToken(ctx); err != nil {
			gincontext.Fail(ctx, err)
			ctx.Abort()
			return
		}
		if ctx.Request != nil {
			ctx.Request = ctx.Request.WithContext(tenantctx.InjectToRequestContext(ctx))
		}
		ctx.Next()
	}
}

type UserClaims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

var (
	auth      *jwtauth.Auth[UserClaims]
	jwtInited bool
)

func initJWT() {
	if config.Conf.JWT.SignKey == "" {
		return
	}
	var err error
	auth, err = jwtauth.New[UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return
	}
	jwtInited = true
}

func extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func setOrganizationContextFromToken(ctx *gin.Context) error {
	userID := gincontext.GetUserID(ctx)
	if userID == 0 {
		return code.GetError(code.TenantContextMissingError)
	}

	userEntity, err := iamdao.NewUserDao().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.TenantContextMissingError)
	}

	if userEntity.UserType != "" {
		ctx.Set(tenantctx.CtxUserType, userEntity.UserType)
	}

	if userEntity.TenantID > 0 {
		ctx.Set(tenantctx.CtxTenantID, userEntity.TenantID)
	}

	if userEntity.UserType == tenantctx.UserTypePlatformAdmin {
		return nil
	}

	if userEntity.TenantID == 0 {
		return code.GetError(code.TenantContextMissingError)
	}

	tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
	if err != nil {
		return err
	}
	if tenantEntity == nil || tenantEntity.ID == 0 || tenantEntity.OrganizationID == 0 {
		return code.GetError(code.TenantContextMissingError)
	}

	ctx.Set(tenantctx.CtxOrganizationID, tenantEntity.OrganizationID)
	return nil
}
