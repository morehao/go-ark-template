package svcpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/organization"
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

type RoleSvc interface {
	Create(ctx *gin.Context, req *dtopermission.RoleCreateReq) (*dtopermission.RoleCreateResp, error)
	Delete(ctx *gin.Context, req *dtopermission.RoleDeleteReq) error
	Update(ctx *gin.Context, req *dtopermission.RoleUpdateReq) error
	Detail(ctx *gin.Context, req *dtopermission.RoleDetailReq) (*dtopermission.RoleDetailResp, error)
	PageList(ctx *gin.Context, req *dtopermission.RolePageListReq) (*dtopermission.RolePageListResp, error)
}

type roleSvc struct {
}

var _ RoleSvc = (*roleSvc)(nil)

func NewRoleSvc() RoleSvc {
	return &roleSvc{}
}

// Create 创建角色管理
func (svc *roleSvc) Create(ctx *gin.Context, req *dtopermission.RoleCreateReq) (*dtopermission.RoleCreateResp, error) {
	insertEntity := &iammodel.RoleEntity{
		DataScope:   req.DataScope,
		Description: req.Description,
		RoleCode:    req.RoleCode,
		RoleName:    req.RoleName,
		RoleType:    req.RoleType,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}

	if err := iamdao.NewRoleDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleCreate] daoRole Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleCreateError)
	}
	return &dtopermission.RoleCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除角色管理
func (svc *roleSvc) Delete(ctx *gin.Context, req *dtopermission.RoleDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	roleEntity, err := iamdao.NewRoleDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleDeleteError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return code.GetError(code.RoleNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, roleEntity.TenantID); err != nil {
		return err
	}

	if err = iamdao.NewRoleDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcpermission.Delete] daoRole Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleDeleteError)
	}
	return nil
}

// Update 更新角色管理
func (svc *roleSvc) Update(ctx *gin.Context, req *dtopermission.RoleUpdateReq) error {
	roleEntity, err := iamdao.NewRoleDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleUpdate] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleUpdateError)
	}
	if roleEntity == nil || roleEntity.ID == 0 {
		return code.GetError(code.RoleNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, roleEntity.TenantID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"data_scope":  req.DataScope,
		"description": req.Description,
		"role_code":   req.RoleCode,
		"role_name":   req.RoleName,
		"role_type":   req.RoleType,
		"sort_order":  req.SortOrder,
		"status":      req.Status,
	}
	if err = iamdao.NewRoleDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleUpdate] daoRole UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.RoleUpdateError)
	}
	return nil
}

// Detail 根据id获取角色管理
func (svc *roleSvc) Detail(ctx *gin.Context, req *dtopermission.RoleDetailReq) (*dtopermission.RoleDetailResp, error) {
	roleEntity, err := iamdao.NewRoleDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RoleDetail] daoRole GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleGetDetailError)
	}
	// 判断是否存在
	if roleEntity == nil || roleEntity.ID == 0 {
		return nil, code.GetError(code.RoleNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, roleEntity.TenantID); err != nil {
		return nil, err
	}
	resp := &dtopermission.RoleDetailResp{
		ID: roleEntity.ID,
		RoleBaseInfo: objpermission.RoleBaseInfo{
			TenantID:    roleEntity.TenantID,
			DataScope:   roleEntity.DataScope,
			Description: roleEntity.Description,
			RoleCode:    roleEntity.RoleCode,
			RoleName:    roleEntity.RoleName,
			RoleType:    roleEntity.RoleType,
			SortOrder:   roleEntity.SortOrder,
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
	cond := &iamdao.RoleCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	roleEntityList, total, err := iamdao.NewRoleDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcpermission.RolePageList] daoRole GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.RoleGetPageListError)
	}
	list := make([]dtopermission.RolePageListItem, 0, len(roleEntityList))
	for _, v := range roleEntityList {
		list = append(list, dtopermission.RolePageListItem{
			ID: v.ID,
			RoleBaseInfo: objpermission.RoleBaseInfo{
				TenantID:    v.TenantID,
				DataScope:   v.DataScope,
				Description: v.Description,
				RoleCode:    v.RoleCode,
				RoleName:    v.RoleName,
				RoleType:    v.RoleType,
				SortOrder:   v.SortOrder,
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
