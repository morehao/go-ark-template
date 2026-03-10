package svcuser

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demoapp/dao/daouser"
	"github.com/morehao/goark/apps/demoapp/internal/dto/dtouser"
	"github.com/morehao/goark/apps/demoapp/model"
	"github.com/morehao/goark/apps/demoapp/object/objuser"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/dbaccess/dbes"
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
	userID := gincontext.GetUserID(ctx)
	insertEntity := &model.UserEntity{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	if err := daouser.NewUserDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.UserCreate] daoUser Create fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserCreateError)
	}
	return &dtouser.UserCreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除用户管理
func (svc *userSvc) Delete(ctx *gin.Context, req *dtouser.UserDeleteReq) error {
	userID := gincontext.GetUserID(ctx)

	if err := daouser.NewUserDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] daoUser Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	return nil
}

// Update 更新用户管理
func (svc *userSvc) Update(ctx *gin.Context, req *dtouser.UserUpdateReq) error {
	userID := gincontext.GetUserID(ctx)

	updateEntity := &model.UserEntity{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		UpdatedBy:    userID,
	}
	if err := daouser.NewUserDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] daoUser UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	return nil
}

// Detail 根据id获取用户管理
func (svc *userSvc) Detail(ctx *gin.Context, req *dtouser.UserDetailReq) (*dtouser.UserDetailResp, error) {
	detailEntity, err := daouser.NewUserDao().GetById(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserDetail] daoUser GetById fail, err:%v, req:%s", err, gutil.ToJsonString(req))
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
			CreatedBy: detailEntity.CreatedBy,
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedBy: detailEntity.UpdatedBy,
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取用户管理列表
func (svc *userSvc) PageList(ctx *gin.Context, req *dtouser.UserPageListReq) (*dtouser.UserPageListResp, error) {
	// Redis 调用示例：获取用户列表缓存
	cacheKey := "user:list:cache"
	_, redisErr := dbclient.RedisCli.Get(ctx, cacheKey).Result()
	if redisErr != nil {
		glog.Debugf(ctx, "[svcuser.UserPageList] redis get fail, key:%s, err:%v", cacheKey, redisErr)
	}

	// ES 调用示例：搜索用户
	query := dbes.NewBuilder().SetQuery(dbes.BuildMap("match_all", dbes.Map{})).Build()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err == nil {
		_, esErr := dbclient.DemoES.Search(
			dbclient.DemoES.Search.WithContext(ctx),
			dbclient.DemoES.Search.WithIndex("users"),
			dbclient.DemoES.Search.WithBody(&buf),
		)
		if esErr != nil {
			glog.Debugf(ctx, "[svcuser.UserPageList] es search fail, err:%v", esErr)
		}
	}

	cond := &daouser.UserCond{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	dataList, total, err := daouser.NewUserDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserPageList] daoUser GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
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
				CreatedBy: v.CreatedBy,
				CreatedAt: v.CreatedAt.Unix(),
				UpdatedBy: v.UpdatedBy,
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtouser.UserPageListResp{
		List:  list,
		Total: total,
	}, nil
}
