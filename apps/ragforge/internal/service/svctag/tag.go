package svctag

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtotag"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type TagSvc interface {
	Create(ctx *gin.Context, req *dtotag.TagCreateReq) (*dtotag.TagCreateResp, error)
	Delete(ctx *gin.Context, req *dtotag.TagDeleteReq) error
	Update(ctx *gin.Context, req *dtotag.TagUpdateReq) error
	List(ctx *gin.Context, req *dtotag.TagListReq) (*dtotag.TagListResp, error)
}

type tagSvc struct {
}

var _ TagSvc = (*tagSvc)(nil)

func NewTagSvc() TagSvc {
	return &tagSvc{}
}

func (svc *tagSvc) Create(ctx *gin.Context, req *dtotag.TagCreateReq) (*dtotag.TagCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	insertEntity := &model.TagEntity{
		KbID:     req.KbID,
		TenantID: tenantID,
		Name:     req.Name,
		Color:    req.Color,
	}
	if err := dao.NewTagDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svctag.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TagCreateError)
	}
	return &dtotag.TagCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *tagSvc) Delete(ctx *gin.Context, req *dtotag.TagDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewTagDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svctag.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TagDeleteError)
	}
	return nil
}

func (svc *tagSvc) Update(ctx *gin.Context, req *dtotag.TagUpdateReq) error {
	updateEntity := &model.TagEntity{
		Name:  req.Name,
		Color: req.Color,
	}
	if err := dao.NewTagDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svctag.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TagUpdateError)
	}
	return nil
}

func (svc *tagSvc) List(ctx *gin.Context, req *dtotag.TagListReq) (*dtotag.TagListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	cond := &dao.TagCond{
		BaseCond: &genericdao.BaseCond{},
		KbID:     req.KbID,
		TenantID: tenantID,
	}
	dataList, err := dao.NewTagDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svctag.List] dao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TagListError)
	}
	list := make([]dtotag.TagItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtotag.TagItem{
			ID:    v.ID,
			KbID:  v.KbID,
			Name:  v.Name,
			Color: v.Color,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				CreatedAt: v.CreatedAt.Unix(),
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtotag.TagListResp{
		List: list,
	}, nil
}
