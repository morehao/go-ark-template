package svcorg

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	org "github.com/morehao/goark/iam/core/org"
	"github.com/morehao/goark/iam/core/user"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/iam/object/objorg"
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
	Create(ctx *gin.Context, req *dtoorg.OrgCreateReq) (*dtoorg.OrgCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorg.OrgDeleteReq) error
	Update(ctx *gin.Context, req *dtoorg.OrgUpdateReq) error
	GetOrgConfig(ctx *gin.Context, req *dtoorg.GetOrganizationConfigsReq) (*dtoorg.GetOrgConfigResp, error)
	Detail(ctx *gin.Context, req *dtoorg.OrgDetailReq) (*dtoorg.OrgDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.OrgPageListReq) (*dtoorg.OrgPageListResp, error)
	ListConfigDefinitions(ctx *gin.Context) (*dtoorg.ListConfigDefinitionsResp, error)
}

type organizationSvc struct {
}

var _ OrganizationSvc = (*organizationSvc)(nil)

func NewOrganizationSvc() OrganizationSvc {
	return &organizationSvc{}
}

func (svc *organizationSvc) Create(ctx *gin.Context, req *dtoorg.OrgCreateReq) (*dtoorg.OrgCreateResp, error) {
	operatorID := gincontext.GetUserID(ctx)

	if req.DisplayCode == "" {
		glog.Errorf(ctx, "[svcorg.OrgCreate] displayCode is required, req:%s", gutil.ToJsonString(req))
		return nil, code.GetError(code.ErrInvalidParam)
	}
	if len(req.DisplayCode) > 32 {
		glog.Errorf(ctx, "[svcorg.OrgCreate] displayCode too long, req:%s", gutil.ToJsonString(req))
		return nil, code.GetError(code.ErrInvalidParam)
	}
	if req.OrgName == "" {
		glog.Errorf(ctx, "[svcorg.OrgCreate] orgName is required, req:%s", gutil.ToJsonString(req))
		return nil, code.GetError(code.ErrInvalidParam)
	}
	if len(req.OrgName) > 64 {
		glog.Errorf(ctx, "[svcorg.OrgCreate] orgName too long, req:%s", gutil.ToJsonString(req))
		return nil, code.GetError(code.ErrInvalidParam)
	}

	existOrg, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
		DisplayCode: req.DisplayCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgCreate] check displayCode fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OrganizationCreateError)
	}
	if existOrg != nil && existOrg.ID != 0 {
		glog.Errorf(ctx, "[svcorg.OrgCreate] displayCode already exists, displayCode:%s, req:%s", req.DisplayCode, gutil.ToJsonString(req))
		return nil, code.GetError(code.OrganizationCodeDuplicateError)
	}

	if len(req.AppIDs) > 0 {
		apps, err := dao.NewApplicationDao().GetListByCond(ctx, &dao.ApplicationCond{
			BaseCond: &genericdao.BaseCond{IDs: req.AppIDs},
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.OrgCreate] GetListByCond apps fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return nil, code.GetError(code.ApplicationGetDetailError)
		}
		appMap := apps.ToMap()
		if len(appMap) != len(req.AppIDs) {
			glog.Errorf(ctx, "[svcorg.OrgCreate] apps not all exist, appIDs:%v, req:%s", req.AppIDs, gutil.ToJsonString(req))
			return nil, code.GetError(code.ApplicationNotExistError)
		}
		for _, appID := range req.AppIDs {
			if app, ok := appMap[appID]; !ok || app.Status != string(model.AppStatusEnabled) {
				glog.Errorf(ctx, "[svcorg.OrgCreate] app invalid, appID:%d, req:%s", appID, gutil.ToJsonString(req))
				return nil, code.GetError(code.ApplicationInvalidError)
			}
		}
	}

	platformTenant, err := org.GetPlatformTenant(ctx)
	if err != nil || platformTenant == nil || platformTenant.ID == 0 {
		glog.Errorf(ctx, "[svcorg.OrgCreate] GetPlatformTenant fail, err:%v", err)
		return nil, code.GetError(code.TenantCreateError)
	}

	platformDept, err := org.GetPlatformDept(ctx, platformTenant.ID)
	if err != nil || platformDept == nil || platformDept.ID == 0 {
		glog.Errorf(ctx, "[svcorg.OrgCreate] GetPlatformDept fail, err:%v", err)
		return nil, code.GetError(code.TenantCreateError)
	}

	insertEntity := &model.OrganizationEntity{
		Code:        generateOrgCode(),
		DisplayCode: req.DisplayCode,
		OrgName:     req.OrgName,
		Description: req.Description,
		Logo:        req.Logo,
		Sequence:    req.Sequence,
		Status:      model.OrgStatus(req.Status),
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	var adminID uint
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewOrganizationDao().WithTx(tx).Insert(ctx, insertEntity); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgCreate] Insert org fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.OrganizationCreateError)
		}

		if len(req.AppIDs) > 0 {
			for _, appID := range req.AppIDs {
				orgAppEntity := &model.OrganizationApplicationEntity{
					OrgID:     insertEntity.ID,
					AppID:     appID,
					CreatedBy: operatorID,
				}
				if err := dao.NewOrganizationApplicationDao().WithTx(tx).Insert(ctx, orgAppEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.OrgCreate] Insert orgApp fail, err:%v, appID:%d", err, appID)
					return code.GetError(code.OrganizationCreateError)
				}
			}
		}

		if req.Admin != nil && req.Admin.Username != "" {
			if req.Admin.Mobile == "" && req.Admin.Email == "" {
				glog.Errorf(ctx, "[svcorg.OrgCreate] admin mobile or email is required, req:%s", gutil.ToJsonString(req))
				return code.GetError(code.ErrInvalidParam)
			}

			params := &user.CreatePersonParams{
				Mobile:     req.Admin.Mobile,
				Email:      req.Admin.Email,
				RealName:   req.Admin.RealName,
				Username:   req.Admin.Username,
				OperatorID: operatorID,
				TenantID:   platformTenant.ID,
				DeptID:     platformDept.ID,
				Status:     model.UserStatusEnabled,
				UserType:   model.UserTypeTenantAdmin,
			}
			result, err := user.CreatePersonWithUser(ctx, tx, params)
			if err != nil {
				glog.Errorf(ctx, "[svcorg.OrgCreate] CreatePersonWithUser fail, err:%v, req:%s", err, gutil.ToJsonString(req))
				return err
			}
			adminID = result.UserID
		}

		if len(req.Configs) > 0 {
			for _, cfg := range req.Configs {
				meta := model.GetOrgConfigMetaByKey(cfg.Key)
				if meta == nil {
					glog.Errorf(ctx, "[svcorg.OrgCreate] config key not found, key:%s, req:%s", cfg.Key, gutil.ToJsonString(req))
					return code.GetError(code.OrganizationConfigError)
				}
				configValue := cfg.Value
				if configValue == "" {
					configValue = meta.DefaultValue
				}
				if !meta.ValidateValue(configValue) {
					glog.Errorf(ctx, "[svcorg.OrgCreate] config value invalid, key:%s, value:%s, req:%s", cfg.Key, cfg.Value, gutil.ToJsonString(req))
					return code.GetError(code.OrganizationConfigError)
				}
				configEntity := &model.OrganizationConfigEntity{
					ConfigGroup: meta.Group,
					ConfigKey:   meta.Key,
					ValueType:   meta.Type,
					ConfigValue: configValue,
					Description: meta.Description,
					OrgID:       insertEntity.ID,
				}
				if err := dao.NewOrganizationConfigDao().WithTx(tx).Insert(ctx, configEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.OrgCreate] Insert config fail, err:%v, key:%s", err, cfg.Key)
					return err
				}
			}
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.OrgCreate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return nil, code.GetError(code.OrganizationCreateError)
	}
	return &dtoorg.OrgCreateResp{
		OrgID:   insertEntity.ID,
		AdminID: adminID,
	}, nil
}

