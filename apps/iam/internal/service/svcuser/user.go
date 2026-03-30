package svcuser

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/organization"
	"github.com/morehao/goark/apps/iam/core/user"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
	"github.com/morehao/goark/apps/iam/internal/tenantctx"
	"github.com/morehao/goark/apps/iam/object/objuser"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
	"gorm.io/gorm"
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
	tenantID := tenantctx.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	primaryDeptID, err := svc.getOrCreatePrimaryDeptID(ctx, tenantID, req.PrimaryDeptID)
	if err != nil {
		return nil, err
	}

	if err := svc.checkUsernameUnique(ctx, tenantID, req.Username); err != nil {
		return nil, err
	}

	params := &user.CreatePersonParams{
		Mobile:      strings.TrimSpace(req.Mobile),
		Email:       strings.TrimSpace(req.Email),
		RealName:    req.RealName,
		OperatorID:  operatorID,
		TenantID:    tenantID,
		DeptID:      primaryDeptID,
		Username:    req.Username,
		UserType:    req.UserType,
		Status:      req.Status,
		EmployeeNo:  req.EmployeeNo,
		JobLevel:    req.JobLevel,
		Position:    req.Position,
		LastLoginIp: req.LastLoginIp,
		LoginCount:  req.LoginCount,
	}

	var result *user.CreatePersonResult
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = user.CreatePersonWithUser(ctx, tx, params)
		if err != nil {
			return err
		}
		if err := svc.createUserDeptRelations(ctx, tx, tenantID, result.UserID, primaryDeptID, req.SecondaryDeptIDs, operatorID); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcuser.Create] Transaction fail, err:%v", txErr)
		return nil, code.GetError(code.UserCreateError)
	}

	return &dtouser.UserCreateResp{
		UserID:   result.UserID,
		PersonID: result.PersonID,
	}, nil
}

// Delete 删除用户管理
func (svc *userSvc) Delete(ctx *gin.Context, req *dtouser.UserDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	userEntity, err := iamdao.NewUserDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, userEntity.TenantID); err != nil {
		return err
	}

	if err = iamdao.NewUserDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] daoUser Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	return nil
}

// Update 更新用户管理
func (svc *userSvc) Update(ctx *gin.Context, req *dtouser.UserUpdateReq) error {
	userEntity, err := iamdao.NewUserDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, userEntity.TenantID); err != nil {
		return err
	}
	updateMap := map[string]any{
		"dept_id":       req.DeptID,
		"employee_no":   req.EmployeeNo,
		"job_level":     req.JobLevel,
		"last_login_ip": req.LastLoginIp,
		"login_count":   req.LoginCount,
		"person_id":     req.PersonID,
		"position":      req.Position,
		"status":        req.Status,
		"user_type":     req.UserType,
		"username":      req.Username,
	}
	if err = iamdao.NewUserDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] daoUser UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	return nil
}

