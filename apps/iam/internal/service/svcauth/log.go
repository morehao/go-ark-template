package svcauth

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type LoginLogSvc interface {
	Create(ctx *gin.Context, req *dtoauth.LoginLogCreateReq) (*dtoauth.LoginLogCreateResp, error)
	PageList(ctx *gin.Context, req *dtoauth.LoginLogPageListReq) (*dtoauth.LoginLogPageListResp, error)
}

type loginLogSvc struct {
}

var _ LoginLogSvc = (*loginLogSvc)(nil)

func NewLoginLogSvc() LoginLogSvc {
	return &loginLogSvc{}
}

func (svc *loginLogSvc) Create(ctx *gin.Context, req *dtoauth.LoginLogCreateReq) (*dtoauth.LoginLogCreateResp, error) {
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
		glog.Errorf(ctx, "[svcauth.LoginLogCreate] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogCreateError)
	}

	return &dtoauth.LoginLogCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *loginLogSvc) PageList(ctx *gin.Context, req *dtoauth.LoginLogPageListReq) (*dtoauth.LoginLogPageListResp, error) {
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
		glog.Errorf(ctx, "[svcauth.LoginLogPageList] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogGetPageListError)
	}

	list := make([]dtoauth.LoginLogPageListItem, 0, len(logList))
	for _, v := range logList {
		list = append(list, dtoauth.LoginLogPageListItem{
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

	return &dtoauth.LoginLogPageListResp{
		List:  list,
		Total: total,
	}, nil
}

type OperationLogSvc interface {
	Create(ctx *gin.Context, req *dtoauth.OperationLogCreateReq) (*dtoauth.OperationLogCreateResp, error)
	PageList(ctx *gin.Context, req *dtoauth.OperationLogPageListReq) (*dtoauth.OperationLogPageListResp, error)
}

type operationLogSvc struct {
}

var _ OperationLogSvc = (*operationLogSvc)(nil)

func NewOperationLogSvc() OperationLogSvc {
	return &operationLogSvc{}
}

func (svc *operationLogSvc) Create(ctx *gin.Context, req *dtoauth.OperationLogCreateReq) (*dtoauth.OperationLogCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.OperationLogEntity{
		TenantID:       tenantID,
		UserID:         req.UserID,
		Username:       req.Username,
		Module:         req.Module,
		Operation:      req.Operation,
		Method:         req.Method,
		RequestID:      req.RequestID,
		RequestURL:     req.RequestURL,
		RequestParams:  req.RequestParams,
		ResponseResult: req.ResponseResult,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		Status:         req.Status,
		ErrorMsg:       req.ErrorMsg,
		ExecuteTime:    req.ExecuteTime,
		CreatedBy:      operatorID,
		UpdatedBy:      operatorID,
	}

	if err := dao.NewOperationLogDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcauth.OperationLogCreate] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OperationLogCreateError)
	}

	return &dtoauth.OperationLogCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *operationLogSvc) PageList(ctx *gin.Context, req *dtoauth.OperationLogPageListReq) (*dtoauth.OperationLogPageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &dao.OperationLogCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID:   tenantID,
		UserID:     req.UserID,
		Username:   req.Username,
		Module:     req.Module,
		Operation:  req.Operation,
		Status:     req.Status,
	}

	logList, total, err := dao.NewOperationLogDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.OperationLogPageList] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OperationLogGetPageListError)
	}

	list := make([]dtoauth.OperationLogPageListItem, 0, len(logList))
	for _, v := range logList {
		list = append(list, dtoauth.OperationLogPageListItem{
			ID:          v.ID,
			UserID:      v.UserID,
			Username:    v.Username,
			Module:      v.Module,
			Operation:   v.Operation,
			Method:      v.Method,
			RequestURL:  v.RequestURL,
			IPAddress:   v.IPAddress,
			Status:      v.Status,
			ErrorMsg:    v.ErrorMsg,
			ExecuteTime: v.ExecuteTime,
		})
	}

	return &dtoauth.OperationLogPageListResp{
		List:  list,
		Total: total,
	}, nil
}