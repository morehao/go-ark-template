package svcpermission

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/ginext"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type AuthorizationSvc interface {
	AssignRolesToUser(ctx *gin.Context, req *dtopermission.UserAssignRolesReq) error
	RemoveRolesFromUser(ctx *gin.Context, req *dtopermission.UserRemoveRolesReq) error
	ListUserRoles(ctx *gin.Context, req *dtopermission.UserRoleListReq) (*dtopermission.UserRoleListResp, error)
	GetUserPermissions(ctx *gin.Context) (*dtopermission.UserPermissionsResp, error)
}

type authorizationSvc struct{}

var _ AuthorizationSvc = (*authorizationSvc)(nil)

func NewAuthorizationSvc() AuthorizationSvc {
	return &authorizationSvc{}
}

// AssignRolesToUser 为用户分配角色
func (svc *authorizationSvc) AssignRolesToUser(ctx *gin.Context, req *dtopermission.UserAssignRolesReq) error {
	tenantID := ginext.GetTenantID(ctx)
	userRoleDao := iamdao.NewUserRoleDao()

	existingList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
		BaseCond: &genericdao.BaseCond{},
		UserID:   req.UserID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.AssignRolesToUser] userRoleDao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserAssignRolesError)
	}

	existingMap := make(map[uint]bool)
	for _, e := range existingList {
		existingMap[e.RoleID] = true
	}

	newEntities := make(iammodel.UserRoleEntityList, 0, len(req.RoleIDs))
	for _, roleID := range req.RoleIDs {
		if existingMap[roleID] {
			continue
		}
		newEntities = append(newEntities, iammodel.UserRoleEntity{
			UserID:   req.UserID,
			RoleID:   roleID,
			TenantID: tenantID,
		})
	}

	if len(newEntities) > 0 {
		if err := userRoleDao.BatchInsert(ctx, newEntities); err != nil {
			glog.Errorf(ctx, "[svcpermission.AssignRolesToUser] userRoleDao BatchInsert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return code.GetError(code.UserAssignRolesError)
		}
	}
	return nil
}

// RemoveRolesFromUser 移除用户角色
func (svc *authorizationSvc) RemoveRolesFromUser(ctx *gin.Context, req *dtopermission.UserRemoveRolesReq) error {
	userID := gincontext.GetUserID(ctx)
	tenantID := ginext.GetTenantID(ctx)
	userRoleDao := iamdao.NewUserRoleDao()

	existingList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
		BaseCond: &genericdao.BaseCond{},
		UserID:   req.UserID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RemoveRolesFromUser] userRoleDao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserRemoveRolesError)
	}

	roleIDSet := make(map[uint]bool)
	for _, id := range req.RoleIDs {
		roleIDSet[id] = true
	}

	for _, e := range existingList {
		if roleIDSet[e.RoleID] {
			if err := userRoleDao.Delete(ctx, e.ID, userID); err != nil {
				glog.Errorf(ctx, "[svcpermission.RemoveRolesFromUser] userRoleDao Delete fail, err:%v, id:%d", err, e.ID)
				return code.GetError(code.UserRemoveRolesError)
			}
		}
	}
	return nil
}

// ListUserRoles 获取用户角色列表
func (svc *authorizationSvc) ListUserRoles(ctx *gin.Context, req *dtopermission.UserRoleListReq) (*dtopermission.UserRoleListResp, error) {
	tenantID := ginext.GetTenantID(ctx)
	userRoleDao := iamdao.NewUserRoleDao()

	userRoleList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
		BaseCond: &genericdao.BaseCond{},
		UserID:   req.UserID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.ListUserRoles] userRoleDao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserListRolesError)
	}

	if len(userRoleList) == 0 {
		return &dtopermission.UserRoleListResp{List: []dtopermission.UserRoleItem{}}, nil
	}

	roleIDs := make([]uint, 0, len(userRoleList))
	for _, ur := range userRoleList {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	roleDao := iamdao.NewRoleDao()
	roles, err := roleDao.GetListByCond(ctx, &iamdao.RoleCond{
		BaseCond: &genericdao.BaseCond{
			IDs: roleIDs,
		},
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.ListUserRoles] roleDao GetListByCond fail, err:%v, roleIDs:%v", err, roleIDs)
		return nil, code.GetError(code.UserListRolesError)
	}

	list := make([]dtopermission.UserRoleItem, 0, len(roles))
	for _, r := range roles {
		list = append(list, dtopermission.UserRoleItem{
			RoleID:   r.ID,
			RoleName: r.RoleName,
			RoleCode: r.RoleCode,
			RoleType: r.RoleType,
		})
	}

	return &dtopermission.UserRoleListResp{List: list}, nil
}

