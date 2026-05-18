package svcuser

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/internal/dto/dtouser"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type LoginLogSvc interface {
	Create(ctx *gin.Context, req *dtouser.LoginLogCreateReq) (*dtouser.LoginLogCreateResp, error)
	PageList(ctx *gin.Context, req *dtouser.LoginLogPageListReq) (*dtouser.LoginLogPageListResp, error)
}

type loginLogSvc struct {
}

var _ LoginLogSvc = (*loginLogSvc)(nil)

func NewLoginLogSvc() LoginLogSvc {
	return &loginLogSvc{}
}

func (svc *loginLogSvc) Create(ctx *gin.Context, req *dtouser.LoginLogCreateReq) (*dtouser.LoginLogCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.LoginLogEntity{
		TenantID:     tenantID,
		UserID:       req.UserID,
		Username:     req.Username,
		LoginType:    req.LoginType,
		LoginStatus:  req.LoginStatus,
		LoginMessage: req.LoginMessage,
		IPAddress:    req.IPAddress,
		Location:     req.Location,
		Browser:      req.Browser,
		OS:           req.OS,
		CreatedBy:    operatorID,
		UpdatedBy:    operatorID,
	}

	if err := dao.NewLoginLogDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.LoginLogCreate] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogCreateError)
	}

	return &dtouser.LoginLogCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *loginLogSvc) PageList(ctx *gin.Context, req *dtouser.LoginLogPageListReq) (*dtouser.LoginLogPageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &dao.LoginLogCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID:    tenantID,
		UserID:      req.UserID,
		Username:    req.Username,
		LoginType:   req.LoginType,
		LoginStatus: req.LoginStatus,
	}

	logList, total, err := dao.NewLoginLogDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.LoginLogPageList] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogGetPageListError)
	}

	list := make([]dtouser.LoginLogPageListItem, 0, len(logList))
	for _, v := range logList {
		list = append(list, dtouser.LoginLogPageListItem{
			ID:           v.ID,
			UserID:       v.UserID,
			Username:     v.Username,
			LoginType:    v.LoginType,
			LoginStatus:  v.LoginStatus,
			LoginMessage: v.LoginMessage,
			IPAddress:    v.IPAddress,
			Location:     v.Location,
			Browser:      v.Browser,
			OS:           v.OS,
		})
	}

	return &dtouser.LoginLogPageListResp{
		List:  list,
		Total: total,
	}, nil
}