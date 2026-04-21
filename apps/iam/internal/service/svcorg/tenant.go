package svcorg

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
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
	operatorID := gincontext.GetUserID(ctx)
	insertEntity := &iammodel.TenantEntity{
		Address:                 req.Address,
		ContactEmail:            req.ContactEmail,
		ContactPhone:            req.ContactPhone,
		LegalPerson:             req.LegalPerson,
		Logo:                    req.Logo,
		OrgID:                   req.OrgID,
		ShortName:               req.ShortName,
		Status:                  iammodel.TenantStatus(req.Status),
		TenantCode:              req.TenantCode,
		TenantName:              req.TenantName,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
		CreatedBy:               operatorID,
		UpdatedBy:               operatorID,
	}

	adminInfo := req.AdminInfo
	mobile := strings.TrimSpace(adminInfo.Mobile)
	email := strings.TrimSpace(adminInfo.Email)

	var tenantID uint
	var deptID uint
	var result *user.CreatePersonResult
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := iamdao.NewTenantDao().WithTx(tx).Insert(ctx, insertEntity); err != nil {
			return err
		}
		tenantID = insertEntity.ID

		deptEntity := &iammodel.DepartmentEntity{
			TenantID:  tenantID,
			DeptCode:  req.TenantCode,
			DeptName:  req.TenantName,
			DeptLevel: 1,
			ParentID:  0,
			Status:    iammodel.DeptStatusEnabled,
			CreatedBy: operatorID,
			UpdatedBy: operatorID,
		}
		if err := iamdao.NewDepartmentDao().WithTx(tx).Insert(ctx, deptEntity); err != nil {
			return err
		}
		deptID = deptEntity.ID

		deptPathMap := map[string]any{
			"dept_path": fmt.Sprintf("/%d/", deptEntity.ID),
		}
		if err := iamdao.NewDepartmentDao().WithTx(tx).UpdateMap(ctx, deptEntity.ID, deptPathMap); err != nil {
			return err
		}

		params := &user.CreatePersonParams{
			Mobile:     mobile,
			Email:      email,
			RealName:   adminInfo.RealName,
			OperatorID: operatorID,
			TenantID:   tenantID,
			DeptID:     deptID,
			Username:   adminInfo.Username,
			UserType:   iammodel.UserTypeTenantAdmin,
			Status:     iammodel.UserStatusEnabled,
		}
		var err error
		result, err = user.CreatePersonWithUser(ctx, tx, params)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.TenantCreate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}

	return &dtoorg.TenantCreateResp{
		ID:       tenantID,
		AdminID:  result.UserID,
		PersonID: result.PersonID,
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
	updateMap := map[string]any{
		"address":                    req.Address,
		"contact_email":              req.ContactEmail,
		"contact_phone":              req.ContactPhone,
		"legal_person":               req.LegalPerson,
		"logo":                       req.Logo,
		"org_id":                     req.OrgID,
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
	resp := &dtoorg.TenantDetailResp{
		ID: tenantEntity.ID,
		TenantBaseInfo: objorg.TenantBaseInfo{
			Address:                 tenantEntity.Address,
			ContactEmail:            tenantEntity.ContactEmail,
			ContactPhone:            tenantEntity.ContactPhone,
			LegalPerson:             tenantEntity.LegalPerson,
			Logo:                    tenantEntity.Logo,
			OrgID:                   tenantEntity.OrgID,
			ShortName:               tenantEntity.ShortName,
			Status:                  string(tenantEntity.Status),
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
	tenantID := gincontext.GetTenantID(ctx)

	isPlatformAdmin := gincontext.GetUserType(ctx) == "platform_admin"

	cond := &iamdao.TenantCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}

	if !isPlatformAdmin && tenantID > 0 {
		cond.ID = tenantID
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
				OrgID:                   v.OrgID,
				ShortName:               v.ShortName,
				Status:                  string(v.Status),
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