// GetUserPermissions 获取当前用户权限
func (svc *authorizationSvc) GetUserPermissions(ctx *gin.Context) (*dtopermission.UserPermissionsResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := ginext.GetTenantID(ctx)

	// Get user's roles
	userRoleDao := iamdao.NewUserRoleDao()
	userRoleList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
		BaseCond: &genericdao.BaseCond{},
		UserID:   userID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.GetUserPermissions] userRoleDao GetListByCond fail, err:%v, userID:%d", err, userID)
		return nil, code.GetError(code.UserGetPermissionsError)
	}

	if len(userRoleList) == 0 {
		return &dtopermission.UserPermissionsResp{
			Menus:       []dtopermission.MenuTreeNode{},
			Permissions: []string{},
		}, nil
	}

	roleIDs := make([]uint, 0, len(userRoleList))
	for _, ur := range userRoleList {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// Get menus for all roles
	roleMenuDao := iamdao.NewRoleMenuDao()
	roleMenuList, err := roleMenuDao.GetListByCond(ctx, &iamdao.RoleMenuCond{
		BaseCond: &genericdao.BaseCond{},
		RoleIDs:  roleIDs,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.GetUserPermissions] roleMenuDao GetListByCond fail, err:%v, roleIDs:%v", err, roleIDs)
		return nil, code.GetError(code.UserGetPermissionsError)
	}

	if len(roleMenuList) == 0 {
		return &dtopermission.UserPermissionsResp{
			Menus:       []dtopermission.MenuTreeNode{},
			Permissions: []string{},
		}, nil
	}

	// Deduplicate menu IDs
	menuIDMap := make(map[uint]bool)
	for _, rm := range roleMenuList {
		menuIDMap[rm.MenuID] = true
	}
	menuIDs := make([]uint, 0, len(menuIDMap))
	for id := range menuIDMap {
		menuIDs = append(menuIDs, id)
	}

	menuDao := iamdao.NewMenuDao()
	menus, err := menuDao.GetListByCond(ctx, &iamdao.MenuCond{
		BaseCond: &genericdao.BaseCond{
			IDs: menuIDs,
		},
	})
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.GetUserPermissions] menuDao GetListByCond fail, err:%v, menuIDs:%v", err, menuIDs)
		return nil, code.GetError(code.UserGetPermissionsError)
	}

	// Build permissions list and menu tree
	permissions := make([]string, 0)
	for _, m := range menus {
		if m.Permission != "" {
			permissions = append(permissions, m.Permission)
		}
	}

	menuTree := buildMenuTree(menus, 0)

	return &dtopermission.UserPermissionsResp{
		Menus:       menuTree,
		Permissions: permissions,
	}, nil
}

func buildMenuTree(menus iammodel.MenuEntityList, parentID uint) []dtopermission.MenuTreeNode {
	nodes := make([]dtopermission.MenuTreeNode, 0)
	for _, m := range menus {
		if m.ParentID == parentID {
			node := dtopermission.MenuTreeNode{
				ID:            m.ID,
				MenuName:      m.MenuName,
				MenuCode:      m.MenuCode,
				MenuType:      m.MenuType,
				ParentID:      m.ParentID,
				RoutePath:     m.RoutePath,
				ComponentPath: m.ComponentPath,
				Permission:    m.Permission,
				Icon:          m.Icon,
				SortOrder:     m.SortOrder,
				Visibility:    m.Visibility,
				Children:      buildMenuTree(menus, m.ID),
			}
			nodes = append(nodes, node)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].SortOrder < nodes[j].SortOrder
	})
	return nodes
}
