package svcorg

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/user"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/model"
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

	if len(req.AppIDs) > 0 {
		orgApps, err := dao.NewOrganizationApplicationDao().GetListByCond(ctx, &dao.OrganizationApplicationCond{
			OrgID: req.OrgID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.TenantCreate] GetListByCond orgApps fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return nil, code.GetError(code.TenantCreateError)
		}
		orgAppMap := make(map[uint]bool)
		for _, orgApp := range orgApps {
			orgAppMap[orgApp.AppID] = true
		}
		for _, appID := range req.AppIDs {
			if !orgAppMap[appID] {
				glog.Errorf(ctx, "[svcorg.TenantCreate] app not in org scope, appID:%d, req:%s", appID, gutil.ToJsonString(req))
				return nil, code.GetError(code.ApplicationInvalidError)
			}
		}
	}

	var tenantLevel int32 = 1
	var tenantPath string = "/"

	if req.ParentID > 0 {
		parentTenant, err := dao.NewTenantDao().GetByID(ctx, req.ParentID)
		if err != nil || parentTenant == nil || parentTenant.ID == 0 {
			glog.Errorf(ctx, "[svcorg.TenantCreate] parent tenant not found, parentID:%d", req.ParentID)
			return nil, code.GetError(code.TenantNotExistError)
		}
		if parentTenant.OrgID != req.OrgID {
			return nil, code.GetError(code.TenantScopeForbiddenError)
		}
		tenantLevel = parentTenant.TenantLevel + 1
		tenantPath = fmt.Sprintf("%s%d/", parentTenant.TenantPath, parentTenant.ID)
	} else {
		tenantPath = fmt.Sprintf("/%d/", 0)
	}

	insertEntity := &model.TenantEntity{
		Address:                 req.Address,
		ContactEmail:            req.ContactEmail,
		ContactPhone:            req.ContactPhone,
		LegalPerson:             req.LegalPerson,
		Logo:                    req.Logo,
		OrgID:                   req.OrgID,
		ParentID:                req.ParentID,
		ShortName:               req.ShortName,
		Status:                  model.TenantStatus(req.Status),
		TenantCode:              req.TenantCode,
		TenantLevel:             tenantLevel,
		TenantName:              req.TenantName,
		TenantPath:              tenantPath,
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
		if err := dao.NewTenantDao().WithTx(tx).Insert(ctx, insertEntity); err != nil {
			return err
		}
		tenantID = insertEntity.ID

		updatePathMap := map[string]any{
			"tenant_path": fmt.Sprintf("/%d/", insertEntity.ID),
		}
		if err := dao.NewTenantDao().WithTx(tx).UpdateMap(ctx, insertEntity.ID, updatePathMap); err != nil {
			return err
		}

		deptEntity := &model.DepartmentEntity{
			TenantID:  tenantID,
			DeptCode:  req.TenantCode,
			DeptName:  req.TenantName,
			DeptLevel: 1,
			ParentID:  0,
			Status:    model.DeptStatusEnabled,
			CreatedBy: operatorID,
			UpdatedBy: operatorID,
		}
		if err := dao.NewDepartmentDao().WithTx(tx).Insert(ctx, deptEntity); err != nil {
			return err
		}
		deptID = deptEntity.ID

		deptPathMap := map[string]any{
			"dept_path": fmt.Sprintf("/%d/", deptEntity.ID),
		}
		if err := dao.NewDepartmentDao().WithTx(tx).UpdateMap(ctx, deptEntity.ID, deptPathMap); err != nil {
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
			UserType:   model.UserTypeTenantAdmin,
			Status:     model.UserStatusEnabled,
		}
		var err error
		result, err = user.CreatePersonWithUser(ctx, tx, params)
		if err != nil {
			return err
		}

		if len(req.AppIDs) > 0 {
			for _, appID := range req.AppIDs {
				tenantAppEntity := &model.TenantApplicationEntity{
					TenantID:  tenantID,
					AppID:     appID,
					CreatedBy: operatorID,
				}
				if err := dao.NewTenantApplicationDao().WithTx(tx).Insert(ctx, tenantAppEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.TenantCreate] Insert tenantApp fail, err:%v, appID:%d", err, appID)
					return err
				}
			}
		}

		return nil
	})

	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.TenantCreate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantCreateError)
	}

	return &dtoorg.TenantCreateResp{
		TenantID: tenantID,
		AdminID:  result.UserID,
		PersonID: result.PersonID,
	}, nil
}

