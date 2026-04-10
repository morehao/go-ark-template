package svcuser

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/core/user"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
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
	AssignDepartment(ctx *gin.Context, req *dtouser.UserDepartmentAssignReq) error
	RemoveDepartment(ctx *gin.Context, req *dtouser.UserDepartmentRemoveReq) error
	ListDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsReq) (*dtouser.UserDepartmentsResp, error)
	AssignRoles(ctx *gin.Context, req *dtouser.UserAssignRolesReq) error
	RemoveRoles(ctx *gin.Context, req *dtouser.UserRemoveRolesReq) error
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
		UserType:    iammodel.UserType(req.UserType),
		Status:      iammodel.UserStatus(req.Status),
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

func (svc *userSvc) AssignDepartment(ctx *gin.Context, req *dtouser.UserDepartmentAssignReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	userEntity, err := iamdao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignDepartment] user not found, userID:%d", req.UserID)
		return code.GetError(code.UserNotExistError)
	}

	deptEntity, err := iamdao.NewDepartmentDao().GetByID(ctx, req.DepartmentID)
	if err != nil || deptEntity == nil || deptEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignDepartment] department not found, departmentID:%d", req.DepartmentID)
		return code.GetError(code.DepartmentNotExistError)
	}

	if deptEntity.TenantID != tenantID {
		return code.GetError(code.TenantScopeForbiddenError)
	}

	userDeptDao := iamdao.NewUserDepartmentDao()
	cond := &iamdao.UserDepartmentCond{
		UserID:   req.UserID,
		DeptID:   req.DepartmentID,
		TenantID: tenantID,
	}
	existingDepts, _, err := userDeptDao.GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.AssignDepartment] GetPageListByCond fail, err:%v", err)
		return code.GetError(code.UserUpdateError)
	}

	if len(existingDepts) > 0 {
		return nil
	}

	deptType := req.DeptType
	if deptType == iammodel.UserDeptTypePrimary {
		primaryCond := &iamdao.UserDepartmentCond{
			UserID:   req.UserID,
			TenantID: tenantID,
			DeptType: iammodel.UserDeptTypePrimary,
		}
		primaryDepts, _, _ := userDeptDao.GetPageListByCond(ctx, primaryCond)
		for _, pd := range primaryDepts {
			updateMap := map[string]any{
				"dept_type": iammodel.UserDeptTypeSecondary,
			}
			if err := userDeptDao.UpdateMap(ctx, pd.ID, updateMap); err != nil {
				glog.Errorf(ctx, "[svcuser.AssignDepartment] update primary to secondary fail, err:%v", err)
			}
		}
	}

	userDeptEntity := &iammodel.UserDepartmentEntity{
		TenantID:  tenantID,
		UserID:    req.UserID,
		DeptID:    req.DepartmentID,
		DeptType:  deptType,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if err := userDeptDao.Insert(ctx, userDeptEntity); err != nil {
		glog.Errorf(ctx, "[svcuser.AssignDepartment] Insert fail, err:%v", err)
		return code.GetError(code.UserUpdateError)
	}

	return nil
}

func (svc *userSvc) RemoveDepartment(ctx *gin.Context, req *dtouser.UserDepartmentRemoveReq) error {
	tenantID := gincontext.GetTenantID(ctx)

	userDeptDao := iamdao.NewUserDepartmentDao()
	cond := &iamdao.UserDepartmentCond{
		UserID:   req.UserID,
		DeptID:   req.DepartmentID,
		TenantID: tenantID,
	}
	existingDepts, _, err := userDeptDao.GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.RemoveDepartment] GetPageListByCond fail, err:%v", err)
		return code.GetError(code.UserUpdateError)
	}

	if len(existingDepts) == 0 {
		return nil
	}

	userID := gincontext.GetUserID(ctx)
	for _, dept := range existingDepts {
		if err := userDeptDao.Delete(ctx, dept.ID, userID); err != nil {
			glog.Errorf(ctx, "[svcuser.RemoveDepartment] Delete fail, err:%v", err)
			return code.GetError(code.UserUpdateError)
		}
	}

	return nil
}

func (svc *userSvc) ListDepartments(ctx *gin.Context, req *dtouser.UserDepartmentsReq) (*dtouser.UserDepartmentsResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	userDeptDao := iamdao.NewUserDepartmentDao()
	cond := &iamdao.UserDepartmentCond{
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
		deptEntity, err := iamdao.NewDepartmentDao().GetByID(ctx, ud.DeptID)
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
	userEntity, err := iamdao.NewUserDao().GetByID(ctx, req.UserID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.AssignRoles] user not found, userID:%d", req.UserID)
		return code.GetError(code.UserNotExistError)
	}

	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		userRoleDao := iamdao.NewUserRoleDao().WithTx(tx)

		// 删除该用户的所有已有角色关联
		existingList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
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
			entity := &iammodel.UserRoleEntity{
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

// RemoveRoles 移除用户角色
func (svc *userSvc) RemoveRoles(ctx *gin.Context, req *dtouser.UserRemoveRolesReq) error {
	operatorID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	userRoleDao := iamdao.NewUserRoleDao()
	for _, roleID := range req.RoleIDs {
		existingList, err := userRoleDao.GetListByCond(ctx, &iamdao.UserRoleCond{
			UserID:   req.UserID,
			RoleID:   roleID,
			TenantID: tenantID,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcuser.RemoveRoles] GetListByCond fail, err:%v, userID:%d, roleID:%d", err, req.UserID, roleID)
			return code.GetError(code.UserUpdateError)
		}
		for _, existing := range existingList {
			if err := userRoleDao.Delete(ctx, existing.ID, operatorID); err != nil {
				glog.Errorf(ctx, "[svcuser.RemoveRoles] Delete fail, err:%v, id:%d", err, existing.ID)
				return code.GetError(code.UserUpdateError)
			}
		}
	}
	return nil
}

// ListRoles 查询用户角色列表
func (svc *userSvc) ListRoles(ctx *gin.Context, req *dtouser.UserRolesReq) (*dtouser.UserRolesResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	userRoleList, err := iamdao.NewUserRoleDao().GetListByCond(ctx, &iamdao.UserRoleCond{
		UserID:   req.UserID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.ListRoles] GetListByCond fail, err:%v, userID:%d", err, req.UserID)
		return nil, code.GetError(code.UserGetDetailError)
	}

	list := make([]dtouser.UserRoleItem, 0, len(userRoleList))
	for _, ur := range userRoleList {
		roleEntity, err := iamdao.NewRoleDao().GetByID(ctx, ur.RoleID)
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
