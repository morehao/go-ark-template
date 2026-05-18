package dao

import (
	"context"
	"errors"

	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type InviteCodeCond struct {
	*genericdao.BaseCond
	OrgID    uint
	TenantID uint
	Code     string
	Status   model.InviteCodeStatus
}

func (c *InviteCodeCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where(tableName + ".org_id = ?", c.OrgID)
	}
	if c.TenantID > 0 {
		db.Where(tableName + ".tenant_id = ?", c.TenantID)
	}
	if c.Code != "" {
		db.Where(tableName + ".code = ?", c.Code)
	}
	if c.Status != "" {
		db.Where(tableName + ".status = ?", c.Status)
	}
}

type InviteCodeDao struct {
	*genericdao.GenericDao[model.InviteCodeEntity, model.InviteCodeEntityList]
}

func NewInviteCodeDao() *InviteCodeDao {
	return &InviteCodeDao{
		GenericDao: genericdao.NewGenericDao[model.InviteCodeEntity, model.InviteCodeEntityList](
			model.TableNameInviteCode, "InviteCodeDao",
			dbclient.IamDB,
		),
	}
}

func (dao *InviteCodeDao) IncrUseCount(ctx context.Context, id uint) (int64, error) {
	db := dao.DB(ctx).Table(model.TableNameInviteCode)
	result := db.Where("id = ? AND status = ? AND (expired_at IS NULL OR expired_at > NOW()) AND (max_use_count = 0 OR use_count < max_use_count)", id, model.InviteCodeStatusActive).
		Update("use_count", gorm.Expr("use_count + 1"))
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("invite code not found or already disabled")
	}
	return result.RowsAffected, nil
}
