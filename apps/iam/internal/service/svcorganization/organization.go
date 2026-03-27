package svcorganization

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/organization"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorganization"
	"github.com/morehao/goark/apps/iam/object/objorganization"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type OrganizationSvc interface {
	Create(ctx *gin.Context, req *dtoorganization.OrganizationCreateReq) (*dtoorganization.OrganizationCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorganization.OrganizationDeleteReq) error
	Update(ctx *gin.Context, req *dtoorganization.OrganizationUpdateReq) error
	Detail(ctx *gin.Context, req *dtoorganization.OrganizationDetailReq) (*dtoorganization.OrganizationDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorganization.OrganizationPageListReq) (*dtoorganization.OrganizationPageListResp, error)
}

type organizationSvc struct {
}

var _ OrganizationSvc = (*organizationSvc)(nil)

func NewOrganizationSvc() OrganizationSvc {
	return &organizationSvc{}
}

func (svc *organizationSvc) Create(ctx *gin.Context, req *dtoorganization.OrganizationCreateReq) (*dtoorganization.OrganizationCreateResp, error) {
	insertEntity := &iammodel.OrganizationEntity{
		Description:      req.Description,
		SortOrder:        req.SortOrder,
		Status:           req.Status,
		OrganizationCode: req.OrganizationCode,
		OrganizationName: req.OrganizationName,
	}

	if err := iamdao.NewOrganizationDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcorganization.OrganizationCreate] daoOrganization Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}
	return &dtoorganization.OrganizationCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *organizationSvc) Delete(ctx *gin.Context, req *dtoorganization.OrganizationDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	organizationEntity, err := iamdao.NewOrganizationDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorganization.Delete] daoOrganization GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	if organizationEntity == nil || organizationEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}
	if err = organization.CheckOrganizationAccess(ctx, organizationEntity.ID); err != nil {
		return err
	}

	if err = iamdao.NewOrganizationDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorganization.Delete] daoOrganization Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

func (svc *organizationSvc) Update(ctx *gin.Context, req *dtoorganization.OrganizationUpdateReq) error {
	if err := organization.CheckOrganizationAccess(ctx, req.ID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"description":       req.Description,
		"sort_order":        req.SortOrder,
		"status":            req.Status,
		"organization_code": req.OrganizationCode,
		"organization_name": req.OrganizationName,
	}
	if err := iamdao.NewOrganizationDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcorganization.OrganizationUpdate] daoOrganization UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

func (svc *organizationSvc) Detail(ctx *gin.Context, req *dtoorganization.OrganizationDetailReq) (*dtoorganization.OrganizationDetailResp, error) {
	if err := organization.CheckOrganizationAccess(ctx, req.ID); err != nil {
		return nil, err
	}
	organizationEntity, err := iamdao.NewOrganizationDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorganization.OrganizationDetail] daoOrganization GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if organizationEntity == nil || organizationEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}
	resp := &dtoorganization.OrganizationDetailResp{
		ID: organizationEntity.ID,
		OrganizationBaseInfo: objorganization.OrganizationBaseInfo{
			Description:      organizationEntity.Description,
			SortOrder:        organizationEntity.SortOrder,
			Status:           organizationEntity.Status,
			OrganizationCode: organizationEntity.OrganizationCode,
			OrganizationName: organizationEntity.OrganizationName,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: organizationEntity.CreatedAt.Unix(),
			UpdatedAt: organizationEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *organizationSvc) PageList(ctx *gin.Context, req *dtoorganization.OrganizationPageListReq) (*dtoorganization.OrganizationPageListResp, error) {
	cond := &iamdao.OrganizationCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	organizationEntityList, total, err := iamdao.NewOrganizationDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorganization.OrganizationPageList] daoOrganization GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtoorganization.OrganizationPageListItem, 0, len(organizationEntityList))
	for _, v := range organizationEntityList {
		list = append(list, dtoorganization.OrganizationPageListItem{
			ID: v.ID,
			OrganizationBaseInfo: objorganization.OrganizationBaseInfo{
				Description:      v.Description,
				SortOrder:        v.SortOrder,
				Status:           v.Status,
				OrganizationCode: v.OrganizationCode,
				OrganizationName: v.OrganizationName,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoorganization.OrganizationPageListResp{
		List:  list,
		Total: total,
	}, nil
}
