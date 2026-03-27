package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/organization"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type TenantSvc interface {
	Create(ctx *gin.Context, req *dtoorg.TenantCreateReq) (*dtoorg.TenantCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorg.TenantDeleteReq) error
	Update(ctx *gin.Context, req *dtoorg.TenantUpdateReq) error
	Detail(ctx *gin.Context, req *dtoorg.TenantDetailReq) (*dtoorg.TenantDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.TenantPageListReq) (*dtoorg.TenantPageListResp, error)
}

type tenantSvc struct {
}

var _ TenantSvc = (*tenantSvc)(nil)

func NewTenantSvc() TenantSvc {
	return &tenantSvc{}
}

func (svc *tenantSvc) Create(ctx *gin.Context, req *dtoorg.TenantCreateReq) (*dtoorg.TenantCreateResp, error) {
	insertEntity := &iammodel.TenantEntity{
		Address:                 req.Address,
		ContactEmail:            req.ContactEmail,
		ContactPhone:            req.ContactPhone,
		LegalPerson:             req.LegalPerson,
		Logo:                    req.Logo,
		OrganizationID:          req.OrganizationID,
		ShortName:               req.ShortName,
		Status:                  req.Status,
		TenantCode:              req.TenantCode,
		TenantName:              req.TenantName,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
	}

	if err := iamdao.NewTenantDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.TenantCreate] daoTenant Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}
	return &dtoorg.TenantCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *tenantSvc) Delete(ctx *gin.Context, req *dtoorg.TenantDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}
	if err = organization.CheckOrganizationAccess(ctx, tenantEntity.OrganizationID); err != nil {
		return err
	}

	if err = iamdao.NewTenantDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoTenant Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

func (svc *tenantSvc) Update(ctx *gin.Context, req *dtoorg.TenantUpdateReq) error {
	tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantUpdate] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}
	if err = organization.CheckOrganizationAccess(ctx, tenantEntity.OrganizationID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"address":                    req.Address,
		"contact_email":              req.ContactEmail,
		"contact_phone":              req.ContactPhone,
		"legal_person":               req.LegalPerson,
		"logo":                       req.Logo,
		"organization_id":            req.OrganizationID,
		"short_name":                 req.ShortName,
		"status":                     req.Status,
		"tenant_code":                req.TenantCode,
		"tenant_name":                req.TenantName,
		"unified_social_credit_code": req.UnifiedSocialCreditCode,
	}
	if err = iamdao.NewTenantDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcorg.TenantUpdate] daoTenant UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

func (svc *tenantSvc) Detail(ctx *gin.Context, req *dtoorg.TenantDetailReq) (*dtoorg.TenantDetailResp, error) {
	tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantDetail] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}
	if err = organization.CheckOrganizationAccess(ctx, tenantEntity.OrganizationID); err != nil {
		return nil, err
	}
	resp := &dtoorg.TenantDetailResp{
		ID: tenantEntity.ID,
		TenantBaseInfo: objorg.TenantBaseInfo{
			Address:                 tenantEntity.Address,
			ContactEmail:            tenantEntity.ContactEmail,
			ContactPhone:            tenantEntity.ContactPhone,
			LegalPerson:             tenantEntity.LegalPerson,
			Logo:                    tenantEntity.Logo,
			OrganizationID:          tenantEntity.OrganizationID,
			ShortName:               tenantEntity.ShortName,
			Status:                  tenantEntity.Status,
			TenantCode:              tenantEntity.TenantCode,
			TenantName:              tenantEntity.TenantName,
			UnifiedSocialCreditCode: tenantEntity.UnifiedSocialCreditCode,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: tenantEntity.CreatedAt.Unix(),
			UpdatedAt: tenantEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *tenantSvc) PageList(ctx *gin.Context, req *dtoorg.TenantPageListReq) (*dtoorg.TenantPageListResp, error) {
	cond := &iamdao.TenantCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	tenantEntityList, total, err := iamdao.NewTenantDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantPageList] daoTenant GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtoorg.TenantPageListItem, 0, len(tenantEntityList))
	for _, v := range tenantEntityList {
		list = append(list, dtoorg.TenantPageListItem{
			ID: v.ID,
			TenantBaseInfo: objorg.TenantBaseInfo{
				Address:                 v.Address,
				ContactEmail:            v.ContactEmail,
				ContactPhone:            v.ContactPhone,
				LegalPerson:             v.LegalPerson,
				Logo:                    v.Logo,
				OrganizationID:          v.OrganizationID,
				ShortName:               v.ShortName,
				Status:                  v.Status,
				TenantCode:              v.TenantCode,
				TenantName:              v.TenantName,
				UnifiedSocialCreditCode: v.UnifiedSocialCreditCode,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoorg.TenantPageListResp{
		List:  list,
		Total: total,
	}, nil
}
