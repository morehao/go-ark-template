package strategy

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

type domainStrategy struct {
	*strategyCommon
}

func NewDomainStrategy() RegisterStrategy {
	return &domainStrategy{
		strategyCommon: newStrategyCommon(),
	}
}

func (s *domainStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	domain := resolveDomain(ctx)
	if domain == "" {
		return nil, code.GetError(code.TenantNotExistError)
	}

	tenant, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
		OrgID:  orgEntity.ID,
		Domain: domain,
	})
	if err != nil {
		glog.Errorf(ctx, "[domainStrategy.PreRegister] GetByCond fail, domain:%s, err:%v", domain, err)
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant == nil {
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant.Status != model.TenantStatusEnabled {
		return nil, code.GetError(code.AuthRegisterError)
	}

	return s.createRegisterResult(ctx, orgEntity.ID, tenant.ID, req)
}

func (s *domainStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint) error {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return err
	}
	return s.assignDefaultRolesAndDepts(ctx, orgEntity.ID, userID)
}

func (s *domainStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyDomain
}


