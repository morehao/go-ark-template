package strategy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

type ssoStrategy struct {
	*strategyCommon
}

func NewSSOStrategy() RegisterStrategy {
	return &ssoStrategy{
		strategyCommon: newStrategyCommon(),
	}
}

func (s *ssoStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	ssoType := req.SSOType
	openID := req.OpenID
	if openID == "" {
		return nil, code.GetError(code.AuthRegisterError)
	}

	tenant, err := s.findTenantBySSO(ctx, orgEntity.ID, ssoType, openID)
	if err != nil {
		return nil, err
	}

	return s.createRegisterResult(ctx, orgEntity.ID, tenant.ID, req, 0)
}

func (s *ssoStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint, result *RegisterResult) error {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return err
	}
	if err := s.assignDefaultRolesAndDepts(ctx, orgEntity.ID, userID); err != nil {
		return err
	}
	return nil
}

func (s *ssoStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategySSO
}

func (s *ssoStrategy) findTenantBySSO(ctx *gin.Context, orgID uint, ssoType, openID string) (*model.TenantEntity, error) {
	ssoBind, err := dao.NewSSOBindDao().GetByCond(ctx, &dao.SSOBindCond{
		OrgID:   orgID,
		SSOType: ssoType,
		OpenID:  openID,
	})
	if err != nil {
		glog.Errorf(ctx, "[findTenantBySSO] GetByCond fail, ssoType:%s, err:%v", ssoType, err)
		return nil, code.GetError(code.AuthRegisterError)
	}
	if ssoBind != nil && ssoBind.TenantID > 0 {
		tenant, err := dao.NewTenantDao().GetByID(ctx, ssoBind.TenantID)
		if err != nil {
			glog.Errorf(ctx, "[findTenantBySSO] GetByID fail, tenantID:%d, err:%v", ssoBind.TenantID, err)
			return nil, code.GetError(code.TenantNotExistError)
		}
		if tenant == nil {
			return nil, code.GetError(code.TenantNotExistError)
		}
		return tenant, nil
	}

	tenantIDStr, err := s.getOrgConfigString(ctx, orgID, model.OrgConfigKeyRegisterSSODefaultTenantID)
	if err != nil {
		glog.Errorf(ctx, "[findTenantBySSO] GetString defaultTenantID fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}
	if tenantIDStr == "" {
		return nil, code.GetError(code.TenantNotExistError)
	}

	var tenantID uint
	if _, err := fmt.Sscanf(tenantIDStr, "%d", &tenantID); err != nil {
		return nil, code.GetError(code.AuthRegisterError)
	}
	tenant, err := dao.NewTenantDao().GetByID(ctx, tenantID)
	if err != nil {
		glog.Errorf(ctx, "[findTenantBySSO] GetByID fail, tenantID:%d, err:%v", tenantID, err)
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant == nil {
		return nil, code.GetError(code.TenantNotExistError)
	}
	return tenant, nil
}
