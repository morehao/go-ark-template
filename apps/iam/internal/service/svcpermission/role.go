package svcpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/apps/iam/object/objpermission"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
	"gorm.io/gorm"
)

type RoleSvc interface {
	Create(ctx *gin.Context, req *dtopermission.RoleCreateReq) (*dtopermission.RoleCreateResp, error)
	Delete(ctx *gin.Context, req *dtopermission.RoleDeleteReq) error
	Update(ctx *gin.Context, req *dtopermission.RoleUpdateReq) error
	Detail(ctx *gin.Context, req *dtopermission.RoleDetailReq) (*dtopermission.RoleDetailResp, error)
	PageList(ctx *gin.Context, req *dtopermission.RolePageListReq) (*dtopermission.RolePageListResp, error)
	AssignMenus(ctx *gin.Context, req *dtopermission.RoleAssignMenusReq) error
	ListMenus(ctx *gin.Context, req *dtopermission.RoleListMenusReq) (*dtopermission.RoleMenuListResp, error)
}

type roleSvc struct {
}

var _ RoleSvc = (*roleSvc)(nil)

func NewRoleSvc() RoleSvc {
	return &roleSvc{}
}

// Create 创建角色管理
func (svc *roleSvc) Create(ctx *gin.Context, req *dtopermission.RoleCreateReq) (*dtopermission.RoleCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.RoleEntity{
		TenantID:   tenantID,
		DataScope:  model.RoleDataScope(req.DataScope),
		Description: req.Description,
		RoleCode:    req.RoleCode,
		RoleName:    req.RoleName,
		RoleType:    model.RoleType(req.RoleType),
		Sequence:   req.Sequence,
		Status:      model.RoleStatus(req.Status),
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	if err := dao.NewRoleDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleCreate] daoRole Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleCreateError)
	}
	return &dtopermission.RoleCreateResp{
		RoleID: insertEntity.ID,
	}, nil
}

// Delete 删除角色管理
func (svc *roleSvc) Delete(ctx *gin.Context, req *dtopermission.RoleDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	roleEntity, err := dao.NewRoleDao().GetByID(ctx, req.RoleID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleDeleteError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return code.GetError(code.RoleNotExistError)
	}

	if err = dao.NewRoleDao().Delete(ctx, req.RoleID, userID); err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoRole Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleDeleteError)
	}
	return nil
}

// Update 更新角色管理
func (svc *roleSvc) Update(ctx *gin.Context, req *dtopermission.RoleUpdateReq) error {
	roleEntity, err := dao.NewRoleDao().GetByID(ctx, req.RoleID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleUpdate] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleUpdateError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return code.GetError(code.RoleNotExistError)
	}
	updateMap := map[string]any{
		"data_scope":  req.DataScope,
		"description": req.Description,
		"role_code":   req.RoleCode,
		"role_name":   req.RoleName,
		"role_type":   req.RoleType,
		"sequence":  req.Sequence,
		"status":      req.Status,
	}
	if err = dao.NewRoleDao().UpdateMap(ctx, req.RoleID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleUpdate] daoRole UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleUpdateError)
	}
	return nil
}

