package svcpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/apps/iam/object/objpermission"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type MenuSvc interface {
	Create(ctx *gin.Context, req *dtopermission.MenuCreateReq) (*dtopermission.MenuCreateResp, error)
	Delete(ctx *gin.Context, req *dtopermission.MenuDeleteReq) error
	Update(ctx *gin.Context, req *dtopermission.MenuUpdateReq) error
	Detail(ctx *gin.Context, req *dtopermission.MenuDetailReq) (*dtopermission.MenuDetailResp, error)
	PageList(ctx *gin.Context, req *dtopermission.MenuPageListReq) (*dtopermission.MenuPageListResp, error)
	Tree(ctx *gin.Context, req *dtopermission.MenuTreeReq) (*dtopermission.MenuTreeResp, error)
}

type menuSvc struct {
}

var _ MenuSvc = (*menuSvc)(nil)

func NewMenuSvc() MenuSvc {
	return &menuSvc{}
}

// Create 创建菜单管理
func (svc *menuSvc) Create(ctx *gin.Context, req *dtopermission.MenuCreateReq) (*dtopermission.MenuCreateResp, error) {
	insertEntity := &iammodel.MenuEntity{
		CacheType:     req.CacheType,
		ComponentPath: req.ComponentPath,
		Icon:          req.Icon,
		LinkType:      req.LinkType,
		MenuCode:      req.MenuCode,
		MenuName:      req.MenuName,
		MenuType:      req.MenuType,
		ParentID:      req.ParentID,
		Permission:    req.Permission,
		RoutePath:     req.RoutePath,
		SortOrder:     req.SortOrder,
		Status:        req.Status,
		Visibility:    req.Visibility,
	}

	if err := iamdao.NewMenuDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuCreate] daoMenu Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuCreateError)
	}
	return &dtopermission.MenuCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除菜单管理
func (svc *menuSvc) Delete(ctx *gin.Context, req *dtopermission.MenuDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	menuEntity, err := iamdao.NewMenuDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoMenu GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuDeleteError)
	}
	if menuEntity == nil || menuEntity.ID == 0 {
		return code.GetError(code.MenuNotExistError)
	}

	if err = iamdao.NewMenuDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoMenu Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuDeleteError)
	}
	return nil
}

// Update 更新菜单管理
func (svc *menuSvc) Update(ctx *gin.Context, req *dtopermission.MenuUpdateReq) error {
	menuEntity, err := iamdao.NewMenuDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuUpdate] daoMenu GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuUpdateError)
	}
	if menuEntity == nil || menuEntity.ID == 0 {
		return code.GetError(code.MenuNotExistError)
	}
	updateMap := map[string]any{
		"cache_type":     req.CacheType,
		"component_path": req.ComponentPath,
		"icon":           req.Icon,
		"link_type":      req.LinkType,
		"menu_code":      req.MenuCode,
		"menu_name":      req.MenuName,
		"menu_type":      req.MenuType,
		"parent_id":      req.ParentID,
		"permission":     req.Permission,
		"route_path":     req.RoutePath,
		"sort_order":     req.SortOrder,
		"status":         req.Status,
		"visibility":     req.Visibility,
	}
	if err = iamdao.NewMenuDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuUpdate] daoMenu UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuUpdateError)
	}
	return nil
}

// Detail 根据id获取菜单管理
func (svc *menuSvc) Detail(ctx *gin.Context, req *dtopermission.MenuDetailReq) (*dtopermission.MenuDetailResp, error) {
	menuEntity, err := iamdao.NewMenuDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuDetail] daoMenu GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuGetDetailError)
	}
	// 判断是否存在
	if menuEntity == nil || menuEntity.ID == 0 {
		return nil, code.GetError(code.MenuNotExistError)
	}
	resp := &dtopermission.MenuDetailResp{
		ID: menuEntity.ID,
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
			SortOrder:     menuEntity.SortOrder,
			Status:        menuEntity.Status,
			Visibility:    menuEntity.Visibility,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: menuEntity.CreatedAt.Unix(),
			UpdatedAt: menuEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取菜单管理列表
func (svc *menuSvc) PageList(ctx *gin.Context, req *dtopermission.MenuPageListReq) (*dtopermission.MenuPageListResp, error) {
	cond := &iamdao.MenuCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	menuEntityList, total, err := iamdao.NewMenuDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuPageList] daoMenu GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuGetPageListError)
	}
	list := make([]dtopermission.MenuPageListItem, 0, len(menuEntityList))
	for _, v := range menuEntityList {
		list = append(list, dtopermission.MenuPageListItem{
			ID: v.ID,
			MenuBaseInfo: objpermission.MenuBaseInfo{
				CacheType:     v.CacheType,
				TenantID:      v.TenantID,
				ComponentPath: v.ComponentPath,
				Icon:          v.Icon,
				LinkType:      v.LinkType,
				MenuCode:      v.MenuCode,
				MenuName:      v.MenuName,
				MenuType:      v.MenuType,
				ParentID:      v.ParentID,
				Permission:    v.Permission,
				RoutePath:     v.RoutePath,
				SortOrder:     v.SortOrder,
				Status:        v.Status,
				Visibility:    v.Visibility,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtopermission.MenuPageListResp{
		List:  list,
		Total: total,
	}, nil
}

// Tree 获取菜单树
func (svc *menuSvc) Tree(ctx *gin.Context, req *dtopermission.MenuTreeReq) (*dtopermission.MenuTreeResp, error) {
	cond := &iamdao.MenuCond{}
	allMenus, err := iamdao.NewMenuDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuTree] daoMenu GetListByCond fail, err:%v", err)
		return nil, code.GetError(code.MenuGetPageListError)
	}

	var parentID uint
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	menuMap := make(map[uint][]iammodel.MenuEntity)
	for _, menu := range allMenus {
		menuMap[menu.ParentID] = append(menuMap[menu.ParentID], menu)
	}

	var buildTree func(parentID uint) []dtopermission.MenuTreeNode
	buildTree = func(parentID uint) []dtopermission.MenuTreeNode {
		var nodes []dtopermission.MenuTreeNode
		children, ok := menuMap[parentID]
		if !ok {
			return nodes
		}
		for _, menu := range children {
			node := dtopermission.MenuTreeNode{
				ID: menu.ID,
				MenuBaseInfo: objpermission.MenuBaseInfo{
					CacheType:     menu.CacheType,
					TenantID:      menu.TenantID,
					ComponentPath: menu.ComponentPath,
					Icon:          menu.Icon,
					LinkType:      menu.LinkType,
					MenuCode:      menu.MenuCode,
					MenuName:      menu.MenuName,
					MenuType:      menu.MenuType,
					ParentID:      menu.ParentID,
					Permission:    menu.Permission,
					RoutePath:     menu.RoutePath,
					SortOrder:     menu.SortOrder,
					Status:        menu.Status,
					Visibility:    menu.Visibility,
				},
				OperatorBaseInfo: gobject.OperatorBaseInfo{
					UpdatedAt: menu.UpdatedAt.Unix(),
				},
				Children: buildTree(menu.ID),
			}
			nodes = append(nodes, node)
		}
		return nodes
	}

	tree := buildTree(parentID)
	return &dtopermission.MenuTreeResp{
		List: tree,
	}, nil
}
