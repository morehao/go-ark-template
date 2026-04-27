package svcauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/internal/service/svcuser"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/pkg/token"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/gauth/jwtauth"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
	"gorm.io/gorm"
)

const (
	tempTokenExpireDuration = 10 * time.Minute
	tokenIssuer             = "iam"
)

type AuthSvc interface {
	LoginByPassword(ctx *gin.Context, req *dtoauth.LoginByPasswordReq) (*dtoauth.LoginByPasswordResp, error)
	SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error)
	Logout(ctx *gin.Context, refreshToken string) error
	RefreshToken(ctx *gin.Context, req *dtoauth.RefreshTokenReq) (*dtoauth.RefreshTokenResp, error)
	Register(ctx *gin.Context, req *dtoauth.RegisterReq) (*dtoauth.RegisterResp, error)
	UnlockAccount(ctx *gin.Context, req *dtoauth.UnlockAccountReq) error
}

type authSvc struct {
}

var _ AuthSvc = (*authSvc)(nil)

func NewAuthSvc() AuthSvc {
	return &authSvc{}
}

// LoginByPassword 密码登录
func (svc *authSvc) LoginByPassword(ctx *gin.Context, req *dtoauth.LoginByPasswordReq) (*dtoauth.LoginByPasswordResp, error) {
	account := strings.TrimSpace(req.Account)

	orgEntity, err := svc.getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	// 查找自然人(通过手机号或邮箱)
	personEntity, err := svc.findPersonByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	// 查询该自然人关联的所有用户账号
	userList, err := dao.NewUserDao().GetListByCond(ctx, &dao.UserCond{
		PersonID: personEntity.ID,
		Status:   model.UserStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.LoginByPassword] GetListByCond fail, err:%v, personID:%d", err, personEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
	}
	userList, err = svc.filterUsersByOrg(ctx, userList, orgEntity.ID)
	if err != nil {
		return nil, err
	}
	if len(userList) == 0 {
		return nil, code.GetError(code.AuthNoTenantError)
	}

	// 检查用户锁定状态
	if err := svcuser.NewPasswordSvc().CheckUserLockStatus(ctx, userList[0].ID); err != nil {
		return nil, err
	}

	// 验证密码
	if err := gcrypto.ComparePasswordHash(personEntity.PasswordHash, req.Password); err != nil {
		// 记录登录失败
		svcuser.NewPasswordSvc().RecordLoginFail(ctx, userList[0].ID, orgEntity.ID)
		glog.Errorf(ctx, "[svcauth.LoginByPassword] password mismatch, account:%s", account)
		return nil, code.GetError(code.AuthPasswordError)
	}

	// 清除登录失败计数
	svcuser.NewPasswordSvc().ClearLoginFail(ctx, userList[0].ID)

	// 统一返回临时token + 租户列表
	tempToken, err := svc.generateTempToken(personEntity.ID, orgEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantList, err := svc.buildTenantList(ctx, userList)
	if err != nil {
		return nil, err
	}

	return &dtoauth.LoginByPasswordResp{
		TempToken:        tempToken,
		NeedSelectTenant: true,
		TenantList:       tenantList,
		PersonID:         personEntity.ID,
		RealName:         personEntity.RealName,
	}, nil
}

// SelectTenant 选择租户
func (svc *authSvc) SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error) {
	if gincontext.GetUserType(ctx) != "temp" {
		return nil, code.GetError(code.AuthTempTokenRequiredError)
	}

	// 从context获取当前用户信息(临时token中的personID通过UserID字段传递)
	personID := gincontext.GetUserID(ctx)
	orgID := gincontext.GetOrgID(ctx)
	if personID == 0 || orgID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	tenantEntity, err := dao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByID tenant fail, err:%v, tenantID:%d", err, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 || tenantEntity.OrgID != orgID {
		return nil, code.GetError(code.AuthTenantNotInOrgError)
	}

	userEntity, err := dao.NewUserDao().GetByCond(ctx, &dao.UserCond{
		PersonID: personID,
		TenantID: req.TenantID,
		Status:   model.UserStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByCond fail, err:%v, personID:%d, tenantID:%d", err, personID, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	personEntity, err := dao.NewPersonDao().GetByID(ctx, userEntity.PersonID)
	if err != nil || personEntity == nil || personEntity.ID == 0 {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByID person fail, err:%v, personID:%d", err, userEntity.PersonID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	token, refreshToken, err := svc.generateTokenPair(ctx, *userEntity, personEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantEntity, _ = dao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
	tenantName := ""
	if tenantEntity != nil {
		tenantName = tenantEntity.TenantName
	}

	svc.updateLoginInfo(ctx, userEntity)

	return &dtoauth.SelectTenantResp{
		Token:        token,
		RefreshToken: refreshToken,
		UserInfo: dtoauth.LoginUserInfo{
			UserID:     userEntity.ID,
			PersonID:   personEntity.ID,
			Username:   userEntity.Username,
			RealName:   personEntity.RealName,
			UserType:   userEntity.UserType,
			TenantID:   userEntity.TenantID,
			TenantName: tenantName,
		},
	}, nil
}

// Logout 登出(token加黑名单)
func (svc *authSvc) Logout(ctx *gin.Context, refreshToken string) error {
	authToken := ctx.GetHeader("Authorization")
	if authToken != "" {
		if err := token.AddTokenToBlacklist(ctx.Request.Context(), authToken, token.TokenExpireDuration); err != nil {
			return code.GetError(code.AuthLogoutError)
		}
	}

	if refreshToken != "" {
		if err := token.AddRefreshTokenToBlacklist(ctx.Request.Context(), refreshToken); err != nil {
			glog.Errorf(ctx, "[svcauth.Logout] AddRefreshTokenToBlacklist fail, err:%v", err)
		}
	}

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func IsTokenBlacklisted(ctx *gin.Context, tokenStr string) bool {
	return token.IsTokenBlacklisted(ctx.Request.Context(), tokenStr)
}

func (svc *authSvc) findPersonByAccount(ctx *gin.Context, account string) (*model.PersonEntity, error) {
	// 先按手机号查询
	personEntity, err := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{
		Mobile: account,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.findPersonByAccount] GetByCond by mobile fail, err:%v, account:%s", err, account)
		return nil, code.GetError(code.AuthLoginError)
	}
	if personEntity != nil && personEntity.ID > 0 {
		return personEntity, nil
	}

	// 再按邮箱查询
	personEntity, err = dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{
		Email: account,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.findPersonByAccount] GetByCond by email fail, err:%v, account:%s", err, account)
		return nil, code.GetError(code.AuthLoginError)
	}
	if personEntity != nil && personEntity.ID > 0 {
		return personEntity, nil
	}

	return nil, code.GetError(code.AuthPersonNotFoundError)
}

func (svc *authSvc) generateTempToken(personID uint, orgID uint) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	customData := gobject.UserClaims{
		UserID:    personID,
		PersonID:  personID,
		TenantID:  0,
		OrgID:     orgID,
		DeptID:    0,
		RoleIDs:   nil,
		UserType:  "temp",
		TokenType: gobject.TokenTypeTemp,
	}

	token, err := jwtAuth.Issue(
		fmt.Sprintf("%d", personID),
		tokenIssuer,
		time.Now().Add(tempTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return token, nil
}

func (svc *authSvc) generateTokenPair(ctx *gin.Context, userEntity model.UserEntity, personID uint) (token string, refreshToken string, err error) {
	var orgID uint
	if userEntity.TenantID > 0 {
		tenantEntity, _ := dao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if tenantEntity != nil {
			orgID = tenantEntity.OrgID
		}
	}

	deptID, roleIDs := svc.getUserDeptAndRoles(ctx, userEntity.ID, userEntity.TenantID)

	token, err = svc.generateTokenWithOrgID(ctx, userEntity, personID, orgID, deptID, roleIDs)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = svc.generateRefreshToken(ctx, userEntity.ID, personID, userEntity.UserType, userEntity.TenantID, orgID)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (svc *authSvc) generateTokenWithOrgID(ctx *gin.Context, userEntity model.UserEntity, personID uint, orgID uint, deptID uint, roleIDs []uint) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	customData := gobject.UserClaims{
		UserID:    userEntity.ID,
		PersonID:  personID,
		TenantID:  userEntity.TenantID,
		OrgID:     orgID,
		DeptID:    deptID,
		RoleIDs:   roleIDs,
		UserType:  string(userEntity.UserType),
		TokenType: gobject.TokenTypeAuth,
	}

	newToken, err := jwtAuth.Issue(
		fmt.Sprintf("%d", userEntity.ID),
		tokenIssuer,
		time.Now().Add(token.TokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return newToken, nil
}

func (svc *authSvc) generateRefreshToken(ctx *gin.Context, userID uint, personID uint, userType model.UserType, tenantID uint, orgID uint) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	deptID, roleIDs := svc.getUserDeptAndRoles(ctx, userID, tenantID)

	customData := gobject.UserClaims{
		UserID:    userID,
		PersonID:  personID,
		TenantID:  tenantID,
		OrgID:     orgID,
		DeptID:    deptID,
		RoleIDs:   roleIDs,
		UserType:  string(userType),
		TokenType: gobject.TokenTypeRefresh,
	}

	refreshTokenStr, err := jwtAuth.Issue(
		fmt.Sprintf("%d", userID),
		tokenIssuer,
		time.Now().Add(token.RefreshTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return refreshTokenStr, nil
}

func (svc *authSvc) RefreshToken(ctx *gin.Context, req *dtoauth.RefreshTokenReq) (*dtoauth.RefreshTokenResp, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	claims, err := jwtAuth.Parse(req.RefreshToken)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.RefreshToken] Parse refreshToken fail, err:%v", err)
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	if claims.CustomData.UserType != "refresh" {
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	if isRefreshTokenBlacklisted(ctx, req.RefreshToken) {
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	userID := claims.CustomData.UserID
	userEntity, err := dao.NewUserDao().GetByID(ctx, userID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcauth.RefreshToken] GetByID user fail, err:%v, userID:%d", err, userID)
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	if userEntity.Status != model.UserStatusEnabled {
		return nil, code.GetError(code.AuthAccountDisabledError)
	}

	token, refreshToken, err := svc.generateTokenPair(ctx, *userEntity, userEntity.PersonID)
	if err != nil {
		return nil, err
	}

	if err := addRefreshTokenToBlacklist(ctx, req.RefreshToken); err != nil {
		glog.Errorf(ctx, "[svcauth.RefreshToken] Add refreshToken to blacklist fail, err:%v", err)
	}

	return &dtoauth.RefreshTokenResp{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func isRefreshTokenBlacklisted(ctx *gin.Context, tokenStr string) bool {
	if dbclient.RedisCli == nil {
		return false
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		return false
	}
	key := token.RefreshTokenBlacklistKeyPrefix + token.HashToken(tokenStr)
	exists, err := dbclient.RedisCli.Exists(ctx.Request.Context(), key).Result()
	if err != nil {
		glog.Errorf(ctx, "[svcauth.isRefreshTokenBlacklisted] Redis Exists fail, err:%v", err)
		return false
	}
	return exists > 0
}

func addRefreshTokenToBlacklist(ctx *gin.Context, tokenStr string) error {
	return token.AddRefreshTokenToBlacklist(ctx.Request.Context(), tokenStr)
}

func (svc *authSvc) buildTenantList(ctx *gin.Context, userList model.UserEntityList) ([]dtoauth.TenantListItem, error) {
	tenants := make([]dtoauth.TenantListItem, 0, len(userList))
	for _, u := range userList {
		tenantEntity, err := dao.NewTenantDao().GetByID(ctx, u.TenantID)
		if err != nil || tenantEntity == nil || tenantEntity.ID == 0 {
			continue
		}
		orgName := ""
		orgEntity, _ := dao.NewOrganizationDao().GetByID(ctx, tenantEntity.OrgID)
		if orgEntity != nil {
			orgName = orgEntity.OrgName
		}
		tenants = append(tenants, dtoauth.TenantListItem{
			TenantID:   tenantEntity.ID,
			TenantName: tenantEntity.TenantName,
			OrgID:      tenantEntity.OrgID,
			OrgName:    orgName,
		})
	}
	return tenants, nil
}

func (svc *authSvc) getCurrentOrg(ctx *gin.Context) (*model.OrganizationEntity, error) {
	domain := resolveDomain(ctx)
	if domain == "" {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByCond(ctx, &dao.OrganizationCond{
		Domain: domain,
		Status: model.OrgStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getCurrentOrg] daoOrg GetByCond fail, err:%v, domain:%s", err, domain)
		return nil, code.GetError(code.AuthLoginError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}
	return orgEntity, nil
}

func (svc *authSvc) filterUsersByOrg(ctx *gin.Context, userList model.UserEntityList, orgID uint) (model.UserEntityList, error) {
	filtered := make(model.UserEntityList, 0, len(userList))
	for _, userEntity := range userList {
		tenantEntity, err := dao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if err != nil {
			glog.Errorf(ctx, "[svcauth.filterUsersByOrg] daoTenant GetByID fail, err:%v, tenantID:%d", err, userEntity.TenantID)
			return nil, code.GetError(code.AuthLoginError)
		}
		if tenantEntity == nil || tenantEntity.ID == 0 || tenantEntity.Status != model.TenantStatusEnabled {
			continue
		}
		if tenantEntity.OrgID != orgID {
			continue
		}
		filtered = append(filtered, userEntity)
	}
	return filtered, nil
}

func resolveDomain(ctx *gin.Context) string {
	host := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Host"))
	if host == "" && ctx.Request != nil {
		host = strings.TrimSpace(ctx.Request.Host)
	}
	host = strings.TrimSpace(strings.Split(host, ",")[0])
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.Split(host, "/")[0]
	if strings.Count(host, ":") == 1 {
		host = strings.Split(host, ":")[0]
	}
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	return host
}

func resolveDomainFromHost(ctx *gin.Context) string {
	host := ctx.Request.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return strings.TrimSpace(host)
}

func (svc *authSvc) getTenantByDomain(ctx *gin.Context, orgID uint, domain string) (*model.TenantEntity, error) {
	return dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
		OrgID:  orgID,
		Domain: domain,
		Status: model.TenantStatusEnabled,
	})
}

func (svc *authSvc) updateLoginInfo(ctx *gin.Context, userEntity *model.UserEntity) {
	now := time.Now()
	updateMap := map[string]any{
		"last_login_at": now,
		"last_login_ip": gincontext.GetClientIP(ctx),
		"login_count":   userEntity.LoginCount + 1,
	}
	if err := dao.NewUserDao().UpdateMap(ctx, userEntity.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcauth.updateLoginInfo] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
	}
}

func (svc *authSvc) getUserDeptAndRoles(ctx *gin.Context, userID uint, tenantID uint) (deptID uint, roleIDs []uint) {
	deptID = 0
	roleIDs = []uint{}

	userDeptList, err := dao.NewUserDepartmentDao().GetListByCond(ctx, &dao.UserDepartmentCond{
		UserID:   userID,
		TenantID: tenantID,
		DeptType: model.UserDeptTypePrimary,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getUserDeptAndRoles] GetListByCond user department fail, err:%v, userID:%d", err, userID)
	} else if len(userDeptList) > 0 {
		deptID = userDeptList[0].DeptID
	}

	userRoleList, err := dao.NewUserRoleDao().GetListByCond(ctx, &dao.UserRoleCond{
		UserID:   userID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getUserDeptAndRoles] GetListByCond user role fail, err:%v, userID:%d", err, userID)
	} else {
		for _, ur := range userRoleList {
			roleIDs = append(roleIDs, ur.RoleID)
		}
	}

	return deptID, roleIDs
}

func (svc *authSvc) Register(ctx *gin.Context, req *dtoauth.RegisterReq) (*dtoauth.RegisterResp, error) {
	orgEntity, err := svc.getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	registerEnabled, err := svc.getOrgConfigBool(ctx, orgEntity.ID, model.OrgConfigKeyRegisterEnabled)
	if err != nil || !registerEnabled {
		return nil, code.GetError(code.AuthRegisterDisabled)
	}

	domain := resolveDomainFromHost(ctx)

	tenantEntity, err := svc.getTenantByDomain(ctx, orgEntity.ID, domain)
	if err != nil || tenantEntity == nil {
		return nil, code.GetError(code.TenantNotExistError)
	}

	if tenantEntity.Status != model.TenantStatusEnabled {
		return nil, code.GetError(code.AuthRegisterError)
	}

	identityType, _ := svc.getOrgConfigString(ctx, orgEntity.ID, model.OrgConfigKeyRegisterIdentityType)
	if identityType == "" {
		identityType = string(model.RegisterIdentityTypeEmail)
	}
	if err := svc.validateRegisterIdentity(ctx, req, model.RegisterIdentityType(identityType)); err != nil {
		return nil, err
	}

	requireApproval, _ := svc.getOrgConfigBool(ctx, orgEntity.ID, model.OrgConfigKeyRegisterRequireApproval)
	userStatus := model.UserStatusEnabled
	message := "注册成功"
	if requireApproval {
		userStatus = model.UserStatusPending
		message = "注册成功，等待管理员审核"
	}

	passwordHash, err := gcrypto.GeneratePasswordHash(req.Password)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.Register] GeneratePasswordHash fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}

	email := strings.TrimSpace(req.Email)
	personEntity, _ := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{Email: email})
	personExists := personEntity != nil && personEntity.ID > 0

	var userID, personID uint
	txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
		var newPersonID uint
		if personExists {
			newPersonID = personEntity.ID
		} else {
			newPerson := &model.PersonEntity{
				Mobile:       strings.TrimSpace(req.Mobile),
				Email:        email,
				RealName:     req.RealName,
				PasswordHash: passwordHash,
				CreatedBy:    0,
				UpdatedBy:    0,
			}
			if err := dao.NewPersonDao().WithTx(tx).Insert(ctx, newPerson); err != nil {
				return err
			}
			newPersonID = newPerson.ID
		}

		userEntity := &model.UserEntity{
			TenantID:  tenantEntity.ID,
			PersonID:  newPersonID,
			Username:  req.Username,
			UserType:  model.UserTypeNormal,
			Status:    userStatus,
			CreatedBy: 0,
			UpdatedBy: 0,
		}
		if err := dao.NewUserDao().WithTx(tx).Insert(ctx, userEntity); err != nil {
			return err
		}
		userID = userEntity.ID

		return nil
	})
	if txErr != nil {
		glog.Errorf(ctx, "[svcauth.Register] Transaction fail, err:%v", txErr)
		return nil, code.GetError(code.AuthRegisterError)
	}

	return &dtoauth.RegisterResp{
		UserID:       userID,
		PersonID:     personID,
		Status:       string(userStatus),
		PersonExists: personExists,
		Message:      message,
	}, nil
}

func (svc *authSvc) getOrgConfigBool(ctx *gin.Context, orgID uint, configKey string) (bool, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getOrgConfigBool] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return false, err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return false, nil
	}
	return strings.ToLower(configEntity.ConfigValue) == "true", nil
}

