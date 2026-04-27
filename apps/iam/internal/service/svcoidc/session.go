package svcoidc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/glog"
)

const (
	SsoSessionExpireDuration = 24 * time.Hour
)

type SsoSessionSvc interface {
	CreateSession(ctx context.Context, personID, orgID uint) (string, error)
	GetValidSession(ctx context.Context, sessionID string) (*model.SsoSessionEntity, error)
	RefreshSession(ctx context.Context, sessionID string) error
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteSessionByPersonID(ctx context.Context, personID uint) error
}

type ssoSessionSvc struct {
}

func NewSsoSessionSvc() SsoSessionSvc {
	return &ssoSessionSvc{}
}

func (svc *ssoSessionSvc) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (svc *ssoSessionSvc) CreateSession(ctx context.Context, personID, orgID uint) (string, error) {
	sessionID, err := svc.generateSessionID()
	if err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.CreateSession] generateSessionID fail, err:%v", err)
		return "", err
	}

	now := time.Now()
	entity := &model.SsoSessionEntity{
		SessionID:      sessionID,
		PersonID:       personID,
		OrgID:          orgID,
		LoginTime:      now,
		LastActiveTime: now,
		ExpiresAt:      now.Add(SsoSessionExpireDuration),
	}

	if err := dao.NewSsoSessionDao().Insert(ctx, entity); err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.CreateSession] Insert fail, err:%v", err)
		return "", err
	}

	return sessionID, nil
}

func (svc *ssoSessionSvc) GetValidSession(ctx context.Context, sessionID string) (*model.SsoSessionEntity, error) {
	if sessionID == "" {
		return nil, nil
	}

	entity, err := dao.NewSsoSessionDao().GetBySessionID(ctx, sessionID)
	if err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.GetValidSession] GetBySessionID fail, err:%v", err)
		return nil, err
	}
	if entity == nil || entity.ID == 0 {
		return nil, nil
	}

	if time.Now().After(entity.ExpiresAt) {
		return nil, nil
	}

	return entity, nil
}

func (svc *ssoSessionSvc) RefreshSession(ctx context.Context, sessionID string) error {
	entity, err := dao.NewSsoSessionDao().GetBySessionID(ctx, sessionID)
	if err != nil || entity == nil || entity.ID == 0 {
		return err
	}

	newExpiresAt := time.Now().Add(SsoSessionExpireDuration)
	if err := dbclient.IamDB(ctx).Model(&model.SsoSessionEntity{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active_time": time.Now(),
			"expires_at":       newExpiresAt,
		}).Error; err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.RefreshSession] Update fail, err:%v", err)
		return err
	}

	return nil
}

func (svc *ssoSessionSvc) DeleteSession(ctx context.Context, sessionID string) error {
	if err := dao.NewSsoSessionDao().DeleteBySessionID(ctx, sessionID); err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.DeleteSession] DeleteBySessionID fail, err:%v", err)
		return err
	}
	return nil
}

func (svc *ssoSessionSvc) DeleteSessionByPersonID(ctx context.Context, personID uint) error {
	if err := dao.NewSsoSessionDao().DeleteByPersonID(ctx, personID); err != nil {
		glog.Errorf(ctx, "[ssoSessionSvc.DeleteSessionByPersonID] DeleteByPersonID fail, err:%v", err)
		return err
	}
	return nil
}
