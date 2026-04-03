package svcorg

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type DepartmentSvc interface {
	Create(ctx *gin.Context, req *dtoorg.DepartmentCreateReq) (*dtoorg.DepartmentCreateResp, error)
	Delete(ctx *gin.Context, req *dtoorg.DepartmentDeleteReq) error
	Update(ctx *gin.Context, req *dtoorg.DepartmentUpdateReq) error
	Detail(ctx *gin.Context, req *dtoorg.DepartmentDetailReq) (*dtoorg.DepartmentDetailResp, error)
	PageList(ctx *gin.Context, req *dtoorg.DepartmentPageListReq) (*dtoorg.DepartmentPageListResp, error)
	Tree(ctx *gin.Context, req *dtoorg.DepartmentTreeReq) (*dtoorg.DepartmentTreeResp, error)
}

type departmentSvc struct {
}

var _ DepartmentSvc = (*departmentSvc)(nil)

func NewDepartmentSvc() DepartmentSvc {
	return &departmentSvc{}
}

// Create 创建部门管理
func (svc *departmentSvc) Create(ctx *gin.Context, req *dtoorg.DepartmentCreateReq) (*dtoorg.DepartmentCreateResp, error) {
	operatorID := gincontext.GetUserID(ctx)

	tenantID := gincontext.GetTenantID(ctx)

	var deptLevel int32 = 1
	var deptPath string = "/"

	if req.ParentID > 0 {
		parentDept, err := iamdao.NewDepartmentDao().GetByID(ctx, req.ParentID)
		if err != nil || parentDept == nil || parentDept.ID == 0 {
			glog.Errorf(ctx, "[svcorg.DepartmentCreate] parent department not found, parentID:%d", req.ParentID)
			return nil, code.GetError(code.DepartmentNotExistError)
		}
		if parentDept.TenantID != tenantID {
			return nil, code.GetError(code.TenantScopeForbiddenError)
		}
		deptLevel = parentDept.DeptLevel + 1
		deptPath = fmt.Sprintf("%s%d/", parentDept.DeptPath, parentDept.ID)
	} else {
		deptPath = fmt.Sprintf("/%d/", 0)
	}

	insertEntity := &iammodel.DepartmentEntity{
		TenantID:  tenantID,
		DeptCode:  req.DeptCode,
		DeptLevel: deptLevel,
		DeptName:  req.DeptName,
		DeptPath:  deptPath,
		LeaderID:  req.LeaderID,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := iamdao.NewDepartmentDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentCreate] daoDepartment Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentCreateError)
	}

	updatePathMap := map[string]any{
		"dept_path": fmt.Sprintf("/%d/", insertEntity.ID),
	}
	if err := iamdao.NewDepartmentDao().UpdateMap(ctx, insertEntity.ID, updatePathMap); err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentCreate] update dept_path fail, err:%v", err)
	}

	return &dtoorg.DepartmentCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除部门管理
func (svc *departmentSvc) Delete(ctx *gin.Context, req *dtoorg.DepartmentDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	departmentEntity, err := iamdao.NewDepartmentDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoDepartment GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentDeleteError)
	}
	if departmentEntity == nil || departmentEntity.ID == 0 {
		return code.GetError(code.DepartmentNotExistError)
	}

	if err = iamdao.NewDepartmentDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoDepartment Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentDeleteError)
	}
	return nil
}

// Update 更新部门管理
func (svc *departmentSvc) Update(ctx *gin.Context, req *dtoorg.DepartmentUpdateReq) error {
	departmentEntity, err := iamdao.NewDepartmentDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentUpdate] daoDepartment GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentUpdateError)
	}
	if departmentEntity == nil || departmentEntity.ID == 0 {
		return code.GetError(code.DepartmentUpdateError)
	}
	updateMap := map[string]any{
		"dept_code":  req.DeptCode,
		"dept_level": req.DeptLevel,
		"dept_name":  req.DeptName,
		"dept_path":  req.DeptPath,
		"leader_id":  req.LeaderID,
		"parent_id":  req.ParentID,
		"sort_order": req.SortOrder,
		"status":     req.Status,
	}
	if err = iamdao.NewDepartmentDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentUpdate] daoDepartment UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentUpdateError)
	}
	return nil
}

