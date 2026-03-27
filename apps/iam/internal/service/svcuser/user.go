package svcuser

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/core/organization"
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
	"github.com/morehao/golib/gcrypto"
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

	mobile := strings.TrimSpace(req.Mobile)
	email := strings.TrimSpace(req.Email)
	passwordHash := ""
	if mobile != "" {
		passwordHash = svc.generatePassword(mobile)
	}
	personCreateEntity := &iammodel.PersonEntity{
		Mobile:       mobile,
		Email:        email,
		RealName:     req.RealName,
		PasswordHash: passwordHash,
		CreatedBy:    operatorID,
		UpdatedBy:    operatorID,
	}
	userEntity := &iammodel.UserEntity{
		TenantID:    tenantID,
		DeptID:      req.DeptID,
		EmployeeNo:  req.EmployeeNo,
		JobLevel:    req.JobLevel,
		LastLoginIp: req.LastLoginIp,
		LoginCount:  req.LoginCount,
		Position:    req.Position,
		Status:      req.Status,
		UserType:    req.UserType,
		Username:    req.Username,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	var userID uint
	var personID uint
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		personID, err = svc.getOrCreatePersonWithTx(ctx, tx, mobile, email, personCreateEntity)
		if err != nil {
			return err
		}
		userEntity.PersonID = personID
		if err := iamdao.NewUserDao().WithTx(tx).Insert(ctx, userEntity); err != nil {
			return err
		}
		userID = userEntity.ID
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcuser.Create] Transaction fail, err:%v, mobile:%s, email:%s", txErr, mobile, email)
		return nil, code.GetError(code.UserCreateError)
	}

	return &dtouser.UserCreateResp{
		ID:       userID,
		PersonID: personID,
	}, nil
}

func (svc *userSvc) getOrCreatePersonWithTx(ctx *gin.Context, tx *gorm.DB, mobile, email string, personCreateEntity *iammodel.PersonEntity) (uint, error) {
	personEntity, err := iamdao.NewPersonDao().GetByCond(ctx, &iamdao.PersonCond{
		Mobile: mobile,
		Email:  email,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getOrCreatePerson] daoPerson GetByCond fail, err:%v, mobile:%s, email:%s", err, mobile, email)
		return 0, code.GetError(code.UserCreateError)
	}
	if personEntity != nil && personEntity.ID != 0 {
		return personEntity.ID, nil
	}
	if err := iamdao.NewPersonDao().WithTx(tx).Insert(ctx, personCreateEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.getOrCreatePerson] daoPerson Insert fail, err:%v, mobile:%s, email:%s", err, mobile, email)
		return 0, code.GetError(code.UserCreateError)
	}
	return personCreateEntity.ID, nil
}

func (svc *userSvc) generatePassword(mobile string) string {
	prefix := "pwd"
	if config.Conf != nil && config.Conf.Password.Prefix != "" {
		prefix = config.Conf.Password.Prefix
	}
	suffix := mobile
	if len(mobile) > 8 {
		suffix = mobile[len(mobile)-8:]
	}
	hash, _ := gcrypto.GeneratePasswordHash(prefix + suffix)
	return hash
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
