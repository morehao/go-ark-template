package ctroidc

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/apps/iam/internal/service/svcoidc"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type logoutCtr struct {
	logoutSvc svcoidc.LogoutSvc
}

func NewLogoutCtr() *logoutCtr {
	return &logoutCtr{
		logoutSvc: svcoidc.NewLogoutSvc(),
	}
}

// Logout
// @Tags OIDC
// @Summary 单点登出
// @accept application/json
// @Produce application/json
// @Param req body dtooidc.LogoutReq true "登出请求"
// @Success 200 {object} gincontext.DtoRender
// @Router /v1/iam/oidc/logout [POST]
func (ctr *logoutCtr) Logout(ctx *gin.Context) {
	var req dtooidc.LogoutReq
	if err := ctx.ShouldBind(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	sessionID := ctx.GetHeader("X-Sso-Session-Id")

	if err := ctr.logoutSvc.Logout(ctx, sessionID, req.RefreshToken); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	gincontext.Success(ctx, nil)
}