// Detail 根据id获取部门管理
func (svc *departmentSvc) Detail(ctx *gin.Context, req *dtoorg.DepartmentDetailReq) (*dtoorg.DepartmentDetailResp, error) {
	departmentEntity, err := iamdao.NewDepartmentDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentDetail] daoDepartment GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentGetDetailError)
	}
	// 判断是否存在
	if departmentEntity == nil || departmentEntity.ID == 0 {
		return nil, code.GetError(code.DepartmentNotExistError)
	}
	resp := &dtoorg.DepartmentDetailResp{
		ID: departmentEntity.ID,
		DepartmentBaseInfo: objorg.DepartmentBaseInfo{
			TenantID:  departmentEntity.TenantID,
			DeptCode:  departmentEntity.DeptCode,
			DeptLevel: departmentEntity.DeptLevel,
			DeptName:  departmentEntity.DeptName,
			DeptPath:  departmentEntity.DeptPath,
			LeaderID:  departmentEntity.LeaderID,
			ParentID:  departmentEntity.ParentID,
			SortOrder: departmentEntity.SortOrder,
			Status:    departmentEntity.Status,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: departmentEntity.CreatedAt.Unix(),
			UpdatedAt: departmentEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取部门管理列表
func (svc *departmentSvc) PageList(ctx *gin.Context, req *dtoorg.DepartmentPageListReq) (*dtoorg.DepartmentPageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &iamdao.DepartmentCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID: tenantID,
	}
	departmentEntityList, total, err := iamdao.NewDepartmentDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentPageList] daoDepartment GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentGetPageListError)
	}
	list := make([]dtoorg.DepartmentPageListItem, 0, len(departmentEntityList))
	for _, v := range departmentEntityList {
		list = append(list, dtoorg.DepartmentPageListItem{
			ID: v.ID,
			DepartmentBaseInfo: objorg.DepartmentBaseInfo{
				TenantID:  v.TenantID,
				DeptCode:  v.DeptCode,
				DeptLevel: v.DeptLevel,
				DeptName:  v.DeptName,
				DeptPath:  v.DeptPath,
				LeaderID:  v.LeaderID,
				ParentID:  v.ParentID,
				SortOrder: v.SortOrder,
				Status:    v.Status,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoorg.DepartmentPageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *departmentSvc) Tree(ctx *gin.Context, req *dtoorg.DepartmentTreeReq) (*dtoorg.DepartmentTreeResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &iamdao.DepartmentCond{
		TenantID: tenantID,
	}
	allDepts, _, err := iamdao.NewDepartmentDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentTree] daoDepartment GetPageListByCond fail, err:%v", err)
		return nil, code.GetError(code.DepartmentGetPageListError)
	}

	var parentID uint
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	deptMap := make(map[uint][]iammodel.DepartmentEntity)
	for _, dept := range allDepts {
		deptMap[dept.ParentID] = append(deptMap[dept.ParentID], dept)
	}

	var buildTree func(parentID uint) []dtoorg.DepartmentTreeNode
	buildTree = func(parentID uint) []dtoorg.DepartmentTreeNode {
		var nodes []dtoorg.DepartmentTreeNode
		children, ok := deptMap[parentID]
		if !ok {
			return nodes
		}
		for _, dept := range children {
			node := dtoorg.DepartmentTreeNode{
				ID: dept.ID,
				DepartmentBaseInfo: objorg.DepartmentBaseInfo{
					TenantID:  dept.TenantID,
					DeptCode:  dept.DeptCode,
					DeptLevel: dept.DeptLevel,
					DeptName:  dept.DeptName,
					DeptPath:  dept.DeptPath,
					LeaderID:  dept.LeaderID,
					ParentID:  dept.ParentID,
					SortOrder: dept.SortOrder,
					Status:    dept.Status,
				},
				OperatorBaseInfo: gobject.OperatorBaseInfo{
					UpdatedAt: dept.UpdatedAt.Unix(),
				},
				Children: buildTree(dept.ID),
			}
			nodes = append(nodes, node)
		}
		return nodes
	}

	tree := buildTree(parentID)
	return &dtoorg.DepartmentTreeResp{
		List: tree,
	}, nil
}
