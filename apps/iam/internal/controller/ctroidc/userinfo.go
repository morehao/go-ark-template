package ctroidc

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/internal/service/svcoidc"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type userinfoCtr struct {
	tokenSvc    svcoidc.TokenSvc
	userInfoSvc svcoidc.UserInfoSvc
}

func NewUserinfoCtr() *userinfoCtr {
	return &userinfoCtr{
		tokenSvc:    svcoidc.NewTokenSvc(),
		userInfoSvc: svcoidc.NewUserInfoSvc(),
	}
}

// UserInfo
// @Tags OIDC
// @Summary 获取用户信息
// @Security BearerAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=dtooidc.UserInfoResp}
// @Router /v1/iam/oidc/userinfo [GET]
func (ctr *userinfoCtr) UserInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		gincontext.Fail(ctx, code.GetError(code.AuthTokenRequiredError))
		return
	}

	tokenStr := authHeader[7:]

	tokenEntity, err := ctr.tokenSvc.ValidateAccessToken(ctx, tokenStr)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	resp, err := ctr.userInfoSvc.GetUserInfo(ctx, tokenEntity.PersonID)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	gincontext.Success(ctx, resp)
}
