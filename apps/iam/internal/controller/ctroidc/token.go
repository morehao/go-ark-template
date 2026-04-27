package ctroidc

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/apps/iam/internal/service/svcoidc"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type tokenCtr struct {
	tokenSvc svcoidc.TokenSvc
}

func NewTokenCtr() *tokenCtr {
	return &tokenCtr{
		tokenSvc: svcoidc.NewTokenSvc(),
	}
}

// Token
// @Tags OIDC
// @Summary 获取Token
// @accept application/json
// @Produce application/json
// @Param req formData dtooidc.TokenReq true "Token请求"
// @Success 200 {object} gincontext.DtoRender{data=dtooidc.TokenResp}
// @Router /v1/iam/oidc/token [POST]
func (ctr *tokenCtr) Token(ctx *gin.Context) {
	var req dtooidc.TokenReq
	if err := ctx.ShouldBind(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	resp, err := ctr.tokenSvc.ExchangeCodeForToken(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	gincontext.Success(ctx, resp)
}

// RefreshToken
// @Tags OIDC
// @Summary 刷新Token
// @accept application/json
// @Produce application/json
// @Param req formData dtooidc.TokenRefreshReq true "刷新Token请求"
// @Success 200 {object} gincontext.DtoRender{data=dtooidc.TokenResp}
// @Router /v1/iam/oidc/refreshToken [POST]
func (ctr *tokenCtr) RefreshToken(ctx *gin.Context) {
	var req dtooidc.TokenRefreshReq
	if err := ctx.ShouldBind(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	resp, err := ctr.tokenSvc.RefreshAccessToken(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	gincontext.Success(ctx, resp)
}
