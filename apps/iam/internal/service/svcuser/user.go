package svcuser

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/user"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
	"github.com/morehao/goark/apps/iam/model"
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
	AssignDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsAssignReq) error
	ListDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsReq) (*dtouser.UserDepartmentsResp, error)
	AssignRoles(ctx *gin.Context, req *dtouser.UserAssignRolesReq) error
	ListRoles(ctx *gin.Context, req *dtouser.UserRolesReq) (*dtouser.UserRolesResp, error)
}

type userSvc struct {
}

var _ UserSvc = (*userSvc)(nil)

func NewUserSvc() UserSvc {
	return &userSvc{}
}

// Create 创建用户管理
func (svc *userSvc) Create(ctx *gin.Context, req *dtouser.UserCreateReq) (*dtouser.UserCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
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
		UserType:    model.UserType(req.UserType),
		Status:      model.UserStatus(req.Status),
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
	userEntity, err := dao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
	}

	if err = dao.NewUserDao().Delete(ctx, req.UserID, userID); err != nil {
		glog.Errorf(ctx, "[svcuser.Delete] daoUser Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserDeleteError)
	}
	return nil
}

// Update 更新用户管理
func (svc *userSvc) Update(ctx *gin.Context, req *dtouser.UserUpdateReq) error {
	userEntity, err := dao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
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
	if err = dao.NewUserDao().UpdateMap(ctx, req.UserID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcuser.UserUpdate] daoUser UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.UserUpdateError)
	}
	return nil
}

// Detail 根据id获取用户管理
func (svc *userSvc) Detail(ctx *gin.Context, req *dtouser.UserDetailReq) (*dtouser.UserDetailResp, error) {
	userEntity, err := dao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserDetail] daoUser GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetDetailError)
	}
	// 判断是否存在
	if userEntity == nil || userEntity.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}
	resp := &dtouser.UserDetailResp{
		UserID: userEntity.ID,
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
			Status:      string(userEntity.Status),
			UserType:    string(userEntity.UserType),
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
	glog.Infof(ctx, "[svcuser.UserPageList] req:%s", gutil.ToJsonString(req))
	cond := &dao.UserCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	userEntityList, total, err := dao.NewUserDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.UserPageList] daoUser GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.UserGetPageListError)
	}
	list := make([]dtouser.UserPageListItem, 0, len(userEntityList))
	for _, v := range userEntityList {
		list = append(list, dtouser.UserPageListItem{
			UserID: v.ID,
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
				Status:      string(v.Status),
				UserType:    string(v.UserType),
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
	deptEntity, err := dao.NewDepartmentDao().GetByCond(ctx, &dao.DepartmentCond{
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
	userEntity, err := dao.NewUserDao().GetByCond(ctx, &dao.UserCond{
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
	userDeptDao := dao.NewUserDepartmentDao().WithTx(tx)

	primaryDeptEntity := &model.UserDepartmentEntity{
		TenantID:  tenantID,
		UserID:    userID,
		DeptID:    primaryDeptID,
		DeptType:  model.UserDeptTypePrimary,
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
		secondaryDeptEntity := &model.UserDepartmentEntity{
			TenantID:  tenantID,
			UserID:    userID,
			DeptID:    deptID,
			DeptType:  model.UserDeptTypeSecondary,
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

func (svc *userSvc) AssignDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsAssignReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	userEntity, err := dao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignDepartments] user not found, userID:%d", req.UserID)
		return code.GetError(code.UserNotExistError)
	}

	primaryDeptEntity, err := dao.NewDepartmentDao().GetByID(ctx, req.PrimaryDeptID)
	if err != nil || primaryDeptEntity == nil || primaryDeptEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignDepartments] primary department not found, primaryDeptID:%d", req.PrimaryDeptID)
		return code.GetError(code.DepartmentNotExistError)
	}

	if primaryDeptEntity.TenantID != tenantID {
		return code.GetError(code.TenantScopeForbiddenError)
	}

	secondaryDeptIDs := make([]uint, 0)
	for _, deptID := range req.SecondaryDeptIDs {
		if deptID != req.PrimaryDeptID {
			secondaryDeptIDs = append(secondaryDeptIDs, deptID)
		}
	}

	for _, deptID := range secondaryDeptIDs {
		deptEntity, err := dao.NewDepartmentDao().GetByID(ctx, deptID)
		if err != nil || deptEntity == nil || deptEntity.ID == 0 {
			glog.Errorf(ctx, "[svcuser.AssignDepartments] secondary department not found, deptID:%d", deptID)
			return code.GetError(code.DepartmentNotExistError)
		}
		if deptEntity.TenantID != tenantID {
			return code.GetError(code.TenantScopeForbiddenError)
		}
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		userDeptDao := dao.NewUserDepartmentDao().WithTx(tx)

		existingList, err := userDeptDao.GetListByCond(ctx, &dao.UserDepartmentCond{
			UserID:   req.UserID,
			TenantID: tenantID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcuser.AssignDepartments] GetListByCond fail, err:%v, userID:%d", err, req.UserID)
			return code.GetError(code.UserUpdateError)
		}

		for _, existing := range existingList {
			if err := userDeptDao.Delete(ctx, existing.ID, operatorID); err != nil {
				glog.Errorf(ctx, "[svcuser.AssignDepartments] Delete fail, err:%v, id:%d", err, existing.ID)
				return code.GetError(code.UserUpdateError)
			}
		}

		primaryDeptEntity := &model.UserDepartmentEntity{
			TenantID:  tenantID,
			UserID:    req.UserID,
			DeptID:    req.PrimaryDeptID,
			DeptType:  model.UserDeptTypePrimary,
			CreatedBy: operatorID,
			UpdatedBy: operatorID,
		}
		if err := userDeptDao.Insert(ctx, primaryDeptEntity); err != nil {
			glog.Errorf(ctx, "[svcuser.AssignDepartments] Insert primary dept fail, err:%v, userID:%d, deptID:%d", err, req.UserID, req.PrimaryDeptID)
			return code.GetError(code.UserUpdateError)
		}

		for _, deptID := range secondaryDeptIDs {
			secondaryDeptEntity := &model.UserDepartmentEntity{
				TenantID:  tenantID,
				UserID:    req.UserID,
				DeptID:    deptID,
				DeptType:  model.UserDeptTypeSecondary,
				CreatedBy: operatorID,
				UpdatedBy: operatorID,
			}
			if err := userDeptDao.Insert(ctx, secondaryDeptEntity); err != nil {
				glog.Errorf(ctx, "[svcuser.AssignDepartments] Insert secondary dept fail, err:%v, userID:%d, deptID:%d", err, req.UserID, deptID)
				return code.GetError(code.UserUpdateError)
			}
		}

		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcuser.AssignDepartments] Transaction fail, err:%v", txErr)
		return code.GetError(code.UserUpdateError)
	}

	return nil
}

func (svc *userSvc) ListDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsReq) (*dtouser.UserDepartmentsResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	userDeptDao := dao.NewUserDepartmentDao()
	cond := &dao.UserDepartmentCond{
		UserID:   req.UserID,
		TenantID: tenantID,
	}
	userDepts, _, err := userDeptDao.GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.ListDepartments] GetPageListByCond fail, err:%v", err)
		return nil, code.GetError(code.UserGetDetailError)
	}

	list := make([]dtouser.UserDepartmentItem, 0, len(userDepts))
	for _, ud := range userDepts {
		deptEntity, err := dao.NewDepartmentDao().GetByID(ctx, ud.DeptID)
		if err != nil || deptEntity == nil {
			continue
		}
		list = append(list, dtouser.UserDepartmentItem{
			DepartmentID:   ud.DeptID,
			DepartmentName: deptEntity.DeptName,
			DeptType:       string(ud.DeptType),
		})
	}

	return &dtouser.UserDepartmentsResp{
		List: list,
	}, nil
}

