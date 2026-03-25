package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/tenant"
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

type CompanySvc interface {
	Create(ctx *gin.Context, req *dtoorg.CompanyCreateReq) (*dtoorg.CompanyCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorg.CompanyDeleteReq) error
	Update(ctx *gin.Context, req *dtoorg.CompanyUpdateReq) error
	Detail(ctx *gin.Context, req *dtoorg.CompanyDetailReq) (*dtoorg.CompanyDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.CompanyPageListReq) (*dtoorg.CompanyPageListResp, error)
}

type companySvc struct {
}

var _ CompanySvc = (*companySvc)(nil)

func NewCompanySvc() CompanySvc {
	return &companySvc{}
}

// Create 创建公司管理
func (svc *companySvc) Create(ctx *gin.Context, req *dtoorg.CompanyCreateReq) (*dtoorg.CompanyCreateResp, error) {
	insertEntity := &iammodel.CompanyEntity{
		Address:                 req.Address,
		CompanyCode:             req.CompanyCode,
		CompanyName:             req.CompanyName,
		ContactEmail:            req.ContactEmail,
		ContactPhone:            req.ContactPhone,
		LegalPerson:             req.LegalPerson,
		Logo:                    req.Logo,
		ShortName:               req.ShortName,
		Status:                  req.Status,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
	}

	if err := iamdao.NewCompanyDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyCreate] daoCompany Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.CompanyCreateError)
	}
	return &dtoorg.CompanyCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除公司管理
func (svc *companySvc) Delete(ctx *gin.Context, req *dtoorg.CompanyDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	companyEntity, err := iamdao.NewCompanyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoCompany GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyDeleteError)
	}
	if companyEntity == nil || companyEntity.ID == 0 {
		return code.GetError(code.CompanyNotExistError)
	}
	if err = tenant.CheckTenantAccess(ctx, companyEntity.TenantID); err != nil {
		return err
	}

	if err = iamdao.NewCompanyDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoCompany Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyDeleteError)
	}
	return nil
}

// Update 更新公司管理
func (svc *companySvc) Update(ctx *gin.Context, req *dtoorg.CompanyUpdateReq) error {
	companyEntity, err := iamdao.NewCompanyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyUpdate] daoCompany GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyUpdateError)
	}
	if companyEntity == nil || companyEntity.ID == 0 {
		return code.GetError(code.CompanyNotExistError)
	}
	if err = tenant.CheckTenantAccess(ctx, companyEntity.TenantID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"address":                    req.Address,
		"company_code":               req.CompanyCode,
		"company_name":               req.CompanyName,
		"contact_email":              req.ContactEmail,
		"contact_phone":              req.ContactPhone,
		"legal_person":               req.LegalPerson,
		"logo":                       req.Logo,
		"short_name":                 req.ShortName,
		"status":                     req.Status,
		"unified_social_credit_code": req.UnifiedSocialCreditCode,
	}
	if err = iamdao.NewCompanyDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyUpdate] daoCompany UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyUpdateError)
	}
	return nil
}

// Detail 根据id获取公司管理
func (svc *companySvc) Detail(ctx *gin.Context, req *dtoorg.CompanyDetailReq) (*dtoorg.CompanyDetailResp, error) {
	companyEntity, err := iamdao.NewCompanyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyDetail] daoCompany GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.CompanyGetDetailError)
	}
	// 判断是否存在
	if companyEntity == nil || companyEntity.ID == 0 {
		return nil, code.GetError(code.CompanyNotExistError)
	}
	if err = tenant.CheckTenantAccess(ctx, companyEntity.TenantID); err != nil {
		return nil, err
	}
	resp := &dtoorg.CompanyDetailResp{
		ID: companyEntity.ID,
		CompanyBaseInfo: objorg.CompanyBaseInfo{
			Address:                 companyEntity.Address,
			CompanyCode:             companyEntity.CompanyCode,
			CompanyName:             companyEntity.CompanyName,
			ContactEmail:            companyEntity.ContactEmail,
			ContactPhone:            companyEntity.ContactPhone,
			LegalPerson:             companyEntity.LegalPerson,
			Logo:                    companyEntity.Logo,
			ShortName:               companyEntity.ShortName,
			Status:                  companyEntity.Status,
			TenantID:                companyEntity.TenantID,
			UnifiedSocialCreditCode: companyEntity.UnifiedSocialCreditCode,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: companyEntity.CreatedAt.Unix(),
			UpdatedAt: companyEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取公司管理列表
func (svc *companySvc) PageList(ctx *gin.Context, req *dtoorg.CompanyPageListReq) (*dtoorg.CompanyPageListResp, error) {
	cond := &iamdao.CompanyCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	companyEntityList, total, err := iamdao.NewCompanyDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyPageList] daoCompany GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.CompanyGetPageListError)
	}
	list := make([]dtoorg.CompanyPageListItem, 0, len(companyEntityList))
	for _, v := range companyEntityList {
		list = append(list, dtoorg.CompanyPageListItem{
			ID: v.ID,
			CompanyBaseInfo: objorg.CompanyBaseInfo{
				Address:                 v.Address,
				CompanyCode:             v.CompanyCode,
				CompanyName:             v.CompanyName,
				ContactEmail:            v.ContactEmail,
				ContactPhone:            v.ContactPhone,
				LegalPerson:             v.LegalPerson,
				Logo:                    v.Logo,
				ShortName:               v.ShortName,
				Status:                  v.Status,
				TenantID:                v.TenantID,
				UnifiedSocialCreditCode: v.UnifiedSocialCreditCode,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoorg.CompanyPageListResp{
		List:  list,
		Total: total,
	}, nil
}
