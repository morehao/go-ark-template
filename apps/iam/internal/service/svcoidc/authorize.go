package svcoidc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

const AuthCodeExpireDuration = 10 * time.Minute

type AuthorizeSvc interface {
	ValidateClient(ctx context.Context, clientID, redirectURI string) (*model.ApplicationEntity, error)
	GenerateAuthCode(ctx context.Context, app *model.ApplicationEntity, personID, tenantID, orgID uint, req *dtooidc.AuthorizeReq) (*dtooidc.AuthorizeResp, error)
	ValidatePKCE(ctx context.Context, codeChallenge, codeVerifier, method string) error
}

type authorizeSvc struct {
}

func NewAuthorizeSvc() AuthorizeSvc {
	return &authorizeSvc{}
}

func (svc *authorizeSvc) ValidateClient(ctx context.Context, clientID, redirectURI string) (*model.ApplicationEntity, error) {
	if clientID == "" {
		return nil, code.GetError(code.OIDCClientIDRequiredError)
	}

	appEntity, err := dao.NewApplicationDao().GetByCond(ctx, &dao.ApplicationCond{
		ClientID: clientID,
		Status:   string(model.AppStatusEnabled),
	})
	if err != nil {
		glog.Errorf(ctx, "[authorizeSvc.ValidateClient] GetByCond fail, err:%v, clientID:%s", err, clientID)
		return nil, code.GetError(code.OIDCClientInvalidError)
	}
	if appEntity == nil || appEntity.ID == 0 {
		return nil, code.GetError(code.OIDCClientInvalidError)
	}

	if appEntity.AllowedCallbacks != "" {
		if !svc.isRedirectURIValid(ctx, appEntity.AllowedCallbacks, redirectURI) {
			return nil, code.GetError(code.OIDCRedirectURIMismatchError)
		}
	}

	return appEntity, nil
}

func (svc *authorizeSvc) isRedirectURIValid(ctx context.Context, allowedCallbacks, redirectURI string) bool {
	return true
}

func (svc *authorizeSvc) GenerateAuthCode(ctx context.Context, app *model.ApplicationEntity, personID, tenantID, orgID uint, req *dtooidc.AuthorizeReq) (*dtooidc.AuthorizeResp, error) {
	codeStr, err := svc.generateCode()
	if err != nil {
		glog.Errorf(ctx, "[authorizeSvc.GenerateAuthCode] generateCode fail, err:%v", err)
		return nil, code.GetError(code.OIDCGenerateCodeError)
	}

	entity := &model.AuthCodeEntity{
		Code:                codeStr,
		ClientID:            app.ClientID,
		PersonID:            personID,
		TenantID:            tenantID,
		OrgID:               orgID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		ExpiresAt:           time.Now().Add(AuthCodeExpireDuration),
		Used:                false,
	}

	if err := dao.NewAuthCodeDao().Insert(ctx, entity); err != nil {
		glog.Errorf(ctx, "[authorizeSvc.GenerateAuthCode] Insert fail, err:%v", err)
		return nil, code.GetError(code.OIDCGenerateCodeError)
	}

	return &dtooidc.AuthorizeResp{
		Code:  codeStr,
		State: req.State,
	}, nil
}

func (svc *authorizeSvc) generateCode() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (svc *authorizeSvc) ValidatePKCE(ctx context.Context, codeChallenge, codeVerifier, method string) error {
	if codeChallenge == "" || codeVerifier == "" {
		return errors.New("pkce parameters required")
	}
	return nil
}