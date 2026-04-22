package svcapplication

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoapplication"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/apps/iam/object/objapplication"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type ApplicationSvc interface {
	Create(ctx *gin.Context, req *dtoapplication.ApplicationCreateReq) (*dtoapplication.ApplicationCreateResp, error)
	Delete(ctx *gin.Context, req *dtoapplication.ApplicationDeleteReq) error
	Update(ctx *gin.Context, req *dtoapplication.ApplicationUpdateReq) error
	Detail(ctx *gin.Context, req *dtoapplication.ApplicationDetailReq) (*dtoapplication.ApplicationDetailResp, error)
	PageList(ctx *gin.Context, req *dtoapplication.ApplicationPageListReq) (*dtoapplication.ApplicationPageListResp, error)
}

type applicationSvc struct {
}

var _ ApplicationSvc = (*applicationSvc)(nil)

func NewApplicationSvc() ApplicationSvc {
	return &applicationSvc{}
}

// Create 创建应用管理
func (svc *applicationSvc) Create(ctx *gin.Context, req *dtoapplication.ApplicationCreateReq) (*dtoapplication.ApplicationCreateResp, error) {
	insertEntity := &model.ApplicationEntity{
		AppCode:     req.AppCode,
		AppName:     req.AppName,
		AppType:     req.AppType,
		CallbackUrl: req.CallbackUrl,
		Description: req.Description,
		HomepageUrl: req.HomepageUrl,
		Logo:        req.Logo,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}

	if err := dao.NewApplicationDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationCreate] dao Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ApplicationCreateError)
	}
	return &dtoapplication.ApplicationCreateResp{
		AppID: insertEntity.ID,
	}, nil
}

// Delete 删除应用管理
func (svc *applicationSvc) Delete(ctx *gin.Context, req *dtoapplication.ApplicationDeleteReq) error {
	applicationEntity, err := dao.NewApplicationDao().GetByID(ctx, req.AppID)
	if err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationDelete] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ApplicationDeleteError)
	}
	if applicationEntity == nil || applicationEntity.ID == 0 {
		return code.GetError(code.ApplicationNotExistError)
	}

	userID := gincontext.GetUserID(ctx)

	if err := dao.NewApplicationDao().Delete(ctx, req.AppID, userID); err != nil {
		glog.Errorf(ctx, "[svcapplication.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ApplicationDeleteError)
	}
	return nil
}

// Update 更新应用管理
func (svc *applicationSvc) Update(ctx *gin.Context, req *dtoapplication.ApplicationUpdateReq) error {
	applicationEntity, err := dao.NewApplicationDao().GetByID(ctx, req.AppID)
	if err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationUpdate] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ApplicationUpdateError)
	}
	if applicationEntity == nil || applicationEntity.ID == 0 {
		return code.GetError(code.ApplicationNotExistError)
	}

	updateMap := map[string]any{}
	if err := dao.NewApplicationDao().UpdateMap(ctx, req.AppID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationUpdate] dao UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ApplicationUpdateError)
	}
	return nil
}

// Detail 根据id获取应用管理
func (svc *applicationSvc) Detail(ctx *gin.Context, req *dtoapplication.ApplicationDetailReq) (*dtoapplication.ApplicationDetailResp, error) {
	applicationEntity, err := dao.NewApplicationDao().GetByID(ctx, req.AppID)
	if err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationDetail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ApplicationGetDetailError)
	}
	if applicationEntity == nil || applicationEntity.ID == 0 {
		return nil, code.GetError(code.ApplicationNotExistError)
	}
	resp := &dtoapplication.ApplicationDetailResp{
		AppID: applicationEntity.ID,
		ApplicationBaseInfo: objapplication.ApplicationBaseInfo{
			AppCode:     applicationEntity.AppCode,
			AppName:     applicationEntity.AppName,
			AppType:     applicationEntity.AppType,
			CallbackUrl: applicationEntity.CallbackUrl,
			Description: applicationEntity.Description,
			HomepageUrl: applicationEntity.HomepageUrl,
			Logo:        applicationEntity.Logo,
			SortOrder:   applicationEntity.SortOrder,
			Status:      applicationEntity.Status,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: applicationEntity.CreatedAt.Unix(),
			UpdatedAt: applicationEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取应用管理列表
func (svc *applicationSvc) PageList(ctx *gin.Context, req *dtoapplication.ApplicationPageListReq) (*dtoapplication.ApplicationPageListResp, error) {
	cond := &dao.ApplicationCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	applicationEntityList, total, err := dao.NewApplicationDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcapplication.ApplicationPageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ApplicationGetPageListError)
	}
	list := make([]dtoapplication.ApplicationPageListItem, 0, len(applicationEntityList))
	for _, v := range applicationEntityList {
		list = append(list, dtoapplication.ApplicationPageListItem{
			AppID: v.ID,
			ApplicationBaseInfo: objapplication.ApplicationBaseInfo{
				AppCode:     v.AppCode,
				AppName:     v.AppName,
				AppType:     v.AppType,
				CallbackUrl: v.CallbackUrl,
				Description: v.Description,
				HomepageUrl: v.HomepageUrl,
				Logo:        v.Logo,
				SortOrder:   v.SortOrder,
				Status:      v.Status,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoapplication.ApplicationPageListResp{
		List:  list,
		Total: total,
	}, nil
}
