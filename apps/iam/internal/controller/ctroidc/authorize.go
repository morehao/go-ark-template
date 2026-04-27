package ctroidc

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/apps/iam/internal/service/svcoidc"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/goark/pkg/code"
)

type authorizeCtr struct {
	authorizeSvc svcoidc.AuthorizeSvc
	sessionSvc   svcoidc.SsoSessionSvc
}

func NewAuthorizeCtr() *authorizeCtr {
	return &authorizeCtr{
		authorizeSvc: svcoidc.NewAuthorizeSvc(),
		sessionSvc:   svcoidc.NewSsoSessionSvc(),
	}
}

// Authorize
// @Tags OIDC
// @Summary 授权入口
// @accept application/json
// @Produce application/json
// @Param req query dtooidc.AuthorizeReq true "授权请求"
// @Success 200 {object} gincontext.DtoRender
// @Router /v1/iam/oidc/authorize [GET]
func (ctr *authorizeCtr) Authorize(ctx *gin.Context) {
	var req dtooidc.AuthorizeReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	sessionID := ctx.GetHeader("X-Sso-Session-Id")

	if sessionID != "" {
		session, err := ctr.sessionSvc.GetValidSession(ctx, sessionID)
		if err != nil || session == nil || session.ID == 0 {
			sessionID = ""
		} else {
			app, err := ctr.authorizeSvc.ValidateClient(ctx, req.ClientID, req.RedirectURI)
			if err != nil {
				gincontext.Fail(ctx, err)
				return
			}

			resp, err := ctr.authorizeSvc.GenerateAuthCode(ctx, app, session.PersonID, 0, session.OrgID, &req)
			if err != nil {
				gincontext.Fail(ctx, err)
				return
			}

			redirectURI := req.RedirectURI + "?code=" + resp.Code + "&state=" + resp.State
			ctx.Redirect(302, redirectURI)
			return
		}
	}

	gincontext.Fail(ctx, code.GetError(code.AuthSessionRequiredError))
}