// Detail 根据id获取角色管理
func (svc *roleSvc) Detail(ctx *gin.Context, req *dtopermission.RoleDetailReq) (*dtopermission.RoleDetailResp, error) {
	roleEntity, err := dao.NewRoleDao().GetByID(ctx, req.RoleID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleDetail] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleGetDetailError)
	}
	// 判断是否存在
	if roleEntity == nil || roleEntity.ID == 0 {
		return nil, code.GetError(code.RoleNotExistError)
	}
	resp := &dtopermission.RoleDetailResp{
		RoleID: roleEntity.ID,
		RoleBaseInfo: objpermission.RoleBaseInfo{
			TenantID:    roleEntity.TenantID,
			DataScope:   roleEntity.DataScope,
			Description: roleEntity.Description,
			RoleCode:    roleEntity.RoleCode,
			RoleName:    roleEntity.RoleName,
			RoleType:    roleEntity.RoleType,
			Sequence:   roleEntity.Sequence,
			Status:      roleEntity.Status,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: roleEntity.CreatedAt.Unix(),
			UpdatedAt: roleEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取角色管理列表
func (svc *roleSvc) PageList(ctx *gin.Context, req *dtopermission.RolePageListReq) (*dtopermission.RolePageListResp, error) {
	cond := &dao.RoleCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	roleEntityList, total, err := dao.NewRoleDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RolePageList] daoRole GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleGetPageListError)
	}
	list := make([]dtopermission.RolePageListItem, 0, len(roleEntityList))
	for _, v := range roleEntityList {
		list = append(list, dtopermission.RolePageListItem{
			RoleID: v.ID,
			RoleBaseInfo: objpermission.RoleBaseInfo{
				TenantID:    v.TenantID,
				DataScope:   v.DataScope,
				Description: v.Description,
				RoleCode:    v.RoleCode,
				RoleName:    v.RoleName,
				RoleType:    v.RoleType,
				Sequence:   v.Sequence,
				Status:      v.Status,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtopermission.RolePageListResp{
		List:  list,
		Total: total,
	}, nil
}

// AssignMenus 角色分配菜单(全量替换)
func (svc *roleSvc) AssignMenus(ctx *gin.Context, req *dtopermission.RoleAssignMenusReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	// 检查角色是否存在
	roleEntity, err := dao.NewRoleDao().GetByID(ctx, req.RoleID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.AssignMenus] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleUpdateError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return code.GetError(code.RoleNotExistError)
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		roleMenuDao := dao.NewRoleMenuDao().WithTx(tx)

		// 删除该角色的所有已有菜单关联
		existingList, err := roleMenuDao.GetListByCond(ctx, &dao.RoleMenuCond{
			RoleID:   req.RoleID,
			TenantID: tenantID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcpermission.AssignMenus] GetListByCond fail, err:%v, roleID:%d", err, req.RoleID)
			return code.GetError(code.RoleUpdateError)
		}
		for _, existing := range existingList {
			if err := roleMenuDao.Delete(ctx, existing.ID, operatorID); err != nil {
				glog.Errorf(ctx, "[svcpermission.AssignMenus] Delete fail, err:%v, id:%d", err, existing.ID)
				return code.GetError(code.RoleUpdateError)
			}
		}

		// 批量插入新的菜单关联
		for _, menuID := range req.MenuIDs {
			entity := &model.RoleMenuEntity{
				TenantID:  tenantID,
				RoleID:    req.RoleID,
				MenuID:    menuID,
				CreatedBy: operatorID,
				UpdatedBy: operatorID,
			}
			if err := roleMenuDao.Insert(ctx, entity); err != nil {
				glog.Errorf(ctx, "[svcpermission.AssignMenus] Insert fail, err:%v, roleID:%d, menuID:%d", err, req.RoleID, menuID)
				return code.GetError(code.RoleUpdateError)
			}
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcpermission.AssignMenus] Transaction fail, err:%v", txErr)
		return code.GetError(code.RoleUpdateError)
	}
	return nil
}

// ListMenus 查询角色已分配的菜单列表
func (svc *roleSvc) ListMenus(ctx *gin.Context, req *dtopermission.RoleListMenusReq) (*dtopermission.RoleMenuListResp, error) {
	// 查询角色是否存在
	roleEntity, err := dao.NewRoleDao().GetByID(ctx, req.RoleID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.ListMenus] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleGetDetailError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return nil, code.GetError(code.RoleNotExistError)
	}

	// 查询角色菜单关联
	roleMenuList, err := dao.NewRoleMenuDao().GetListByCond(ctx, &dao.RoleMenuCond{
		RoleID: req.RoleID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.ListMenus] daoRoleMenu GetListByCond fail, err:%v, roleID:%d", err, req.RoleID)
		return nil, code.GetError(code.RoleGetDetailError)
	}

	list := make([]dtopermission.RoleMenuListItem, 0, len(roleMenuList))
	for _, rm := range roleMenuList {
		menuEntity, err := dao.NewMenuDao().GetByID(ctx, rm.MenuID)
		if err != nil || menuEntity == nil || menuEntity.ID == 0 {
			continue
		}
		list = append(list, dtopermission.RoleMenuListItem{
			MenuID: menuEntity.ID,
			MenuBaseInfo: objpermission.MenuBaseInfo{
				CacheType:     menuEntity.CacheType,
				TenantID:      menuEntity.TenantID,
				ComponentPath: menuEntity.ComponentPath,
				Icon:          menuEntity.Icon,
				LinkType:      menuEntity.LinkType,
				MenuCode:      menuEntity.MenuCode,
				MenuName:      menuEntity.MenuName,
				MenuType:      menuEntity.MenuType,
				ParentID:      menuEntity.ParentID,
				Permission:    menuEntity.Permission,
				RoutePath:     menuEntity.RoutePath,
				Sequence:     menuEntity.Sequence,
				Status:        menuEntity.Status,
				Visibility:    menuEntity.Visibility,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: menuEntity.UpdatedAt.Unix(),
			},
		})
	}

	return &dtopermission.RoleMenuListResp{
		List: list,
	}, nil
}