func generateOrgCode() string {
	return uuid.New().String()
}

func (svc *organizationSvc) Delete(ctx *gin.Context, req *dtoorg.OrgDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, req.OrgID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDelete] daoOrg GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.OrganizationDeleteError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return code.GetError(code.OrganizationNotExistError)
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationApplicationEntity{}).Error; err != nil {
			glog.Errorf(ctx, "[svcorg.OrgDelete] Delete orgApps fail, err:%v, orgID:%d", err, req.OrgID)
			return err
		}

		if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationConfigEntity{}).Error; err != nil {
			glog.Errorf(ctx, "[svcorg.OrgDelete] Delete orgConfigs fail, err:%v, orgID:%d", err, req.OrgID)
			return err
		}

		tenants, _, err := dao.NewTenantDao().WithTx(tx).GetPageListByCond(ctx, &dao.TenantCond{
			OrgID: req.OrgID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.OrgDelete] GetPageListByCond tenants fail, err:%v, orgID:%d", err, req.OrgID)
			return err
		}
		if len(tenants) > 0 {
			tenantIDs := make([]uint, 0, len(tenants))
			for _, tenant := range tenants {
				tenantIDs = append(tenantIDs, tenant.ID)
				if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&model.TenantApplicationEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgDelete] Delete tenantApps fail, err:%v, tenantID:%d", err, tenant.ID)
					return err
				}
				if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&model.DepartmentEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgDelete] Delete departments fail, err:%v, tenantID:%d", err, tenant.ID)
					return err
				}
				if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&model.UserEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgDelete] Delete users fail, err:%v, tenantID:%d", err, tenant.ID)
					return err
				}
				if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&model.PersonEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgDelete] Delete persons fail, err:%v, tenantID:%d", err, tenant.ID)
					return err
				}
			}
			if err := tx.Where("tenant_id IN ?", tenantIDs).Delete(&model.TenantEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.OrgDelete] Delete tenants fail, err:%v, orgID:%d", err, req.OrgID)
				return err
			}
		}

		if err := dao.NewOrganizationDao().WithTx(tx).Delete(ctx, req.OrgID, userID); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgDelete] daoOrg Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.OrgDelete] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return code.GetError(code.OrganizationDeleteError)
	}
	return nil
}

