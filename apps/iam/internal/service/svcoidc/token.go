package svcoidc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/gauth/jwtauth"
	"github.com/morehao/golib/glog"
)

const (
	AccessTokenExpireDuration  = 1 * time.Hour
	RefreshTokenExpireDuration = 7 * 24 * time.Hour
	TokenIssuer                = "iam"
)

type TokenSvc interface {
	ExchangeCodeForToken(ctx context.Context, req *dtooidc.TokenReq) (*dtooidc.TokenResp, error)
	RefreshAccessToken(ctx context.Context, req *dtooidc.TokenRefreshReq) (*dtooidc.TokenResp, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*model.TokenEntity, error)
	RevokeToken(ctx context.Context, token string, tokenTypeHint string) error
}

type tokenSvc struct {
}

func NewTokenSvc() TokenSvc {
	return &tokenSvc{}
}

func (svc *tokenSvc) ExchangeCodeForToken(ctx context.Context, req *dtooidc.TokenReq) (*dtooidc.TokenResp, error) {
	if req.GrantType != "authorization_code" {
		return nil, code.GetError(code.OIDCUnsupportedGrantTypeError)
	}

	authCode, err := dao.NewAuthCodeDao().GetByCode(ctx, req.Code)
	if err != nil {
		glog.Errorf(ctx, "[tokenSvc.ExchangeCodeForToken] GetByCode fail, err:%v", err)
		return nil, code.GetError(code.OIDCInvalidCodeError)
	}
	if authCode == nil || authCode.ID == 0 {
		return nil, code.GetError(code.OIDCInvalidCodeError)
	}
	if authCode.Used {
		return nil, code.GetError(code.OIDCInvalidCodeError)
	}
	if time.Now().After(authCode.ExpiresAt) {
		return nil, code.GetError(code.OIDCInvalidCodeError)
	}
	if authCode.RedirectURI != req.RedirectURI {
		return nil, code.GetError(code.OIDCRedirectURIMismatchError)
	}

	appEntity, err := dao.NewApplicationDao().GetByCond(ctx, &dao.ApplicationCond{
		ClientID: authCode.ClientID,
		Status:   string(model.AppStatusEnabled),
	})
	if err != nil || appEntity == nil || appEntity.ID == 0 {
		return nil, code.GetError(code.OIDCClientInvalidError)
	}

	if appEntity.ClientSecret != "" && appEntity.ClientSecret != req.ClientSecret {
		return nil, code.GetError(code.OIDCInvalidClientError)
	}

	if authCode.CodeChallenge != "" && req.CodeVerifier == "" {
		return nil, code.GetError(code.OIDCPKCERequiredError)
	}

	if err := dao.NewAuthCodeDao().MarkUsed(ctx, req.Code); err != nil {
		glog.Errorf(ctx, "[tokenSvc.ExchangeCodeForToken] MarkUsed fail, err:%v", err)
	}

	return svc.generateTokenPair(ctx, authCode.PersonID, authCode.TenantID, authCode.OrgID, authCode.ClientID, authCode.Scope)
}

func (svc *tokenSvc) generateTokenPair(ctx context.Context, personID, tenantID, orgID uint, clientID, scopes string) (*dtooidc.TokenResp, error) {
	accessToken, err := svc.generateAccessToken(ctx, personID, tenantID, orgID, clientID, scopes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := svc.generateRefreshToken(ctx, personID, tenantID, orgID, clientID, scopes)
	if err != nil {
		return nil, err
	}

	return &dtooidc.TokenResp{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenExpireDuration.Seconds()),
		RefreshToken: refreshToken,
		Scope:        scopes,
	}, nil
}

func (svc *tokenSvc) generateAccessToken(ctx context.Context, personID, tenantID, orgID uint, clientID, scopes string) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.OIDCGenerateTokenError)
	}

	customData := gobject.UserClaims{
		PersonID: personID,
		TenantID: tenantID,
		OrgID:    orgID,
		UserType: "access",
	}

	tokenStr, err := jwtAuth.Issue(
		fmt.Sprintf("%d", personID),
		TokenIssuer,
		time.Now().Add(AccessTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.OIDCGenerateTokenError)
	}

	return tokenStr, nil
}

