package dao

import (
	"context"
	"time"

	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"gorm.io/gorm"
)

type TokenDao struct {
	tx *gorm.DB
}

func NewTokenDao() *TokenDao {
	return &TokenDao{}
}

type TokenCond struct {
	TokenID           string
	PersonID          uint
	ClientID          string
	TenantID          uint
	AccessTokenHash   string
	RefreshTokenHash  string
	Revoked           bool
}

func (d *TokenDao) GetByTokenID(ctx context.Context, tokenID string) (*model.TokenEntity, error) {
	var entity model.TokenEntity
	err := d.db(ctx).Where("token_id = ? AND deleted_at IS NULL", tokenID).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *TokenDao) GetByAccessTokenHash(ctx context.Context, hash string) (*model.TokenEntity, error) {
	var entity model.TokenEntity
	err := d.db(ctx).Where("access_token_hash = ? AND token_type = ? AND revoked = ? AND deleted_at IS NULL",
		hash, model.TokenTypeAccess, false).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *TokenDao) GetByRefreshTokenHash(ctx context.Context, hash string) (*model.TokenEntity, error) {
	var entity model.TokenEntity
	err := d.db(ctx).Where("refresh_token_hash = ? AND token_type = ? AND revoked = ? AND deleted_at IS NULL",
		hash, model.TokenTypeRefresh, false).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (d *TokenDao) Insert(ctx context.Context, entity *model.TokenEntity) error {
	return d.db(ctx).Create(entity).Error
}

func (d *TokenDao) RevokeByTokenID(ctx context.Context, tokenID string) error {
	now := time.Now()
	return d.db(ctx).Model(&model.TokenEntity{}).
		Where("token_id = ?", tokenID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		}).Error
}

func (d *TokenDao) RevokeByPersonID(ctx context.Context, personID uint) error {
	now := time.Now()
	return d.db(ctx).Model(&model.TokenEntity{}).
		Where("person_id = ?", personID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		}).Error
}

func (d *TokenDao) RevokeByRefreshTokenHash(ctx context.Context, hash string) error {
	now := time.Now()
	return d.db(ctx).Model(&model.TokenEntity{}).
		Where("refresh_token_hash = ? AND token_type = ?", hash, model.TokenTypeRefresh).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		}).Error
}

func (d *TokenDao) RevokeByAccessTokenHash(ctx context.Context, hash string) error {
	now := time.Now()
	return d.db(ctx).Model(&model.TokenEntity{}).
		Where("access_token_hash = ? AND token_type = ?", hash, model.TokenTypeAccess).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		}).Error
}

func (d *TokenDao) CleanExpired(ctx context.Context) error {
	return d.db(ctx).Where("expires_at < ? AND revoked = ?", time.Now(), true).
		Delete(&model.TokenEntity{}).Error
}

func (d *TokenDao) WithTx(tx *gorm.DB) *TokenDao {
	return &TokenDao{tx: tx}
}

func (d *TokenDao) db(ctx context.Context) *gorm.DB {
	if d.tx != nil {
		return d.tx
	}
	return dbclient.IamDB(ctx)
}
