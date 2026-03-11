package svcorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/genericdao"
	"github.com/morehao/golib/biz/gcontext/gincontext"
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
		CompanyID: req.CompanyID,
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

	if err := iamdao.NewDepartmentDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcorg.Delete] daoDepartment Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentDeleteError)
	}
	return nil
}

// Update 更新部门管理
func (svc *departmentSvc) Update(ctx *gin.Context, req *dtoorg.DepartmentUpdateReq) error {

	updateEntity := &iammodel.DepartmentEntity{
		CompanyID: req.CompanyID,
		DeptCode:  req.DeptCode,
		DeptLevel: req.DeptLevel,
		DeptName:  req.DeptName,
		DeptPath:  req.DeptPath,
		LeaderID:  req.LeaderID,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}
	if err := iamdao.NewDepartmentDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentUpdate] daoDepartment UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.DepartmentUpdateError)
	}
	return nil
}

// Detail 根据id获取部门管理
func (svc *departmentSvc) Detail(ctx *gin.Context, req *dtoorg.DepartmentDetailReq) (*dtoorg.DepartmentDetailResp, error) {
	detailEntity, err := iamdao.NewDepartmentDao().GetById(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentDetail] daoDepartment GetById fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentGetDetailError)
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.DepartmentNotExistError)
	}
	resp := &dtoorg.DepartmentDetailResp{
		ID: detailEntity.ID,
		DepartmentBaseInfo: objorg.DepartmentBaseInfo{
			CompanyID: detailEntity.CompanyID,
			DeptCode:  detailEntity.DeptCode,
			DeptLevel: detailEntity.DeptLevel,
			DeptName:  detailEntity.DeptName,
			DeptPath:  detailEntity.DeptPath,
			LeaderID:  detailEntity.LeaderID,
			ParentID:  detailEntity.ParentID,
			SortOrder: detailEntity.SortOrder,
			Status:    detailEntity.Status,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
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
	dataList, total, err := iamdao.NewDepartmentDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcorg.DepartmentPageList] daoDepartment GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.DepartmentGetPageListError)
	}
	list := make([]dtoorg.DepartmentPageListItem, 0, len(dataList))
	for _, v := range dataList {
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
