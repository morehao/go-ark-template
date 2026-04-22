package svcorg

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	org "github.com/morehao/goark/apps/iam/core/org"
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

	insertEntity := &model.OrganizationEntity{
		Domain:      req.Domain,
		Logo:        req.Logo,
		Description: req.Description,
		Sequence:   req.Sequence,
		Status:      model.OrgStatus(req.Status),
		OrgCode:     req.OrgCode,
		OrgName:     req.OrgName,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	var adminID uint
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewOrganizationDao().WithTx(tx).Insert(ctx, insertEntity); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgCreate] Insert org fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.TenantCreateError)
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
					return code.GetError(code.CompanyCreateError)
				}
			}
		}

		if req.Admin != nil && req.Admin.Username != "" {
			platformTenant, err := org.GetPlatformTenant(ctx)
			if err != nil || platformTenant == nil || platformTenant.ID == 0 {
				glog.Errorf(ctx, "[svcorg.OrgCreate] GetPlatformTenant fail, err:%v", err)
				return code.GetError(code.TenantCreateError)
			}

			platformDept, err := org.GetPlatformDept(ctx, platformTenant.ID)
			if err != nil || platformDept == nil || platformDept.ID == 0 {
				glog.Errorf(ctx, "[svcorg.OrgCreate] GetPlatformDept fail, err:%v", err)
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
					continue
				}
				configValue := cfg.Value
				if configValue == "" {
					configValue = meta.DefaultValue
				}
				if !meta.ValidateValue(configValue) {
					continue
				}
				configEntity := &model.OrganizationConfigEntity{
					ConfigGroup: meta.Group,
					ConfigKey:   meta.Key,
					ValueType:  meta.Type,
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
		return nil, code.GetError(code.TenantCreateError)
	}
	return &dtoorg.OrgCreateResp{
		OrgID:   insertEntity.ID,
		AdminID: adminID,
	}, nil
}

func (svc *organizationSvc) Delete(ctx *gin.Context, req *dtoorg.OrgDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, req.OrgID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDelete] daoOrg GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return code.GetError(code.TenantNotExistError)
	}

	if err = dao.NewOrganizationDao().Delete(ctx, req.OrgID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDelete] daoOrg Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.TenantDeleteError)
	}
	return nil
}

func (svc *organizationSvc) Update(ctx *gin.Context, req *dtoorg.OrgUpdateReq) error {
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

	updateMap := map[string]any{
		"domain":      req.Domain,
		"logo":        req.Logo,
		"description": req.Description,
		"sequence":  req.Sequence,
		"status":      req.Status,
		"org_code":    req.OrgCode,
		"org_name":    req.OrgName,
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewOrganizationDao().WithTx(tx).UpdateMap(ctx, req.OrgID, updateMap); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] daoOrg UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.TenantUpdateError)
		}

		if len(req.AppIDs) > 0 {
			if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationApplicationEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete orgApps fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.TenantUpdateError)
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
					return code.GetError(code.TenantUpdateError)
				}
			}

			tenants, _, err := dao.NewTenantDao().WithTx(tx).GetPageListByCond(ctx, &dao.TenantCond{
				OrgID: req.OrgID,
			})
			if err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] GetPageListByCond tenants fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.TenantUpdateError)
			}
			if len(tenants) > 0 {
				tenantIDs := make([]uint, 0, len(tenants))
				for _, tenant := range tenants {
					tenantIDs = append(tenantIDs, tenant.ID)
				}
				if err := tx.Where("tenant_id IN ? AND app_id NOT IN ?", tenantIDs, req.AppIDs).Delete(&model.TenantApplicationEntity{}).Error; err != nil {
					glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete tenantApps fail, err:%v, orgID:%d", err, req.OrgID)
					return code.GetError(code.TenantUpdateError)
				}
			}
		}

		if len(req.Configs) > 0 {
			if err := tx.Where("org_id = ?", req.OrgID).Delete(&model.OrganizationConfigEntity{}).Error; err != nil {
				glog.Errorf(ctx, "[svcorg.OrgUpdate] Delete configs fail, err:%v, orgID:%d", err, req.OrgID)
				return code.GetError(code.TenantUpdateError)
			}

			for _, cfg := range req.Configs {
				meta := model.GetOrgConfigMetaByKey(cfg.Key)
				if meta == nil {
					continue
				}
				configValue := cfg.Value
				if configValue == "" {
					configValue = meta.DefaultValue
				}
				if !meta.ValidateValue(configValue) {
					continue
				}
				configEntity := &model.OrganizationConfigEntity{
					ConfigGroup: meta.Group,
					ConfigKey:   meta.Key,
					ValueType:  meta.Type,
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
		return code.GetError(code.TenantUpdateError)
	}
	return nil
}

func (svc *organizationSvc) GetOrgConfig(ctx *gin.Context, req *dtoorg.GetOrganizationConfigsReq) (*dtoorg.GetOrgConfigResp, error) {
	domain := resolveDomain(ctx, req.Domain)
	if domain == "" {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
		Domain: domain,
		Status: model.OrgStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetOrgConfig] daoOrg GetByCond fail, err:%v, domain:%s", err, domain)
		return nil, code.GetError(code.AuthLoginError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	configEntityList, err := dao.NewOrganizationConfigDao().GetListByCond(ctx, &dao.OrganizationConfigCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetOrgConfig] daoOrgConfig GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
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
		Domain:  orgEntity.Domain,
		Logo:    orgEntity.Logo,
		Status:  orgEntity.Status,
		Configs: configGroups,
	}, nil
}

func (svc *organizationSvc) Detail(ctx *gin.Context, req *dtoorg.OrgDetailReq) (*dtoorg.OrgDetailResp, error) {
	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, req.OrgID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDetail] daoOrg GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.TenantGetDetailError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.TenantNotExistError)
	}

	configEntityList, err := dao.NewOrganizationConfigDao().GetListByCond(ctx, &dao.OrganizationConfigCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.OrgDetail] daoOrgConfig GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.TenantGetDetailError)
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
		return nil, code.GetError(code.TenantGetDetailError)
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
			return nil, code.GetError(code.TenantGetDetailError)
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
			Domain:      orgEntity.Domain,
			Logo:        orgEntity.Logo,
			Description: orgEntity.Description,
			Sequence:   orgEntity.Sequence,
			Status:      string(orgEntity.Status),
			OrgCode:     orgEntity.OrgCode,
			OrgName:     orgEntity.OrgName,
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
		return nil, code.GetError(code.TenantGetPageListError)
	}
	list := make([]dtoorg.OrgPageListItem, 0, len(orgEntityList))
	for _, v := range orgEntityList {
		list = append(list, dtoorg.OrgPageListItem{
			OrgID: v.ID,
			OrgBaseInfo: objorg.OrgBaseInfo{
				Domain:      v.Domain,
				Logo:        v.Logo,
				Description: v.Description,
				Sequence:   v.Sequence,
				Status:      string(v.Status),
				OrgCode:     v.OrgCode,
				OrgName:     v.OrgName,
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

func resolveDomain(ctx *gin.Context, reqDomain string) string {
	host := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Host"))
	if host == "" && ctx.Request != nil {
		host = strings.TrimSpace(ctx.Request.Host)
	}
	if host == "" {
		host = strings.TrimSpace(reqDomain)
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