func (svc *authSvc) getOrgConfigString(ctx *gin.Context, orgID uint, configKey string) (string, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getOrgConfigString] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return "", err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return "", nil
	}
	return configEntity.ConfigValue, nil
}

func (svc *authSvc) validateRegisterIdentity(ctx *gin.Context, req *dtoauth.RegisterReq, identityType model.RegisterIdentityType) error {
	mobile := strings.TrimSpace(req.Mobile)
	email := strings.TrimSpace(req.Email)

	switch identityType {
	case model.RegisterIdentityTypeMobile:
		if mobile == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
	case model.RegisterIdentityTypeEmail:
		if email == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
	case model.RegisterIdentityTypeBoth:
		if mobile == "" && email == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
	}
	return nil
}

func (svc *authSvc) UnlockAccount(ctx *gin.Context, req *dtoauth.UnlockAccountReq) error {
	account := strings.TrimSpace(req.Account)

	personEntity, err := svc.findPersonByAccount(ctx, account)
	if err != nil {
		return err
	}

	userEntity, err := dao.NewUserDao().GetByCond(ctx, &dao.UserCond{
		PersonID: personEntity.ID,
		Status:   model.UserStatusLocked,
	})
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
	}

	if err := dao.NewUserDao().UpdateMap(ctx, userEntity.ID, map[string]any{
		"status":            model.UserStatusEnabled,
		"login_fail_count":  0,
		"locked_until":      nil,
	}); err != nil {
		glog.Errorf(ctx, "[svcauth.UnlockAccount] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
		return code.GetError(code.UserUpdateError)
	}

	return nil
}
