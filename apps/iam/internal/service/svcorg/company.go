package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/genericdao"
	"github.com/morehao/golib/biz/gcontext/gincontext"
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
		TenantID:                req.TenantID,
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

	if err := iamdao.NewCompanyDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoCompany Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyDeleteError)
	}
	return nil
}

// Update 更新公司管理
func (svc *companySvc) Update(ctx *gin.Context, req *dtoorg.CompanyUpdateReq) error {

	updateEntity := &iammodel.CompanyEntity{
		Address:                 req.Address,
		CompanyCode:             req.CompanyCode,
		CompanyName:             req.CompanyName,
		ContactEmail:            req.ContactEmail,
		ContactPhone:            req.ContactPhone,
		LegalPerson:             req.LegalPerson,
		Logo:                    req.Logo,
		ShortName:               req.ShortName,
		Status:                  req.Status,
		TenantID:                req.TenantID,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
	}
	if err := iamdao.NewCompanyDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyUpdate] daoCompany UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.CompanyUpdateError)
	}
	return nil
}

// Detail 根据id获取公司管理
func (svc *companySvc) Detail(ctx *gin.Context, req *dtoorg.CompanyDetailReq) (*dtoorg.CompanyDetailResp, error) {
	detailEntity, err := iamdao.NewCompanyDao().GetById(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyDetail] daoCompany GetById fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.CompanyGetDetailError)
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.CompanyNotExistError)
	}
	resp := &dtoorg.CompanyDetailResp{
		ID: detailEntity.ID,
		CompanyBaseInfo: objorg.CompanyBaseInfo{
			Address:                 detailEntity.Address,
			CompanyCode:             detailEntity.CompanyCode,
			CompanyName:             detailEntity.CompanyName,
			ContactEmail:            detailEntity.ContactEmail,
			ContactPhone:            detailEntity.ContactPhone,
			LegalPerson:             detailEntity.LegalPerson,
			Logo:                    detailEntity.Logo,
			ShortName:               detailEntity.ShortName,
			Status:                  detailEntity.Status,
			TenantID:                detailEntity.TenantID,
			UnifiedSocialCreditCode: detailEntity.UnifiedSocialCreditCode,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
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
	dataList, total, err := iamdao.NewCompanyDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.CompanyPageList] daoCompany GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.CompanyGetPageListError)
	}
	list := make([]dtoorg.CompanyPageListItem, 0, len(dataList))
	for _, v := range dataList {
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
