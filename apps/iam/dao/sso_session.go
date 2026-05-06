package dao

import (
	"context"
	"time"

	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"gorm.io/gorm"
)

type SsoSessionDao struct {
}

func NewSsoSessionDao() *SsoSessionDao {
	return &SsoSessionDao{}
}

type SsoSessionCond struct {
	SessionID string
	PersonID  uint
	OrgID     uint
}

func (d *SsoSessionDao) GetBySessionID(ctx context.Context, sessionID string) (*model.SsoSessionEntity, error) {
	var entity model.SsoSessionEntity
	err := dbclient.IamDB(ctx).Where("session_id = ? AND deleted_at IS NULL", sessionID).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *SsoSessionDao) GetByCond(ctx context.Context, cond *SsoSessionCond) (*model.SsoSessionEntity, error) {
	query := dbclient.IamDB(ctx).Where("deleted_at IS NULL")
	if cond.SessionID != "" {
		query = query.Where("session_id = ?", cond.SessionID)
	}
	if cond.PersonID > 0 {
		query = query.Where("person_id = ?", cond.PersonID)
	}
	if cond.OrgID > 0 {
		query = query.Where("org_id = ?", cond.OrgID)
	}
	var entity model.SsoSessionEntity
	err := query.First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *SsoSessionDao) Insert(ctx context.Context, entity *model.SsoSessionEntity) error {
	return dbclient.IamDB(ctx).Create(entity).Error
}

func (d *SsoSessionDao) UpdateLastActiveTime(ctx context.Context, sessionID string) error {
	return dbclient.IamDB(ctx).Model(&model.SsoSessionEntity{}).
		Where("session_id = ?", sessionID).
		Update("last_active_time", time.Now()).Error
}

func (d *SsoSessionDao) DeleteBySessionID(ctx context.Context, sessionID string) error {
	return dbclient.IamDB(ctx).Where("session_id = ?", sessionID).Delete(&model.SsoSessionEntity{}).Error
}

func (d *SsoSessionDao) DeleteByPersonID(ctx context.Context, personID uint) error {
	return dbclient.IamDB(ctx).Where("person_id = ?", personID).Delete(&model.SsoSessionEntity{}).Error
}

func (d *SsoSessionDao) CleanExpired(ctx context.Context) error {
	return dbclient.IamDB(ctx).Where("expires_at < ?", time.Now()).Delete(&model.SsoSessionEntity{}).Error
}

func (d *SsoSessionDao) WithTx(tx *gorm.DB) *SsoSessionDao {
	return &SsoSessionDao{}
}