func (svc *tenantSvc) Delete(ctx *gin.Context, req *dtoorg.TenantDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	tenantEntity, err := dao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}

	if err = dao.NewTenantDao().Delete(ctx, req.TenantID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoTenant Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

func (svc *tenantSvc) Update(ctx *gin.Context, req *dtoorg.TenantUpdateReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantEntity, err := dao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantUpdate] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}

	if len(req.AppIDs) > 0 {
		orgApps, err := dao.NewOrganizationApplicationDao().GetListByCond(ctx, &dao.OrganizationApplicationCond{
			OrgID: tenantEntity.OrgID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.TenantUpdate] GetListByCond orgApps fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.TenantUpdateError)
		}
		orgAppMap := make(map[uint]bool)
		for _, orgApp := range orgApps {
			orgAppMap[orgApp.AppID] = true
		}
		for _, appID := range req.AppIDs {
			if !orgAppMap[appID] {
				glog.Errorf(ctx, "[svcorg.TenantUpdate] app not in org scope, appID:%d, req:%s", appID, gutil.ToJsonString(req))
				return code.GetError(code.ApplicationInvalidError)
			}
		}
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
		"updated_by":                 operatorID,
	}

	if req.ParentID != tenantEntity.ParentID {
		if req.ParentID > 0 {
			parentTenant, err := dao.NewTenantDao().GetByID(ctx, req.ParentID)
			if err != nil || parentTenant == nil || parentTenant.ID == 0 {
				glog.Errorf(ctx, "[svcorg.TenantUpdate] parent tenant not found, parentID:%d", req.ParentID)
				return code.GetError(code.TenantNotExistError)
			}
			if parentTenant.OrgID != req.OrgID {
				return code.GetError(code.TenantScopeForbiddenError)
			}
			updateMap["parent_id"] = req.ParentID
			updateMap["tenant_level"] = parentTenant.TenantLevel + 1
			updateMap["tenant_path"] = fmt.Sprintf("%s%d/", parentTenant.TenantPath, parentTenant.ID)
		} else {
			updateMap["parent_id"] = 0
			updateMap["tenant_level"] = 1
			updateMap["tenant_path"] = fmt.Sprintf("/%d/", 0)
		}
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewTenantDao().WithTx(tx).UpdateMap(ctx, req.TenantID, updateMap); err != nil {
			return err
		}

		if req.ParentID != tenantEntity.ParentID {
			if err := svc.updateChildrenTenantPath(ctx, tx, req.TenantID, updateMap["tenant_path"].(string), updateMap["tenant_level"].(int32)); err != nil {
				return err
			}
		}

		if len(req.AppIDs) > 0 {
			if err := tx.Where("tenant_id = ?", req.TenantID).Delete(&model.TenantApplicationEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.TenantUpdate] Delete tenantApps fail, err:%v, tenantID:%d", err, req.TenantID)
				return err
			}
			for _, appID := range req.AppIDs {
				tenantAppEntity := &model.TenantApplicationEntity{
					TenantID:  req.TenantID,
					AppID:     appID,
					CreatedBy: operatorID,
				}
				if err := dao.NewTenantApplicationDao().WithTx(tx).Insert(ctx, tenantAppEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.TenantUpdate] Insert tenantApp fail, err:%v, appID:%d", err, appID)
					return err
				}
			}
		}

		return nil
	})

	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.TenantUpdate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

