package strategy

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

type inviteStrategy struct {
	*strategyCommon
}

func NewInviteStrategy() RegisterStrategy {
	return &inviteStrategy{
		strategyCommon: newStrategyCommon(),
	}
}

func (s *inviteStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	inviteCode := strings.TrimSpace(req.InviteCode)
	if inviteCode == "" {
		return nil, code.GetError(code.InviteCodeRequiredError)
	}

	invite, err := dao.NewInviteCodeDao().GetByCond(ctx, &dao.InviteCodeCond{
		OrgID: orgEntity.ID,
		Code:  inviteCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[inviteStrategy.PreRegister] GetByCond fail, code:%s, err:%v", inviteCode, err)
		return nil, code.GetError(code.InviteCodeInvalidError)
	}
	if invite == nil {
		return nil, code.GetError(code.InviteCodeInvalidError)
	}

	if invite.Status != model.InviteCodeStatusActive {
		return nil, code.GetError(code.InviteCodeInvalidError)
	}
	if invite.ExpiredAt != nil && invite.ExpiredAt.Before(time.Now()) {
		return nil, code.GetError(code.InviteCodeExpiredError)
	}

	maxUse, err := s.getOrgConfigInt(ctx, orgEntity.ID, model.OrgConfigKeyRegisterCodeMaxUse)
	if err != nil {
		glog.Errorf(ctx, "[inviteStrategy.PreRegister] getOrgConfigInt fail, err:%v", err)
	}
	if maxUse > 0 && invite.UseCount >= maxUse {
		return nil, code.GetError(code.InviteCodeUsedUpError)
	}

	tenant, err := dao.NewTenantDao().GetByID(ctx, invite.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[inviteStrategy.PreRegister] GetByID tenant fail, tenantID:%d, err:%v", invite.TenantID, err)
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant == nil {
		return nil, code.GetError(code.TenantNotExistError)
	}
	if tenant.Status != model.TenantStatusEnabled {
		return nil, code.GetError(code.AuthRegisterError)
	}

	return s.createRegisterResult(ctx, orgEntity.ID, tenant.ID, req, invite.ID)
}

func (s *inviteStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint, result *RegisterResult) error {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return err
	}
	if err := s.assignDefaultRolesAndDepts(ctx, orgEntity.ID, userID); err != nil {
		return err
	}
	if result.InviteID > 0 {
		if _, err := dao.NewInviteCodeDao().IncrUseCount(ctx, result.InviteID); err != nil {
			glog.Errorf(ctx, "[inviteStrategy.PostRegister] IncrUseCount fail, inviteID:%d, err:%v", result.InviteID, err)
			return code.GetError(code.UserUpdateError)
		}
	}
	return nil
}

func (s *inviteStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyInvite
}