func (svc *tokenSvc) generateRefreshToken(ctx context.Context, personID, tenantID, orgID uint, clientID, scopes string) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.OIDCGenerateTokenError)
	}

	customData := gobject.UserClaims{
		PersonID: personID,
		TenantID: tenantID,
		OrgID:    orgID,
		UserType: "refresh",
	}

	tokenStr, err := jwtAuth.Issue(
		fmt.Sprintf("%d_refresh", personID),
		TokenIssuer,
		time.Now().Add(RefreshTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.OIDCGenerateTokenError)
	}

	return tokenStr, nil
}

func (svc *tokenSvc) RefreshAccessToken(ctx context.Context, req *dtooidc.TokenRefreshReq) (*dtooidc.TokenResp, error) {
	if req.RefreshToken == "" {
		return nil, code.GetError(code.OIDCInvalidRefreshTokenError)
	}

	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return nil, code.GetError(code.OIDCGenerateTokenError)
	}

	claims, err := jwtAuth.Parse(req.RefreshToken)
	if err != nil {
		glog.Errorf(ctx, "[tokenSvc.RefreshAccessToken] Parse refreshToken fail, err:%v", err)
		return nil, code.GetError(code.OIDCInvalidRefreshTokenError)
	}

	if claims.CustomData.UserType != "refresh" {
		return nil, code.GetError(code.OIDCInvalidRefreshTokenError)
	}

	personID := claims.CustomData.PersonID
	tenantID := claims.CustomData.TenantID
	orgID := claims.CustomData.OrgID

	tokenEntity, _ := dao.NewTokenDao().GetByRefreshTokenHash(ctx, svc.hashToken(req.RefreshToken))
	if tokenEntity != nil && tokenEntity.ClientID != "" {
		if err := dao.NewTokenDao().RevokeByRefreshTokenHash(ctx, svc.hashToken(req.RefreshToken)); err != nil {
			glog.Errorf(ctx, "[tokenSvc.RefreshAccessToken] RevokeByRefreshTokenHash fail, err:%v", err)
		}
	}

	scopes := ""
	if tokenEntity != nil {
		scopes = tokenEntity.Scopes
	}

	return svc.generateTokenPair(ctx, personID, tenantID, orgID, "", scopes)
}

func (svc *tokenSvc) ValidateAccessToken(ctx context.Context, accessToken string) (*model.TokenEntity, error) {
	if accessToken == "" {
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	claims, err := jwtAuth.Parse(accessToken)
	if err != nil {
		glog.Errorf(ctx, "[tokenSvc.ValidateAccessToken] Parse accessToken fail, err:%v", err)
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	if claims.CustomData.UserType != "access" {
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	tokenEntity, err := dao.NewTokenDao().GetByAccessTokenHash(ctx, svc.hashToken(accessToken))
	if err != nil {
		glog.Errorf(ctx, "[tokenSvc.ValidateAccessToken] GetByAccessTokenHash fail, err:%v", err)
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	if tokenEntity == nil {
		return nil, code.GetError(code.AuthTokenInvalidError)
	}

	return tokenEntity, nil
}

func (svc *tokenSvc) RevokeToken(ctx context.Context, token string, tokenTypeHint string) error {
	if token == "" {
		return nil
	}

	if tokenTypeHint == "access_token" || tokenTypeHint == "" {
		if err := dao.NewTokenDao().RevokeByAccessTokenHash(ctx, svc.hashToken(token)); err != nil {
			glog.Errorf(ctx, "[tokenSvc.RevokeToken] RevokeByAccessTokenHash fail, err:%v", err)
			return err
		}
	}

	if tokenTypeHint == "refresh_token" || tokenTypeHint == "" {
		if err := dao.NewTokenDao().RevokeByRefreshTokenHash(ctx, svc.hashToken(token)); err != nil {
			glog.Errorf(ctx, "[tokenSvc.RevokeToken] RevokeByRefreshTokenHash fail, err:%v", err)
			return err
		}
	}

	return nil
}

func (svc *tokenSvc) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
