package svcuser

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demo/demodao"
	"github.com/morehao/goark/apps/demo/demomodel"
	"github.com/morehao/goark/apps/demo/internal/dto/dtouser"
	"github.com/morehao/goark/apps/demo/object/objuser"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type UserSvc interface {
	Create(ctx *gin.Context, req *dtouser.UserCreateReq) (*dtouser.UserCreateResp, error)
	Delete(ctx *gin.Context, req *dtouser.UserDeleteReq) error
	Update(ctx *gin.Context, req *dtouser.UserUpdateReq) error
	Detail(ctx *gin.Context, req *dtouser.UserDetailReq) (*dtouser.UserDetailResp, error)
	PageList(ctx *gin.Context, req *dtouser.UserPageListReq) (*dtouser.UserPageListResp, error)
}

type userSvc struct {
}

var _ UserSvc = (*userSvc)(nil)

func NewUserSvc() UserSvc {
	return &userSvc{}
}

// Create 创建用户管理
func (svc *userSvc) Create(ctx *gin.Context, req *dtouser.UserCreateReq) (*dtouser.UserCreateResp, error) {
	insertEntity := &demomodel.UserEntity{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
	}

	if err := demodao.NewUserDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.UserCreate] demodao Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserCreateError)
	}
	return &dtouser.UserCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除用户管理
func (svc *userSvc) Delete(ctx *gin.Context, req *dtouser.UserDeleteReq) error {
	userID := gincontext.GetUserID(ctx)

	if err := demodao.NewUserDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] demodao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	return nil
}

// Update 更新用户管理
func (svc *userSvc) Update(ctx *gin.Context, req *dtouser.UserUpdateReq) error {

	updateEntity := &demomodel.UserEntity{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
	}
	if err := demodao.NewUserDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] demodao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	return nil
}

// Detail 根据id获取用户管理
func (svc *userSvc) Detail(ctx *gin.Context, req *dtouser.UserDetailReq) (*dtouser.UserDetailResp, error) {

	detailEntity, err := demodao.NewUserDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserDetail] demodao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetDetailError)
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}
	resp := &dtouser.UserDetailResp{
		ID: detailEntity.ID,
		UserBaseInfo: objuser.UserBaseInfo{
			CompanyID:    detailEntity.CompanyID,
			DepartmentID: detailEntity.DepartmentID,
			Name:         detailEntity.Name,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取用户管理列表
func (svc *userSvc) PageList(ctx *gin.Context, req *dtouser.UserPageListReq) (*dtouser.UserPageListResp, error) {
	cond := &demodao.UserCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	dataList, total, err := demodao.NewUserDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserPageList] demodao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetPageListError)
	}
	list := make([]dtouser.UserPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtouser.UserPageListItem{
			ID: v.ID,
			UserBaseInfo: objuser.UserBaseInfo{
				CompanyID:    v.CompanyID,
				DepartmentID: v.DepartmentID,
				Name:         v.Name,
			},
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtouser.UserPageListResp{
		List:  list,
		Total: total,
	}, nil
}