func (svc *tenantSvc) updateChildrenTenantPath(ctx *gin.Context, tx *gorm.DB, parentID uint, parentPath string, parentLevel int32) error {
	childrenCond := &dao.TenantCond{
		ParentID: parentID,
	}
	children, _, err := dao.NewTenantDao().WithTx(tx).GetPageListByCond(ctx, childrenCond)
	if err != nil {
		return err
	}

	for _, child := range children {
		newPath := fmt.Sprintf("%s%d/", parentPath, child.ID)
		newLevel := parentLevel + 1
		childUpdateMap := map[string]any{
			"tenant_path":  newPath,
			"tenant_level": newLevel,
			"updated_by":   gincontext.GetUserID(ctx),
		}
		if err := dao.NewTenantDao().WithTx(tx).UpdateMap(ctx, child.ID, childUpdateMap); err != nil {
			return err
		}
		if err := svc.updateChildrenTenantPath(ctx, tx, child.ID, newPath, newLevel); err != nil {
			return err
		}
	}
	return nil
}

func (svc *tenantSvc) Detail(ctx *gin.Context, req *dtoorg.TenantDetailReq) (*dtoorg.TenantDetailResp, error) {
	tenantEntity, err := dao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantDetail] daoTenant GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}
	var apps []dtoorg.AppInfo
	tenantAppList, err := dao.NewTenantApplicationDao().GetListByCond(ctx, &dao.TenantApplicationCond{
		TenantID: tenantEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantDetail] GetListByCond tenantApps fail, err:%v, tenantID:%d", err, tenantEntity.ID)
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if len(tenantAppList) > 0 {
		appIDs := make([]uint, 0, len(tenantAppList))
		for _, tenantApp := range tenantAppList {
			appIDs = append(appIDs, tenantApp.AppID)
		}
		appEntities, err := dao.NewApplicationDao().GetListByCond(ctx, &dao.ApplicationCond{
			BaseCond: &genericdao.BaseCond{IDs: appIDs},
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.TenantDetail] GetListByCond apps fail, err:%v, tenantID:%d", err, tenantEntity.ID)
			return nil, code.GetError(code.TenantGetDetailError)
		}
		appMap := appEntities.ToMap()
		for _, tenantApp := range tenantAppList {
			if app, ok := appMap[tenantApp.AppID]; ok {
				apps = append(apps, dtoorg.AppInfo{
					AppID:   app.ID,
					AppName: app.AppName,
				})
			}
		}
	}

	resp := &dtoorg.TenantDetailResp{
		TenantID: tenantEntity.ID,
		Apps:     apps,
		TenantBaseInfo: objorg.TenantBaseInfo{
			Address:                 tenantEntity.Address,
			ContactEmail:            tenantEntity.ContactEmail,
			ContactPhone:            tenantEntity.ContactPhone,
			LegalPerson:             tenantEntity.LegalPerson,
			Logo:                    tenantEntity.Logo,
			OrgID:                   tenantEntity.OrgID,
			ParentID:                tenantEntity.ParentID,
			ShortName:               tenantEntity.ShortName,
			Status:                  string(tenantEntity.Status),
			TenantCode:              tenantEntity.TenantCode,
			TenantLevel:             tenantEntity.TenantLevel,
			TenantName:              tenantEntity.TenantName,
			TenantPath:              tenantEntity.TenantPath,
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

	cond := &dao.TenantCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}

	if !isPlatformAdmin && tenantID > 0 {
		cond.ID = tenantID
	}

	tenantEntityList, total, err := dao.NewTenantDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.TenantPageList] daoTenant GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtoorg.TenantPageListItem, 0, len(tenantEntityList))
	for _, v := range tenantEntityList {
		list = append(list, dtoorg.TenantPageListItem{
			TenantID: v.ID,
			TenantBaseInfo: objorg.TenantBaseInfo{
				Address:                 v.Address,
				ContactEmail:            v.ContactEmail,
				ContactPhone:            v.ContactPhone,
				LegalPerson:             v.LegalPerson,
				Logo:                    v.Logo,
				OrgID:                   v.OrgID,
				ParentID:                v.ParentID,
				ShortName:               v.ShortName,
				Status:                  string(v.Status),
				TenantCode:              v.TenantCode,
				TenantLevel:             v.TenantLevel,
				TenantName:              v.TenantName,
				TenantPath:              v.TenantPath,
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