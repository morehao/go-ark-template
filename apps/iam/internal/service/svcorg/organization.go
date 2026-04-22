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
	GetConfigs(ctx *gin.Context, req *dtoorg.GetOrganizationConfigsReq) (*dtoorg.GetOrganizationConfigsResp, error)
	Detail(ctx *gin.Context, req *dtoorg.OrgDetailReq) (*dtoorg.OrgDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.OrgPageListReq) (*dtoorg.OrgPageListResp, error)
	ListConfig(ctx *gin.Context) (*dtoorg.OrgConfigListResp, error)
}

type organizationSvc struct {
}

var _ OrganizationSvc = (*organizationSvc)(nil)

func NewOrganizationSvc() OrganizationSvc {
	return &organizationSvc{}
}

func (svc *organizationSvc) Create(ctx *gin.Context, req *dtoorg.OrgCreateReq) (*dtoorg.OrgCreateResp, error) {
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.OrganizationEntity{
		Domain:      req.Domain,
		Logo:        req.Logo,
		Description: req.Description,
		SortOrder:   req.SortOrder,
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
					ConfigType:  meta.Type,
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
	updateMap := map[string]any{
		"domain":      req.Domain,
		"logo":        req.Logo,
		"description": req.Description,
		"sort_order":  req.SortOrder,
		"status":      req.Status,
		"org_code":    req.OrgCode,
		"org_name":    req.OrgName,
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := dao.NewOrganizationDao().WithTx(tx).UpdateMap(ctx, req.OrgID, updateMap); err != nil {
			glog.Errorf(ctx, "[svcorg.OrgUpdate] daoOrg UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.TenantUpdateError)
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
					ConfigType:  meta.Type,
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

func (svc *organizationSvc) GetConfigs(ctx *gin.Context, req *dtoorg.GetOrganizationConfigsReq) (*dtoorg.GetOrganizationConfigsResp, error) {
	domain := resolveDomain(ctx, req.Domain)
	if domain == "" {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
		Domain: domain,
		Status: model.OrgStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetConfigsByDomain] daoOrg GetByCond fail, err:%v, domain:%s", err, domain)
		return nil, code.GetError(code.AuthLoginError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	configEntityList, err := dao.NewOrganizationConfigDao().GetListByCond(ctx, &dao.OrganizationConfigCond{
		OrgID: orgEntity.ID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcorg.GetConfigsByDomain] daoOrgConfig GetListByCond fail, err:%v, orgID:%d", err, orgEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
	}

	configs := make(map[string]map[string]string)
	for _, v := range configEntityList {
		if _, ok := configs[v.ConfigGroup]; !ok {
			configs[v.ConfigGroup] = make(map[string]string)
		}
		configs[v.ConfigGroup][v.ConfigKey] = v.ConfigValue
	}

	return &dtoorg.GetOrganizationConfigsResp{
		OrgID:   orgEntity.ID,
		OrgName: orgEntity.OrgName,
		Domain:  orgEntity.Domain,
		Logo:    orgEntity.Logo,
		Status:  orgEntity.Status,
		Configs: configs,
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

	resp := &dtoorg.OrgDetailResp{
		OrgID:   orgEntity.ID,
		Configs: configs,
		OrgBaseInfo: objorg.OrgBaseInfo{
			Domain:      orgEntity.Domain,
			Logo:        orgEntity.Logo,
			Description: orgEntity.Description,
			SortOrder:   orgEntity.SortOrder,
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
				SortOrder:   v.SortOrder,
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

func (svc *organizationSvc) ListConfig(ctx *gin.Context) (*dtoorg.OrgConfigListResp, error) {
	configs := make([]dtoorg.OrgConfigMetaResp, 0, len(model.OrgConfigMetaList))
	for _, meta := range model.OrgConfigMetaList {
		options := make([]dtoorg.OrgConfigOptionResp, 0, len(meta.Options))
		for _, opt := range meta.Options {
			options = append(options, dtoorg.OrgConfigOptionResp{
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
	return &dtoorg.OrgConfigListResp{Configs: configs}, nil
}