// AssignRoles 分配用户角色(全量替换)
func (svc *userSvc) AssignRoles(ctx *gin.Context, req *dtouser.UserAssignRolesReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	// 检查用户是否存在
	userEntity, err := dao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignRoles] user not found, userID:%d", req.UserID)
		return code.GetError(code.UserNotExistError)
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		userRoleDao := dao.NewUserRoleDao().WithTx(tx)

		// 删除该用户的所有已有角色关联
		existingList, err := userRoleDao.GetListByCond(ctx, &dao.UserRoleCond{
			UserID:   req.UserID,
			TenantID: tenantID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcuser.AssignRoles] GetListByCond fail, err:%v, userID:%d", err, req.UserID)
			return code.GetError(code.UserUpdateError)
		}
		for _, existing := range existingList {
			if err := userRoleDao.Delete(ctx, existing.ID, operatorID); err != nil {
				glog.Errorf(ctx, "[svcuser.AssignRoles] Delete fail, err:%v, id:%d", err, existing.ID)
				return code.GetError(code.UserUpdateError)
			}
		}

		// 批量插入新的角色关联
		for _, roleID := range req.RoleIDs {
			entity := &model.UserRoleEntity{
				TenantID:  tenantID,
				UserID:    req.UserID,
				RoleID:    roleID,
				CreatedBy: operatorID,
				UpdatedBy: operatorID,
			}
			if err := userRoleDao.Insert(ctx, entity); err != nil {
				glog.Errorf(ctx, "[svcuser.AssignRoles] Insert fail, err:%v, userID:%d, roleID:%d", err, req.UserID, roleID)
				return code.GetError(code.UserUpdateError)
			}
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcuser.AssignRoles] Transaction fail, err:%v", txErr)
		return code.GetError(code.UserUpdateError)
	}
	return nil
}

func (svc *userSvc) ListRoles(ctx *gin.Context, req *dtouser.UserRolesReq) (*dtouser.UserRolesResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	userRoleList, err := dao.NewUserRoleDao().GetListByCond(ctx, &dao.UserRoleCond{
		UserID:   req.UserID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.ListRoles] GetListByCond fail, err:%v, userID:%d", err, req.UserID)
		return nil, code.GetError(code.UserGetDetailError)
	}

	list := make([]dtouser.UserRoleItem, 0, len(userRoleList))
	for _, ur := range userRoleList {
		roleEntity, err := dao.NewRoleDao().GetByID(ctx, ur.RoleID)
		if err != nil || roleEntity == nil || roleEntity.ID == 0 {
			continue
		}
		list = append(list, dtouser.UserRoleItem{
			RoleID:   roleEntity.ID,
			RoleName: roleEntity.RoleName,
			RoleCode: roleEntity.RoleCode,
			RoleType: string(roleEntity.RoleType),
		})
	}

	return &dtouser.UserRolesResp{
		List: list,
	}, nil
}
