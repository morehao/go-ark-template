package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/organization"
	"github.com/morehao/goark/apps/iam/core/user"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
	"gorm.io/gorm"
)

type OrganizationSvc interface {
	Create(ctx *gin.Context, req *dtoorg.OrganizationCreateReq) (*dtoorg.OrganizationCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorg.OrganizationDeleteReq) error
	Update(ctx *gin.Context, req *dtoorg.OrganizationUpdateReq) error
	Detail(ctx *gin.Context, req *dtoorg.OrganizationDetailReq) (*dtoorg.OrganizationDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.OrganizationPageListReq) (*dtoorg.OrganizationPageListResp, error)
}

type organizationSvc struct {
}

var _ OrganizationSvc = (*organizationSvc)(nil)

func NewOrganizationSvc() OrganizationSvc {
	return &organizationSvc{}
}

func (svc *organizationSvc) Create(ctx *gin.Context, req *dtoorg.OrganizationCreateReq) (*dtoorg.OrganizationCreateResp, error) {
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &iammodel.OrganizationEntity{
		Domain:           req.Domain,
		Logo:             req.Logo,
		Description:      req.Description,
		SortOrder:        req.SortOrder,
		Status:           req.Status,
		OrganizationCode: req.OrganizationCode,
		OrganizationName: req.OrganizationName,
		CreatedBy:        operatorID,
		UpdatedBy:        operatorID,
	}

	var adminID uint
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := iamdao.NewOrganizationDao().WithTx(tx).Insert(ctx, insertEntity); err != nil {
			glog.Errorf(ctx, "[svcorg.OrganizationCreate] Insert organization fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.TenantCreateError)
		}

		if req.Admin != nil && req.Admin.Username != "" {
			platformTenant, err := organization.GetPlatformTenant(ctx)
			if err != nil || platformTenant == nil || platformTenant.ID == 0 {
				glog.Errorf(ctx, "[svcorg.OrganizationCreate] GetPlatformTenant fail, err:%v", err)
				return code.GetError(code.TenantCreateError)
			}

			platformDept, err := organization.GetPlatformDept(ctx, platformTenant.ID)
			if err != nil || platformDept == nil || platformDept.ID == 0 {
				glog.Errorf(ctx, "[svcorg.OrganizationCreate] GetPlatformDept fail, err:%v", err)
				return code.GetError(code.TenantCreateError)
			}

			params := &user.CreatePersonParams{
				Mobile:     req.Admin.Mobile,
				Email:      req.Admin.Email,
				RealName:   req.Admin.RealName,
				Username:   req.Admin.Username,
				OperatorID: operatorID,
				TenantID:   platformTenant.ID,
				DeptID:     platformDept.ID,
				Status:     "active",
				UserType:   "tenant_admin",
			}
			result, err := user.CreatePersonWithUser(ctx, tx, params)
			if err != nil {
				glog.Errorf(ctx, "[svcorg.OrganizationCreate] CreatePersonWithUser fail, err:%v, req:%s", err, gutil.ToJsonString(req))
				return err
			}
			adminID = result.UserID
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationCreate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}
	return &dtoorg.OrganizationCreateResp{
		ID:      insertEntity.ID,
		AdminID: adminID,
	}, nil
}

func (svc *organizationSvc) Delete(ctx *gin.Context, req *dtoorg.OrganizationDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	organizationEntity, err := iamdao.NewOrganizationDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationDelete] daoOrganization GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	if organizationEntity == nil || organizationEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}
	if err = organization.CheckOrganizationAccess(ctx, organizationEntity.ID); err != nil {
		return err
	}

	if err = iamdao.NewOrganizationDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationDelete] daoOrganization Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

func (svc *organizationSvc) Update(ctx *gin.Context, req *dtoorg.OrganizationUpdateReq) error {
	if err := organization.CheckOrganizationAccess(ctx, req.ID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"domain":            req.Domain,
		"logo":              req.Logo,
		"description":       req.Description,
		"sort_order":        req.SortOrder,
		"status":            req.Status,
		"organization_code": req.OrganizationCode,
		"organization_name": req.OrganizationName,
	}
	if err := iamdao.NewOrganizationDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationUpdate] daoOrganization UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

func (svc *organizationSvc) Detail(ctx *gin.Context, req *dtoorg.OrganizationDetailReq) (*dtoorg.OrganizationDetailResp, error) {
	if err := organization.CheckOrganizationAccess(ctx, req.ID); err != nil {
		return nil, err
	}
	organizationEntity, err := iamdao.NewOrganizationDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationDetail] daoOrganization GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if organizationEntity == nil || organizationEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}
	resp := &dtoorg.OrganizationDetailResp{
		ID: organizationEntity.ID,
		OrganizationBaseInfo: objorg.OrganizationBaseInfo{
			Domain:           organizationEntity.Domain,
			Logo:             organizationEntity.Logo,
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

func (svc *organizationSvc) PageList(ctx *gin.Context, req *dtoorg.OrganizationPageListReq) (*dtoorg.OrganizationPageListResp, error) {
	cond := &iamdao.OrganizationCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	organizationEntityList, total, err := iamdao.NewOrganizationDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrganizationPageList] daoOrganization GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtoorg.OrganizationPageListItem, 0, len(organizationEntityList))
	for _, v := range organizationEntityList {
		list = append(list, dtoorg.OrganizationPageListItem{
			ID: v.ID,
			OrganizationBaseInfo: objorg.OrganizationBaseInfo{
				Domain:           v.Domain,
				Logo:             v.Logo,
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
	return &dtoorg.OrganizationPageListResp{
		List:  list,
		Total: total,
	}, nil
}