func (svc *organizationSvc) Update(ctx *gin.Context, req *dtoorg.OrgUpdateReq) error {
	if req.DisplayCode != "" {
		if len(req.DisplayCode) > 32 {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] displayCode too long, req:%s", gutil.ToJsonString(req))
			return code.GetError(code.ErrInvalidParam)
		}
		existOrg, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
			DisplayCode: req.DisplayCode,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] check displayCode fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.OrganizationUpdateError)
		}
		if existOrg != nil && existOrg.ID != 0 && existOrg.ID != req.OrgID {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] displayCode already exists, displayCode:%s, req:%s", req.DisplayCode, gutil.ToJsonString(req))
			return code.GetError(code.OrganizationCodeDuplicateError)
		}
	}

	if len(req.AppIDs) > 0 {
		apps, err := dao.NewApplicationDao().GetListByCond(ctx, &dao.ApplicationCond{
			BaseCond: &genericdao.BaseCond{IDs: req.AppIDs},
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] GetListByCond apps fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.ApplicationGetDetailError)
		}
		appMap := apps.ToMap()
		if len(appMap) != len(req.AppIDs) {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] apps not all exist, appIDs:%v, req:%s", req.AppIDs, gutil.ToJsonString(req))
			return code.GetError(code.ApplicationNotExistError)
		}
		for _, appID := range req.AppIDs {
			if app, ok := appMap[appID]; !ok || app.Status != string(model.AppStatusEnabled) {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] app invalid, appID:%d, req:%s", appID, gutil.ToJsonString(req))
				return code.GetError(code.ApplicationInvalidError)
			}
		}
	}

	updateMap := map[string]any{}
	if req.DisplayCode != "" {
		updateMap["display_code"] = req.DisplayCode
	}
	if req.OrgName != "" {
		updateMap["org_name"] = req.OrgName
	}
	updateMap["description"] = req.Description
	updateMap["logo"] = req.Logo
	updateMap["sequence"] = req.Sequence
	updateMap["status"] = req.Status

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewOrganizationDao().WithTx(tx).UpdateMap(ctx, req.OrgID, updateMap); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] daoOrg UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.OrganizationUpdateError)
		}

		if len(req.AppIDs) > 0 {
			if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationApplicationEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete orgApps fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.OrganizationUpdateError)
			}
			operatorID := gincontext.GetUserID(ctx)
			for _, appID := range req.AppIDs {
				orgAppEntity := &model.OrganizationApplicationEntity{
					OrgID:     req.OrgID,
					AppID:     appID,
					CreatedBy: operatorID,
				}
				if err := dao.NewOrganizationApplicationDao().WithTx(tx).Insert(ctx, orgAppEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] Insert orgApp fail, err:%v, appID:%d", err, appID)
					return code.GetError(code.OrganizationUpdateError)
				}
			}

			tenants, _, err := dao.NewTenantDao().WithTx(tx).GetPageListByCond(ctx, &dao.TenantCond{
				OrgID: req.OrgID,
			})
			if err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] GetPageListByCond tenants fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.OrganizationUpdateError)
			}
			if len(tenants) > 0 {
				tenantIDs := make([]uint, 0, len(tenants))
				for _, tenant := range tenants {
					tenantIDs = append(tenantIDs, tenant.ID)
				}
				if err := tx.Where("tenant_id IN ? AND app_id NOT IN ?", tenantIDs, req.AppIDs).Delete(&model.TenantApplicationEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete tenantApps fail, err:%v, orgID:%d", err, req.OrgID)
					return code.GetError(code.OrganizationUpdateError)
				}
			}
		}

		if len(req.Configs) > 0 {
			if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationConfigEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete configs fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.OrganizationUpdateError)
			}

			for _, cfg := range req.Configs {
				meta := model.GetOrgConfigMetaByKey(cfg.Key)
				if meta == nil {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] config key not found, key:%s, req:%s", cfg.Key, gutil.ToJsonString(req))
					return code.GetError(code.OrganizationConfigError)
				}
				configValue := cfg.Value
				if configValue == "" {
					configValue = meta.DefaultValue
				}
				if !meta.ValidateValue(configValue) {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] config value invalid, key:%s, value:%s, req:%s", cfg.Key, cfg.Value, gutil.ToJsonString(req))
					return code.GetError(code.OrganizationConfigError)
				}
				configEntity := &model.OrganizationConfigEntity{
					ConfigGroup: meta.Group,
					ConfigKey:   meta.Key,
					ValueType:   meta.Type,
					ConfigValue: configValue,
					Description: meta.Description,
					OrgID:       req.OrgID,
				}
				if err := dao.NewOrganizationConfigDao().WithTx(tx).Insert(ctx, configEntity); err != nil {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] Insert config fail, err:%v, key:%s", err, cfg.Key)
					return err
				}
			}
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcorg.OrgUpdate] Transaction fail, err:%v, req:%s", txErr, gutil.ToJsonString(req))
		return code.GetError(code.OrganizationUpdateError)
	}
	return nil
}

