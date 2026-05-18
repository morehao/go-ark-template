package svcpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/iam/object/objpermission"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gtree"
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

type menuSequenceComparator struct{}

func (c menuSequenceComparator) Compare(a, b *dtopermission.MenuTreeNode) int {
	if a.Sequence < b.Sequence {
		return -1
	} else if a.Sequence > b.Sequence {
		return 1
	}
	return 0
}

// Create 创建菜单管理
func (svc *menuSvc) Create(ctx *gin.Context, req *dtopermission.MenuCreateReq) (*dtopermission.MenuCreateResp, error) {
	insertEntity := &model.MenuEntity{
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
		Sequence:     req.Sequence,
		Status:        req.Status,
		Visibility:    req.Visibility,
		AccessPolicy:  model.AccessPoliciesToMask(req.AccessPolicies),
	}

	if err := dao.NewMenuDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuCreate] daoMenu Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuCreateError)
	}
	return &dtopermission.MenuCreateResp{
		MenuID: insertEntity.ID,
	}, nil
}

// Delete 删除菜单管理
func (svc *menuSvc) Delete(ctx *gin.Context, req *dtopermission.MenuDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	menuEntity, err := dao.NewMenuDao().GetByID(ctx, req.MenuID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoMenu GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuDeleteError)
	}
	if menuEntity == nil || menuEntity.ID == 0 {
		return code.GetError(code.MenuNotExistError)
	}

	if err = dao.NewMenuDao().Delete(ctx, req.MenuID, userID); err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoMenu Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuDeleteError)
	}
	return nil
}

// Update 更新菜单管理
func (svc *menuSvc) Update(ctx *gin.Context, req *dtopermission.MenuUpdateReq) error {
	menuEntity, err := dao.NewMenuDao().GetByID(ctx, req.MenuID)
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
		"sequence":     req.Sequence,
		"status":         req.Status,
		"visibility":     req.Visibility,
		"access_policy":  model.AccessPoliciesToMask(req.AccessPolicies),
	}
	if err = dao.NewMenuDao().UpdateMap(ctx, req.MenuID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuUpdate] daoMenu UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MenuUpdateError)
	}
	return nil
}

// Detail 根据id获取菜单管理
func (svc *menuSvc) Detail(ctx *gin.Context, req *dtopermission.MenuDetailReq) (*dtopermission.MenuDetailResp, error) {
	menuEntity, err := dao.NewMenuDao().GetByID(ctx, req.MenuID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuDetail] daoMenu GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuGetDetailError)
	}
	// 判断是否存在
	if menuEntity == nil || menuEntity.ID == 0 {
		return nil, code.GetError(code.MenuNotExistError)
	}
	resp := &dtopermission.MenuDetailResp{
		MenuID: menuEntity.ID,
		MenuBaseInfo: objpermission.MenuBaseInfo{
			CacheType:      menuEntity.CacheType,
			TenantID:       menuEntity.TenantID,
			ComponentPath:  menuEntity.ComponentPath,
			Icon:           menuEntity.Icon,
			LinkType:       menuEntity.LinkType,
			MenuCode:       menuEntity.MenuCode,
			MenuName:       menuEntity.MenuName,
			MenuType:       menuEntity.MenuType,
			ParentID:       menuEntity.ParentID,
			Permission:     menuEntity.Permission,
			RoutePath:      menuEntity.RoutePath,
			Sequence:      menuEntity.Sequence,
			Status:         menuEntity.Status,
			Visibility:     menuEntity.Visibility,
			AccessPolicies: menuEntity.AccessPolicy.ToStrings(),
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
	cond := &dao.MenuCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	menuEntityList, total, err := dao.NewMenuDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuPageList] daoMenu GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MenuGetPageListError)
	}
	list := make([]dtopermission.MenuPageListItem, 0, len(menuEntityList))
	for _, v := range menuEntityList {
		list = append(list, dtopermission.MenuPageListItem{
			MenuID: v.ID,
			MenuBaseInfo: objpermission.MenuBaseInfo{
				CacheType:      v.CacheType,
				TenantID:       v.TenantID,
				ComponentPath:  v.ComponentPath,
				Icon:           v.Icon,
				LinkType:       v.LinkType,
				MenuCode:       v.MenuCode,
				MenuName:       v.MenuName,
				MenuType:       v.MenuType,
				ParentID:       v.ParentID,
				Permission:     v.Permission,
				RoutePath:      v.RoutePath,
				Sequence:      v.Sequence,
				Status:         v.Status,
				Visibility:     v.Visibility,
				AccessPolicies: v.AccessPolicy.ToStrings(),
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
	cond := &dao.MenuCond{}
	allMenus, err := dao.NewMenuDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.MenuTree] daoMenu GetListByCond fail, err:%v", err)
		return nil, code.GetError(code.MenuGetPageListError)
	}

	nodes := make([]*dtopermission.MenuTreeNode, len(allMenus))
	for i, menu := range allMenus {
		nodes[i] = &dtopermission.MenuTreeNode{
			MenuID: menu.ID,
			MenuBaseInfo: objpermission.MenuBaseInfo{
				CacheType:      menu.CacheType,
				TenantID:       menu.TenantID,
				ComponentPath:  menu.ComponentPath,
				Icon:           menu.Icon,
				LinkType:       menu.LinkType,
				MenuCode:       menu.MenuCode,
				MenuName:       menu.MenuName,
				MenuType:       menu.MenuType,
				ParentID:       menu.ParentID,
				Permission:     menu.Permission,
				RoutePath:      menu.RoutePath,
				Sequence:      menu.Sequence,
				Status:         menu.Status,
				Visibility:     menu.Visibility,
				AccessPolicies: menu.AccessPolicy.ToStrings(),
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: menu.UpdatedAt.Unix(),
			},
		}
	}

	builder := gtree.NewTreeBuilder[uint, *dtopermission.MenuTreeNode](
		gtree.WithComparator(menuSequenceComparator{}),
	)
	tree := builder.Build(nodes)

	roots := tree.Roots
	if req.ParentID != nil && *req.ParentID != 0 {
		if subtree, ok := tree.NodeMap[*req.ParentID]; ok {
			roots = []*dtopermission.MenuTreeNode{subtree}
		}
	}

	result := make([]dtopermission.MenuTreeNode, len(roots))
	for i, root := range roots {
		result[i] = *root
		result[i].Children = svc.buildJSONChildren(tree, root.MenuID)
	}

	return &dtopermission.MenuTreeResp{List: result}, nil
}

func (svc *menuSvc) buildJSONChildren(tree *gtree.Tree[uint, *dtopermission.MenuTreeNode], parentID uint) []dtopermission.MenuTreeNode {
	children, ok := tree.Children(parentID)
	if !ok || len(children) == 0 {
		return nil
	}
	result := make([]dtopermission.MenuTreeNode, len(children))
	for i, child := range children {
		result[i] = *child
		result[i].Children = svc.buildJSONChildren(tree, child.MenuID)
	}
	return result
}
