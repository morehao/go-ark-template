package svcuser

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
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
	LoginByPassword(ctx *gin.Context, req *dtouser.LoginByPasswordReq) (*dtouser.LoginByPasswordResp, error)
	SelectTenant(ctx *gin.Context, req *dtouser.SelectTenantReq) (*dtouser.SelectTenantResp, error)
	Logout(ctx *gin.Context, refreshToken string) error
	RefreshToken(ctx *gin.Context, req *dtouser.RefreshTokenReq) (*dtouser.RefreshTokenResp, error)
	Register(ctx *gin.Context, req *dtouser.RegisterReq) (*dtouser.RegisterResp, error)
	UnlockAccount(ctx *gin.Context, req *dtouser.UnlockAccountReq) error
}

type authSvc struct {
}

var _ AuthSvc = (*authSvc)(nil)

func NewAuthSvc() AuthSvc {
	return &authSvc{}
}

func (svc *authSvc) LoginByPassword(ctx *gin.Context, req *dtouser.LoginByPasswordReq) (*dtouser.LoginByPasswordResp, error) {
	account := strings.TrimSpace(req.Account)

	orgEntity, err := svc.getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	personEntity, err := svc.findPersonByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	userList, err := dao.NewUserDao().GetListByCond(ctx, &dao.UserCond{
		PersonID: personEntity.ID,
		Status:   model.UserStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.LoginByPassword] GetListByCond fail, err:%v, personID:%d", err, personEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
	}
	userList, err = svc.filterUsersByOrg(ctx, userList, orgEntity.ID)
	if err != nil {
		return nil, err
	}
	if len(userList) == 0 {
		return nil, code.GetError(code.AuthNoTenantError)
	}

	if err := NewPasswordSvc().CheckUserLockStatus(ctx, userList[0].ID); err != nil {
		return nil, err
	}

	if err := gcrypto.ComparePasswordHash(personEntity.PasswordHash, req.Password); err != nil {
		NewPasswordSvc().RecordLoginFail(ctx, userList[0].ID, orgEntity.ID)
		glog.Errorf(ctx, "[svcuser.LoginByPassword] password mismatch, account:%s", account)
		return nil, code.GetError(code.AuthPasswordError)
	}

	NewPasswordSvc().ClearLoginFail(ctx, userList[0].ID)

	tempToken, err := svc.generateTempToken(personEntity.ID, orgEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantList, err := svc.buildTenantList(ctx, userList)
	if err != nil {
		return nil, err
	}

	return &dtouser.LoginByPasswordResp{
		TempToken:        tempToken,
		NeedSelectTenant: true,
		TenantList:       tenantList,
		PersonID:         personEntity.ID,
		RealName:         personEntity.RealName,
	}, nil
}

func (svc *authSvc) SelectTenant(ctx *gin.Context, req *dtouser.SelectTenantReq) (*dtouser.SelectTenantResp, error) {
	if gincontext.GetUserType(ctx) != "temp" {
		return nil, code.GetError(code.AuthTempTokenRequiredError)
	}

	personID := gincontext.GetUserID(ctx)
	orgID := gincontext.GetOrgID(ctx)
	if personID == 0 || orgID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	tenantEntity, err := dao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.SelectTenant] GetByID tenant fail, err:%v, tenantID:%d", err, req.TenantID)
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
		glog.Errorf(ctx, "[svcuser.SelectTenant] GetByCond fail, err:%v, personID:%d, tenantID:%d", err, personID, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	personEntity, err := dao.NewPersonDao().GetByID(ctx, userEntity.PersonID)
	if err != nil || personEntity == nil || personEntity.ID == 0 {
		glog.Errorf(ctx, "[svcuser.SelectTenant] GetByID person fail, err:%v, personID:%d", err, userEntity.PersonID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	tokenStr, refreshToken, err := svc.generateTokenPair(ctx, *userEntity, personEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantEntity, _ = dao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
	tenantName := ""
	if tenantEntity != nil {
		tenantName = tenantEntity.TenantName
	}

	svc.updateLoginInfo(ctx, userEntity)

	return &dtouser.SelectTenantResp{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		UserInfo: dtouser.LoginUserInfo{
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

func (svc *authSvc) Logout(ctx *gin.Context, refreshToken string) error {
	authToken := ctx.GetHeader("Authorization")
	if authToken != "" {
		if err := token.AddTokenToBlacklist(ctx.Request.Context(), authToken, token.TokenExpireDuration); err != nil {
			return code.GetError(code.AuthLogoutError)
		}
	}

	if refreshToken != "" {
		if err := token.AddRefreshTokenToBlacklist(ctx.Request.Context(), refreshToken); err != nil {
			glog.Errorf(ctx, "[svcuser.Logout] AddRefreshTokenToBlacklist fail, err:%v", err)
		}
	}

	return nil
}

func (svc *authSvc) RefreshToken(ctx *gin.Context, req *dtouser.RefreshTokenReq) (*dtouser.RefreshTokenResp, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	claims, err := jwtAuth.Parse(req.RefreshToken)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.RefreshToken] Parse refreshToken fail, err:%v", err)
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
		glog.Errorf(ctx, "[svcuser.RefreshToken] GetByID user fail, err:%v, userID:%d", err, userID)
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	if userEntity.Status != model.UserStatusEnabled {
		return nil, code.GetError(code.AuthAccountDisabledError)
	}

	tokenStr, refreshToken, err := svc.generateTokenPair(ctx, *userEntity, userEntity.PersonID)
	if err != nil {
		return nil, err
	}

	if err := addRefreshTokenToBlacklist(ctx, req.RefreshToken); err != nil {
		glog.Errorf(ctx, "[svcuser.RefreshToken] Add refreshToken to blacklist fail, err:%v", err)
	}

	return &dtouser.RefreshTokenResp{
		Token:        tokenStr,
		RefreshToken: refreshToken,
	}, nil
}

func (svc *authSvc) Register(ctx *gin.Context, req *dtouser.RegisterReq) (*dtouser.RegisterResp, error) {
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
		glog.Errorf(ctx, "[svcuser.Register] GeneratePasswordHash fail, err:%v", err)
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
		glog.Errorf(ctx, "[svcuser.Register] Transaction fail, err:%v", txErr)
		return nil, code.GetError(code.AuthRegisterError)
	}

	return &dtouser.RegisterResp{
		UserID:       userID,
		PersonID:     personID,
		Status:       string(userStatus),
		PersonExists: personExists,
		Message:      message,
	}, nil
}

func (svc *authSvc) UnlockAccount(ctx *gin.Context, req *dtouser.UnlockAccountReq) error {
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
		"status":           model.UserStatusEnabled,
		"login_fail_count": 0,
		"locked_until":     nil,
	}); err != nil {
		glog.Errorf(ctx, "[svcuser.UnlockAccount] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
		return code.GetError(code.UserUpdateError)
	}

	return nil
}

func (svc *authSvc) findPersonByAccount(ctx *gin.Context, account string) (*model.PersonEntity, error) {
	personEntity, err := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{
		Mobile: account,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.findPersonByAccount] GetByCond by mobile fail, err:%v, account:%s", err, account)
		return nil, code.GetError(code.AuthLoginError)
	}
	if personEntity != nil && personEntity.ID > 0 {
		return personEntity, nil
	}

	personEntity, err = dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{
		Email: account,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.findPersonByAccount] GetByCond by email fail, err:%v, account:%s", err, account)
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

	tokenStr, err := jwtAuth.Issue(
		fmt.Sprintf("%d", personID),
		tokenIssuer,
		time.Now().Add(tempTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return tokenStr, nil
}

func (svc *authSvc) generateTokenPair(ctx *gin.Context, userEntity model.UserEntity, personID uint) (string, string, error) {
	var orgID uint
	if userEntity.TenantID > 0 {
		tenantEntity, _ := dao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if tenantEntity != nil {
			orgID = tenantEntity.OrgID
		}
	}

	deptID, roleIDs := svc.getUserDeptAndRoles(ctx, userEntity.ID, userEntity.TenantID)

	tokenStr, err := svc.generateTokenWithOrgID(ctx, userEntity, personID, orgID, deptID, roleIDs)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := svc.generateRefreshToken(ctx, userEntity.ID, personID, userEntity.UserType, userEntity.TenantID, orgID)
	if err != nil {
		return "", "", err
	}

	return tokenStr, refreshToken, nil
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
		glog.Errorf(ctx, "[svcuser.isRefreshTokenBlacklisted] Redis Exists fail, err:%v", err)
		return false
	}
	return exists > 0
}

func addRefreshTokenToBlacklist(ctx *gin.Context, tokenStr string) error {
	return token.AddRefreshTokenToBlacklist(ctx.Request.Context(), tokenStr)
}

func (svc *authSvc) buildTenantList(ctx *gin.Context, userList model.UserEntityList) ([]dtouser.TenantListItem, error) {
	tenants := make([]dtouser.TenantListItem, 0, len(userList))
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
		tenants = append(tenants, dtouser.TenantListItem{
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

	tenant, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
		Domain: domain,
		Status: model.TenantStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getCurrentOrg] daoTenant GetByCond fail, err:%v, domain:%s", err, domain)
		return nil, code.GetError(code.AuthLoginError)
	}
	if tenant == nil || tenant.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, tenant.OrgID)
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getCurrentOrg] daoOrg GetByID fail, err:%v, orgID:%d", err, tenant.OrgID)
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
			glog.Errorf(ctx, "[svcuser.filterUsersByOrg] daoTenant GetByID fail, err:%v, tenantID:%d", err, userEntity.TenantID)
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
		glog.Errorf(ctx, "[svcuser.updateLoginInfo] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
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
		glog.Errorf(ctx, "[svcuser.getUserDeptAndRoles] GetListByCond user department fail, err:%v, userID:%d", err, userID)
	} else if len(userDeptList) > 0 {
		deptID = userDeptList[0].DeptID
	}

	userRoleList, err := dao.NewUserRoleDao().GetListByCond(ctx, &dao.UserRoleCond{
		UserID:   userID,
		TenantID: tenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getUserDeptAndRoles] GetListByCond user role fail, err:%v, userID:%d", err, userID)
	} else {
		for _, ur := range userRoleList {
			roleIDs = append(roleIDs, ur.RoleID)
		}
	}

	return deptID, roleIDs
}

func (svc *authSvc) getOrgConfigBool(ctx *gin.Context, orgID uint, configKey string) (bool, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcuser.getOrgConfigBool] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
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
		glog.Errorf(ctx, "[svcuser.getOrgConfigString] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return "", err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return "", nil
	}
	return configEntity.ConfigValue, nil
}

func (svc *authSvc) validateRegisterIdentity(ctx *gin.Context, req *dtouser.RegisterReq, identityType model.RegisterIdentityType) error {
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