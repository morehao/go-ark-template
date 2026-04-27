package strategy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

type openStrategy struct {
	*strategyCommon
}

func NewOpenStrategy() RegisterStrategy {
	return &openStrategy{
		strategyCommon: newStrategyCommon(),
	}
}

func (s *openStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	tenantIDStr, err := s.getOrgConfigString(ctx, orgEntity.ID, model.OrgConfigKeyRegisterOpenTenantID)
	if err != nil {
		glog.Errorf(ctx, "[openStrategy.PreRegister] GetString tenantID fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}
	if tenantIDStr == "" {
		glog.Errorf(ctx, "[openStrategy.PreRegister] tenantID not configured")
		return nil, code.GetError(code.AuthRegisterError)
	}

	var tenantID uint
	if _, err := fmt.Sscanf(tenantIDStr, "%d", &tenantID); err != nil {
		glog.Errorf(ctx, "[openStrategy.PreRegister] parse tenantID fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}

	tenant, err := dao.NewTenantDao().GetByID(ctx, tenantID)
	if err != nil {
		glog.Errorf(ctx, "[openStrategy.PreRegister] GetByID tenant fail, err:%v, tenantID:%d", err, tenantID)
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant == nil {
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant.Status != model.TenantStatusEnabled {
		return nil, code.GetError(code.AuthRegisterError)
	}

	return s.createRegisterResult(ctx, orgEntity.ID, tenant.ID, req, 0)
}

func (s *openStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint, result *RegisterResult) error {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return err
	}
	return s.assignDefaultRolesAndDepts(ctx, orgEntity.ID, userID)
}

func (s *openStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyOpen
}