package svcoidc

import (
	"context"
	"fmt"

	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtooidc"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/glog"
)

type UserInfoSvc interface {
	GetUserInfo(ctx context.Context, personID uint) (*dtooidc.UserInfoResp, error)
}

type userInfoSvc struct {
}

func NewUserInfoSvc() UserInfoSvc {
	return &userInfoSvc{}
}

func (svc *userInfoSvc) GetUserInfo(ctx context.Context, personID uint) (*dtooidc.UserInfoResp, error) {
	personEntity, err := dao.NewPersonDao().GetByID(ctx, personID)
	if err != nil {
		glog.Errorf(ctx, "[userInfoSvc.GetUserInfo] GetByID person fail, err:%v", err)
		return nil, code.GetError(code.UserNotExistError)
	}
	if personEntity == nil || personEntity.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}

	return &dtooidc.UserInfoResp{
		Subject: fmt.Sprintf("%d", personEntity.ID),
		Name:    personEntity.RealName,
		Email:   personEntity.Email,
		Phone:   personEntity.Mobile,
	}, nil
}
