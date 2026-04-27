package svcoidc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/golib/glog"
)

type LogoutSvc interface {
	Logout(ctx context.Context, sessionID string, refreshToken string) error
}

type logoutSvc struct {
}

func NewLogoutSvc() LogoutSvc {
	return &logoutSvc{}
}

func (svc *logoutSvc) Logout(ctx context.Context, sessionID string, refreshToken string) error {
	if sessionID != "" {
		if err := dao.NewSsoSessionDao().DeleteBySessionID(ctx, sessionID); err != nil {
			glog.Errorf(ctx, "[logoutSvc.Logout] DeleteBySessionID fail, err:%v", err)
		}
	}

	if refreshToken != "" {
		tokenHash := svc.hashToken(refreshToken)
		if err := dao.NewTokenDao().RevokeByRefreshTokenHash(ctx, tokenHash); err != nil {
			glog.Errorf(ctx, "[logoutSvc.Logout] RevokeByRefreshTokenHash fail, err:%v", err)
		}
	}

	return nil
}

func (svc *logoutSvc) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}