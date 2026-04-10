package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/service/svcauth"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/gerror"
)

func TokenBlacklistCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		if svcauth.IsTokenBlacklisted(c, token) {
			gincontext.Abort(c, &gerror.Error{
				Code: gconstant.TokenInvalidErr,
				Msg:  "token已失效",
			})
			return
		}
		c.Next()
	}
}
