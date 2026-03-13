package svctenant

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtotenant"
	"github.com/morehao/goark/apps/iam/object/objtenant"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type TenantSvc interface {
	Create(ctx *gin.Context, req *dtotenant.TenantCreateReq) (*dtotenant.TenantCreateResp, error)
	Delete(ctx *gin.Context, req *dtotenant.TenantDeleteReq) error
	Update(ctx *gin.Context, req *dtotenant.TenantUpdateReq) error
	Detail(ctx *gin.Context, req *dtotenant.TenantDetailReq) (*dtotenant.TenantDetailResp, error)
	PageList(ctx *gin.Context, req *dtotenant.TenantPageListReq) (*dtotenant.TenantPageListResp, error)
}

type tenantSvc struct {
}

var _ TenantSvc = (*tenantSvc)(nil)

func NewTenantSvc() TenantSvc {
	return &tenantSvc{}
}

// Create 创建租户管理
func (svc *tenantSvc) Create(ctx *gin.Context, req *dtotenant.TenantCreateReq) (*dtotenant.TenantCreateResp, error) {
	insertEntity := &iammodel.TenantEntity{
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		TenantCode:  req.TenantCode,
		TenantName:  req.TenantName,
	}

	if err := iamdao.NewTenantDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svctenant.TenantCreate] daoTenant Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}
	return &dtotenant.TenantCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除租户管理
func (svc *tenantSvc) Delete(ctx *gin.Context, req *dtotenant.TenantDeleteReq) error {
	userID := gincontext.GetUserID(ctx)

	if err := iamdao.NewTenantDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svctenant.Delete] daoTenant Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

// Update 更新租户管理
func (svc *tenantSvc) Update(ctx *gin.Context, req *dtotenant.TenantUpdateReq) error {

	updateEntity := &iammodel.TenantEntity{
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		TenantCode:  req.TenantCode,
		TenantName:  req.TenantName,
	}
	if err := iamdao.NewTenantDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svctenant.TenantUpdate] daoTenant UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

// Detail 根据id获取租户管理
func (svc *tenantSvc) Detail(ctx *gin.Context, req *dtotenant.TenantDetailReq) (*dtotenant.TenantDetailResp, error) {
	detailEntity, err := iamdao.NewTenantDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svctenant.TenantDetail] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}
	resp := &dtotenant.TenantDetailResp{
		ID: detailEntity.ID,
		TenantBaseInfo: objtenant.TenantBaseInfo{
			Description: detailEntity.Description,
			SortOrder:   detailEntity.SortOrder,
			Status:      detailEntity.Status,
			TenantCode:  detailEntity.TenantCode,
			TenantName:  detailEntity.TenantName,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取租户管理列表
func (svc *tenantSvc) PageList(ctx *gin.Context, req *dtotenant.TenantPageListReq) (*dtotenant.TenantPageListResp, error) {
	cond := &iamdao.TenantCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	dataList, total, err := iamdao.NewTenantDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svctenant.TenantPageList] daoTenant GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtotenant.TenantPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtotenant.TenantPageListItem{
			ID: v.ID,
			TenantBaseInfo: objtenant.TenantBaseInfo{
				Description: v.Description,
				SortOrder:   v.SortOrder,
				Status:      v.Status,
				TenantCode:  v.TenantCode,
				TenantName:  v.TenantName,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtotenant.TenantPageListResp{
		List:  list,
		Total: total,
	}, nil
}
