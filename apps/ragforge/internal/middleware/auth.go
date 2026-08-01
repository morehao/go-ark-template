package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/config"
	"github.com/morehao/golib/biz/gcontext"
)

type authClaims struct {
	UserID   uint   `json:"userID"`
	TenantID uint   `json:"tenantID"`
	Role     string `json:"role"`
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := resolveClaims(ctx)
		if claims != nil {
			ctx.Set(gcontext.KeyUserID, claims.UserID)
			ctx.Set(gcontext.KeyTenantID, claims.TenantID)
			if claims.Role != "" {
				ctx.Set(gcontext.KeyUserType, claims.Role)
			}
		}
		ctx.Next()
	}
}

func resolveClaims(ctx *gin.Context) *authClaims {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := parseDevToken(token); err == nil {
			return claims
		}
	}

	if userID := ctx.GetHeader("X-User-ID"); userID != "" {
		uid, _ := strconv.ParseUint(userID, 10, 64)
		tid := uint64(1)
		if tenantID := ctx.GetHeader("X-Tenant-ID"); tenantID != "" {
			tid, _ = strconv.ParseUint(tenantID, 10, 64)
		}
		return &authClaims{UserID: uint(uid), TenantID: uint(tid)}
	}

	if config.Conf.Server.Env != "prod" {
		return &authClaims{UserID: 1, TenantID: 1, Role: "admin"}
	}
	return nil
}

func parseDevToken(token string) (*authClaims, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var claims authClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}
