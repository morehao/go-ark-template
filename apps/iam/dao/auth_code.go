package dao

import (
	"context"
	"time"

	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"gorm.io/gorm"
)

type AuthCodeDao struct {
}

func NewAuthCodeDao() *AuthCodeDao {
	return &AuthCodeDao{}
}

type AuthCodeCond struct {
	Code     string
	ClientID string
}

func (d *AuthCodeDao) GetByCode(ctx context.Context, code string) (*model.AuthCodeEntity, error) {
	var entity model.AuthCodeEntity
	err := dbclient.IamDB(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *AuthCodeDao) Insert(ctx context.Context, entity *model.AuthCodeEntity) error {
	return dbclient.IamDB(ctx).Create(entity).Error
}

func (d *AuthCodeDao) MarkUsed(ctx context.Context, code string) error {
	now := time.Now()
	return dbclient.IamDB(ctx).Model(&model.AuthCodeEntity{}).
		Where("code = ?", code).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": now,
		}).Error
}

func (d *AuthCodeDao) CleanExpired(ctx context.Context) error {
	return dbclient.IamDB(ctx).Where("expires_at < ? OR (used = ? AND used_at < ?)",
		time.Now(), true, time.Now().Add(-10*time.Minute)).Delete(&model.AuthCodeEntity{}).Error
}

func (d *AuthCodeDao) WithTx(tx *gorm.DB) *AuthCodeDao {
	return &AuthCodeDao{}
}