// Detail 根据id获取用户管理
func (svc *userSvc) Detail(ctx *gin.Context, req *dtouser.UserDetailReq) (*dtouser.UserDetailResp, error) {
	userEntity, err := iamdao.NewUserDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserDetail] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetDetailError)
	}
	// 判断是否存在
	if userEntity == nil || userEntity.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}
	if err = organization.CheckTenantAccess(ctx, userEntity.TenantID); err != nil {
		return nil, err
	}
	resp := &dtouser.UserDetailResp{
		ID: userEntity.ID,
		UserBaseInfo: objuser.UserBaseInfo{
			TenantID:    userEntity.TenantID,
			DeptID:      userEntity.DeptID,
			EmployeeNo:  userEntity.EmployeeNo,
			EntryDate:   userEntity.EntryDate.Unix(),
			JobLevel:    userEntity.JobLevel,
			LastLoginAt: userEntity.LastLoginAt.Unix(),
			LastLoginIp: userEntity.LastLoginIp,
			LoginCount:  userEntity.LoginCount,
			PersonID:    userEntity.PersonID,
			Position:    userEntity.Position,
			Status:      userEntity.Status,
			UserType:    userEntity.UserType,
			Username:    userEntity.Username,
		},
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: userEntity.CreatedAt.Unix(),
			UpdatedAt: userEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

// PageList 分页获取用户管理列表
func (svc *userSvc) PageList(ctx *gin.Context, req *dtouser.UserPageListReq) (*dtouser.UserPageListResp, error) {
	cond := &iamdao.UserCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	userEntityList, total, err := iamdao.NewUserDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserPageList] daoUser GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetPageListError)
	}
	list := make([]dtouser.UserPageListItem, 0, len(userEntityList))
	for _, v := range userEntityList {
		list = append(list, dtouser.UserPageListItem{
			ID: v.ID,
			UserBaseInfo: objuser.UserBaseInfo{
				TenantID:    v.TenantID,
				DeptID:      v.DeptID,
				EmployeeNo:  v.EmployeeNo,
				EntryDate:   v.EntryDate.Unix(),
				JobLevel:    v.JobLevel,
				LastLoginAt: v.LastLoginAt.Unix(),
				LastLoginIp: v.LastLoginIp,
				LoginCount:  v.LoginCount,
				PersonID:    v.PersonID,
				Position:    v.Position,
				Status:      v.Status,
				UserType:    v.UserType,
				Username:    v.Username,
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

func (svc *userSvc) getOrCreatePrimaryDeptID(ctx *gin.Context, tenantID uint, primaryDeptID uint) (uint, error) {
	if primaryDeptID > 0 {
		return primaryDeptID, nil
	}
	deptEntity, err := iamdao.NewDepartmentDao().GetByCond(ctx, &iamdao.DepartmentCond{
		TenantID:    tenantID,
		ParentID:    0,
		ParentIDNil: true,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getOrCreatePrimaryDeptID] GetByCond fail, err:%v, tenantID:%d", err, tenantID)
		return 0, code.GetError(code.UserCreateError)
	}
	if deptEntity == nil || deptEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.getOrCreatePrimaryDeptID] root dept not found, tenantID:%d", tenantID)
		return 0, code.GetError(code.UserCreateError)
	}
	return deptEntity.ID, nil
}

func (svc *userSvc) checkUsernameUnique(ctx *gin.Context, tenantID uint, username string) error {
	if username == "" {
		return nil
	}
	userEntity, err := iamdao.NewUserDao().GetByCond(ctx, &iamdao.UserCond{
		TenantID: tenantID,
		Username: username,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.checkUsernameUnique] GetByCond fail, err:%v, tenantID:%d, username:%s", err, tenantID, username)
		return code.GetError(code.UserCreateError)
	}
	if userEntity != nil && userEntity.ID > 0 {
		return code.GetError(code.UsernameDuplicateError)
	}
	return nil
}

func (svc *userSvc) createUserDeptRelations(ctx *gin.Context, tx *gorm.DB, tenantID uint, userID uint, primaryDeptID uint, secondaryDeptIDs []uint, operatorID uint) error {
	userDeptDao := iamdao.NewUserDepartmentDao().WithTx(tx)

	primaryDeptEntity := &iammodel.UserDepartmentEntity{
		TenantID:  tenantID,
		UserID:    userID,
		DeptID:    primaryDeptID,
		DeptType:  iammodel.UserDeptTypePrimary,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if err := userDeptDao.Insert(ctx, primaryDeptEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.createUserDeptRelations] Insert primary dept fail, err:%v, userID:%d, deptID:%d", err, userID, primaryDeptID)
		return code.GetError(code.UserCreateError)
	}

	for _, deptID := range secondaryDeptIDs {
		if deptID == primaryDeptID {
			continue
		}
		secondaryDeptEntity := &iammodel.UserDepartmentEntity{
			TenantID:  tenantID,
			UserID:    userID,
			DeptID:    deptID,
			DeptType:  iammodel.UserDeptTypeSecondary,
			CreatedBy: operatorID,
			UpdatedBy: operatorID,
		}
		if err := userDeptDao.Insert(ctx, secondaryDeptEntity); err != nil {
			glog.Errorf(ctx, "[svcuser.createUserDeptRelations] Insert secondary dept fail, err:%v, userID:%d, deptID:%d", err, userID, deptID)
			return code.GetError(code.UserCreateError)
		}
	}
	return nil
}
