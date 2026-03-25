package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/tenant"
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
}

type departmentSvc struct {
}

var _ DepartmentSvc = (*departmentSvc)(nil)

func NewDepartmentSvc() DepartmentSvc {
	return &departmentSvc{}
}

// Create 创建部门管理
func (svc *departmentSvc) Create(ctx *gin.Context, req *dtoorg.DepartmentCreateReq) (*dtoorg.DepartmentCreateResp, error) {
	insertEntity := &iammodel.DepartmentEntity{
		DeptCode:  req.DeptCode,
		DeptLevel: req.DeptLevel,
		DeptName:  req.DeptName,
		DeptPath:  req.DeptPath,
		LeaderID:  req.LeaderID,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	if err := iamdao.NewDepartmentDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentCreate] daoDepartment Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentCreateError)
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
	if err = tenant.CheckCompanyAccess(ctx, departmentEntity.CompanyID); err != nil {
		return err
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
		return code.GetError(code.DepartmentNotExistError)
	}
	if err = tenant.CheckCompanyAccess(ctx, departmentEntity.CompanyID); err != nil {
		return err
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
	if err = tenant.CheckCompanyAccess(ctx, departmentEntity.CompanyID); err != nil {
		return nil, err
	}
	resp := &dtoorg.DepartmentDetailResp{
		ID: departmentEntity.ID,
		DepartmentBaseInfo: objorg.DepartmentBaseInfo{
			CompanyID: departmentEntity.CompanyID,
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
	cond := &iamdao.DepartmentCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
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
				CompanyID: v.CompanyID,
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