func (svc *organizationSvc) GetOrgConfig(ctx *gin.Context, req *dtoorg.GetOrganizationConfigsReq) (*dtoorg.GetOrgConfigResp, error) {
	displayCode := resolveDisplayCode(ctx, req.DisplayCode)
	if displayCode == "" {
		return nil, code.GetError(code.OrganizationNotExistError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
		DisplayCode: displayCode,
		Status:      model.OrgStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetOrgConfig] daoOrg GetByCond fail, err:%v, displayCode:%s", err, displayCode)
		return nil, code.GetError(code.OrganizationGetDetailError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.OrganizationNotExistError)
	}

	configEntityList, err := dao.NewOrganizationConfigDao().GetListByCond(ctx, &dao.OrganizationConfigCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetOrgConfig] daoOrgConfig GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.OrganizationGetDetailError)
	}

	groupMap := make(map[string][]dtoorg.ConfigItemResp)
	for _, v := range configEntityList {
		groupMap[v.ConfigGroup] = append(groupMap[v.ConfigGroup], dtoorg.ConfigItemResp{
			Key:   v.ConfigKey,
			Value: v.ConfigValue,
		})
	}

	configGroups := make([]dtoorg.ConfigGroupResp, 0, len(groupMap))
	for group, items := range groupMap {
		configGroups = append(configGroups, dtoorg.ConfigGroupResp{
			Group: group,
			Items: items,
		})
	}

	return &dtoorg.GetOrgConfigResp{
		OrgID:   orgEntity.ID,
		OrgName: orgEntity.OrgName,
		Logo:    orgEntity.Logo,
		Status:  orgEntity.Status,
		Configs: configGroups,
	}, nil
}

func resolveDisplayCode(ctx *gin.Context, reqDisplayCode string) string {
	host := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Host"))
	if host == "" && ctx.Request != nil {
		host = strings.TrimSpace(ctx.Request.Host)
	}
	if host == "" {
		host = strings.TrimSpace(reqDisplayCode)
	}
	host = strings.TrimSpace(strings.Split(host, ",")[0])
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.Split(host, "/")[0]

	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	} else if strings.Count(host, ":") == 1 {
		host = strings.Split(host, ":")[0]
	}

	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	return host
}

