package svclog

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtolog"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type OperationLogSvc interface {
	Create(ctx *gin.Context, req *dtolog.OperationLogCreateReq) (*dtolog.OperationLogCreateResp, error)
	PageList(ctx *gin.Context, req *dtolog.OperationLogPageListReq) (*dtolog.OperationLogPageListResp, error)
}

type operationLogSvc struct {
}

var _ OperationLogSvc = (*operationLogSvc)(nil)

func NewOperationLogSvc() OperationLogSvc {
	return &operationLogSvc{}
}

func (svc *operationLogSvc) Create(ctx *gin.Context, req *dtolog.OperationLogCreateReq) (*dtolog.OperationLogCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.OperationLogEntity{
		TenantID:      tenantID,
		UserID:        req.UserID,
		Username:      req.Username,
		Module:        req.Module,
		Operation:     req.Operation,
		Method:        req.Method,
		RequestID:     req.RequestID,
		RequestURL:    req.RequestURL,
		RequestParams: req.RequestParams,
		ResponseResult: req.ResponseResult,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
		ExecuteTime:   req.ExecuteTime,
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}

	if err := dao.NewOperationLogDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svclog.OperationLogCreate] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OperationLogCreateError)
	}

	return &dtolog.OperationLogCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *operationLogSvc) PageList(ctx *gin.Context, req *dtolog.OperationLogPageListReq) (*dtolog.OperationLogPageListResp, error) {
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
		glog.Errorf(ctx, "[svclog.OperationLogPageList] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OperationLogGetPageListError)
	}

	list := make([]dtolog.OperationLogPageListItem, 0, len(logList))
	for _, v := range logList {
		list = append(list, dtolog.OperationLogPageListItem{
			ID:            v.ID,
			UserID:        v.UserID,
			Username:      v.Username,
			Module:        v.Module,
			Operation:     v.Operation,
			Method:        v.Method,
			RequestURL:    v.RequestURL,
			IPAddress:     v.IPAddress,
			Status:        v.Status,
			ErrorMsg:      v.ErrorMsg,
			ExecuteTime:   v.ExecuteTime,
		})
	}

	return &dtolog.OperationLogPageListResp{
		List:  list,
		Total: total,
	}, nil
}