func (svc *organizationSvc) Detail(ctx *gin.Context, req *dtoorg.OrgDetailReq) (*dtoorg.OrgDetailResp, error) {
	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, req.OrgID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDetail] daoOrg GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OrganizationGetDetailError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.OrganizationNotExistError)
	}

	configEntityList, err := dao.NewOrganizationConfigDao().GetListByCond(ctx, &dao.OrganizationConfigCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDetail] daoOrgConfig GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.OrganizationGetDetailError)
	}

	configs := make(map[string]map[string]string)
	for _, v := range configEntityList {
		if _, ok := configs[v.ConfigGroup]; !ok {
			configs[v.ConfigGroup] = make(map[string]string)
		}
		configs[v.ConfigGroup][v.ConfigKey] = v.ConfigValue
	}

	orgAppList, err := dao.NewOrganizationApplicationDao().GetListByCond(ctx, &dao.OrganizationApplicationCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDetail] daoOrgApp GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.OrganizationGetDetailError)
	}

	var apps []dtoorg.AppInfo
	if len(orgAppList) > 0 {
		appIDs := make([]uint, 0, len(orgAppList))
		for _, orgApp := range orgAppList {
			appIDs = append(appIDs, orgApp.AppID)
		}
		appEntities, err := dao.NewApplicationDao().GetListByCond(ctx, &dao.ApplicationCond{
			BaseCond: &genericdao.BaseCond{IDs: appIDs},
		})
		if err != nil {
			glog.Errorf(ctx, "[svcorg.OrgDetail] GetListByCond apps fail, err:%v, orgID:%d", err, orgEntity.ID)
			return nil, code.GetError(code.OrganizationGetDetailError)
		}
		appMap := appEntities.ToMap()
		for _, orgApp := range orgAppList {
			if app, ok := appMap[orgApp.AppID]; ok {
				apps = append(apps, dtoorg.AppInfo{
					AppID:   app.ID,
					AppName: app.AppName,
				})
			}
		}
	}

	resp := &dtoorg.OrgDetailResp{
		OrgID:   orgEntity.ID,
		Configs: configs,
		Apps:    apps,
		OrgBaseInfo: objorg.OrgBaseInfo{
			DisplayCode: orgEntity.DisplayCode,
			OrgName:     orgEntity.OrgName,
			Description: orgEntity.Description,
			Logo:        orgEntity.Logo,
			Sequence:    orgEntity.Sequence,
			Status:      string(orgEntity.Status),
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: orgEntity.CreatedAt.Unix(),
			UpdatedAt: orgEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *organizationSvc) PageList(ctx *gin.Context, req *dtoorg.OrgPageListReq) (*dtoorg.OrgPageListResp, error) {
	cond := &dao.OrganizationCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
	}
	orgEntityList, total, err := dao.NewOrganizationDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgPageList] daoOrg GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.OrganizationGetPageListError)
	}
	list := make([]dtoorg.OrgPageListItem, 0, len(orgEntityList))
	for _, v := range orgEntityList {
		list = append(list, dtoorg.OrgPageListItem{
			OrgID: v.ID,
			OrgBaseInfo: objorg.OrgBaseInfo{
				DisplayCode: v.DisplayCode,
				OrgName:     v.OrgName,
				Description: v.Description,
				Logo:        v.Logo,
				Sequence:    v.Sequence,
				Status:      string(v.Status),
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoorg.OrgPageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *organizationSvc) ListConfigDefinitions(ctx *gin.Context) (*dtoorg.ListConfigDefinitionsResp, error) {
	configs := make([]dtoorg.OrgConfigMetaResp, 0, len(model.OrgConfigMetaList))
	for _, meta := range model.OrgConfigMetaList {
		options := make([]dtoorg.OrgConfigOptionsItem, 0, len(meta.Options))
		for _, opt := range meta.Options {
			options = append(options, dtoorg.OrgConfigOptionsItem{
				Value:       opt.Value,
				Description: opt.Description,
			})
		}
		configs = append(configs, dtoorg.OrgConfigMetaResp{
			Group:        meta.Group,
			Key:          meta.Key,
			Type:         meta.Type,
			DefaultValue: meta.DefaultValue,
			Description:  meta.Description,
			Options:      options,
		})
	}
	return &dtoorg.ListConfigDefinitionsResp{Configs: configs}, nil
}
