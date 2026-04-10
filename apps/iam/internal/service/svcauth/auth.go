package svcauth

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/gauth/jwtauth"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
)

const (
	tokenExpireDuration            = 24 * time.Hour
	tempTokenExpireDuration        = 10 * time.Minute
	refreshTokenExpireDuration     = 7 * 24 * time.Hour
	tokenBlacklistKeyPrefix        = "iam:token:blacklist:"
	refreshTokenBlacklistKeyPrefix = "iam:refreshToken:blacklist:"
	tokenIssuer                    = "iam"
)

type AuthSvc interface {
	LoginByPassword(ctx *gin.Context, req *dtoauth.LoginByPasswordReq) (*dtoauth.LoginByPasswordResp, error)
	SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error)
	Logout(ctx *gin.Context, refreshToken string) error
	RefreshToken(ctx *gin.Context, req *dtoauth.RefreshTokenReq) (*dtoauth.RefreshTokenResp, error)
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

	organizationEntity, err := svc.getCurrentOrganization(ctx)
	if err != nil {
		return nil, err
	}

	// 查找自然人(通过手机号或邮箱)
	personEntity, err := svc.findPersonByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := gcrypto.ComparePasswordHash(personEntity.PasswordHash, req.Password); err != nil {
		glog.Errorf(ctx, "[svcauth.LoginByPassword] password mismatch, account:%s", account)
		return nil, code.GetError(code.AuthPasswordError)
	}

	// 查询该自然人关联的所有用户账号
	userList, err := iamdao.NewUserDao().GetListByCond(ctx, &iamdao.UserCond{
		PersonID: personEntity.ID,
		Status:   iammodel.UserStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.LoginByPassword] GetListByCond fail, err:%v, personID:%d", err, personEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
	}
	userList, err = svc.filterUsersByOrganization(ctx, userList, organizationEntity.ID)
	if err != nil {
		return nil, err
	}
	if len(userList) == 0 {
		return nil, code.GetError(code.AuthNoTenantError)
	}

	// 统一返回临时token + 租户列表
	tempToken, err := svc.generateTempToken(personEntity.ID, organizationEntity.ID)
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

	tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, req.TenantID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByID tenant fail, err:%v, tenantID:%d", err, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}
	if tenantEntity == nil || tenantEntity.ID == 0 || tenantEntity.OrganizationID != orgID {
		return nil, code.GetError(code.AuthTenantNotInOrgError)
	}

	userEntity, err := iamdao.NewUserDao().GetByCond(ctx, &iamdao.UserCond{
		PersonID: personID,
		TenantID: req.TenantID,
		Status:   iammodel.UserStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByCond fail, err:%v, personID:%d, tenantID:%d", err, personID, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}
	if userEntity == nil || userEntity.ID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	personEntity, err := iamdao.NewPersonDao().GetByID(ctx, userEntity.PersonID)
	if err != nil || personEntity == nil || personEntity.ID == 0 {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByID person fail, err:%v, personID:%d", err, userEntity.PersonID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	token, refreshToken, err := svc.generateTokenPair(ctx, *userEntity, personEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantEntity, _ = iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
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
	// 1. 将 access token 加入黑名单
	token := ctx.GetHeader("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		if token != "" {
			key := tokenBlacklistKeyPrefix + hashToken(token)
			if err := dbclient.RedisCli.Set(ctx.Request.Context(), key, "1", tokenExpireDuration).Err(); err != nil {
				glog.Errorf(ctx, "[svcauth.Logout] Redis Set access token fail, err:%v", err)
				return code.GetError(code.AuthLogoutError)
			}
		}
	}

	// 2. 将 refreshToken 加入黑名单（如果提供了）
	if refreshToken != "" {
		refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
		if refreshToken != "" {
			key := refreshTokenBlacklistKeyPrefix + hashToken(refreshToken)
			if err := dbclient.RedisCli.Set(ctx.Request.Context(), key, "1", refreshTokenExpireDuration).Err(); err != nil {
				glog.Errorf(ctx, "[svcauth.Logout] Redis Set refresh token fail, err:%v", err)
			}
		}
	}

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func IsTokenBlacklisted(ctx *gin.Context, token string) bool {
	if dbclient.RedisCli == nil {
		return false
	}
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return false
	}
	key := tokenBlacklistKeyPrefix + hashToken(token)
	exists, err := dbclient.RedisCli.Exists(ctx.Request.Context(), key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

func (svc *authSvc) findPersonByAccount(ctx *gin.Context, account string) (*iammodel.PersonEntity, error) {
	// 先按手机号查询
	personEntity, err := iamdao.NewPersonDao().GetByCond(ctx, &iamdao.PersonCond{
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
	personEntity, err = iamdao.NewPersonDao().GetByCond(ctx, &iamdao.PersonCond{
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

func (svc *authSvc) generateToken(ctx *gin.Context, userEntity iammodel.UserEntity, personID uint) (string, error) {
	jwtAuth, err := jwtauth.New[gobject.UserClaims](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	var orgID uint
	if userEntity.TenantID > 0 {
		tenantEntity, _ := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if tenantEntity != nil {
			orgID = tenantEntity.OrganizationID
		}
	}

	deptID, roleIDs := svc.getUserDeptAndRoles(ctx, userEntity.ID, userEntity.TenantID)

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

	token, err := jwtAuth.Issue(
		fmt.Sprintf("%d", userEntity.ID),
		tokenIssuer,
		time.Now().Add(tokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return token, nil
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

func (svc *authSvc) generateTokenPair(ctx *gin.Context, userEntity iammodel.UserEntity, personID uint) (token string, refreshToken string, err error) {
	var orgID uint
	if userEntity.TenantID > 0 {
		tenantEntity, _ := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if tenantEntity != nil {
			orgID = tenantEntity.OrganizationID
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

func (svc *authSvc) generateTokenWithOrgID(ctx *gin.Context, userEntity iammodel.UserEntity, personID uint, orgID uint, deptID uint, roleIDs []uint) (string, error) {
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

	token, err := jwtAuth.Issue(
		fmt.Sprintf("%d", userEntity.ID),
		tokenIssuer,
		time.Now().Add(tokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return token, nil
}

func (svc *authSvc) generateRefreshToken(ctx *gin.Context, userID uint, personID uint, userType iammodel.UserType, tenantID uint, orgID uint) (string, error) {
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

	token, err := jwtAuth.Issue(
		fmt.Sprintf("%d", userID),
		tokenIssuer,
		time.Now().Add(refreshTokenExpireDuration),
		customData,
	)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}
	return token, nil
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
	userEntity, err := iamdao.NewUserDao().GetByID(ctx, userID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		glog.Errorf(ctx, "[svcauth.RefreshToken] GetByID user fail, err:%v, userID:%d", err, userID)
		return nil, code.GetError(code.AuthRefreshTokenInvalidError)
	}

	if userEntity.Status != iammodel.UserStatusEnabled {
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

func isRefreshTokenBlacklisted(ctx *gin.Context, token string) bool {
	if dbclient.RedisCli == nil {
		return false
	}
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return false
	}
	key := refreshTokenBlacklistKeyPrefix + hashToken(token)
	exists, err := dbclient.RedisCli.Exists(ctx.Request.Context(), key).Result()
	if err != nil {
		glog.Errorf(ctx, "[svcauth.isRefreshTokenBlacklisted] Redis Exists fail, err:%v", err)
		return false
	}
	return exists > 0
}

func addRefreshTokenToBlacklist(ctx *gin.Context, token string) error {
	if dbclient.RedisCli == nil {
		return nil
	}
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil
	}
	key := refreshTokenBlacklistKeyPrefix + hashToken(token)
	return dbclient.RedisCli.Set(ctx.Request.Context(), key, "1", refreshTokenExpireDuration).Err()
}

func (svc *authSvc) buildTenantList(ctx *gin.Context, userList iammodel.UserEntityList) ([]dtoauth.TenantListItem, error) {
	tenants := make([]dtoauth.TenantListItem, 0, len(userList))
	for _, u := range userList {
		tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, u.TenantID)
		if err != nil || tenantEntity == nil || tenantEntity.ID == 0 {
			continue
		}
		orgName := ""
		orgEntity, _ := iamdao.NewOrganizationDao().GetByID(ctx, tenantEntity.OrganizationID)
		if orgEntity != nil {
			orgName = orgEntity.OrganizationName
		}
		tenants = append(tenants, dtoauth.TenantListItem{
			TenantID:   tenantEntity.ID,
			TenantName: tenantEntity.TenantName,
			OrgID:      tenantEntity.OrganizationID,
			OrgName:    orgName,
		})
	}
	return tenants, nil
}

func (svc *authSvc) getCurrentOrganization(ctx *gin.Context) (*iammodel.OrganizationEntity, error) {
	domain := resolveDomain(ctx)
	if domain == "" {
		return nil, code.GetError(code.AuthOrganizationNotFoundError)
	}

	organizationEntity, err := iamdao.NewOrganizationDao().GetByCond(ctx, &iamdao.OrganizationCond{
		Domain: domain,
		Status: iammodel.OrgStatusEnabled,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getCurrentOrganization] daoOrganization GetByCond fail, err:%v, domain:%s", err, domain)
		return nil, code.GetError(code.AuthLoginError)
	}
	if organizationEntity == nil || organizationEntity.ID == 0 {
		return nil, code.GetError(code.AuthOrganizationNotFoundError)
	}
	return organizationEntity, nil
}

func (svc *authSvc) filterUsersByOrganization(ctx *gin.Context, userList iammodel.UserEntityList, organizationID uint) (iammodel.UserEntityList, error) {
	filtered := make(iammodel.UserEntityList, 0, len(userList))
	for _, userEntity := range userList {
		tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if err != nil {
			glog.Errorf(ctx, "[svcauth.filterUsersByOrganization] daoTenant GetByID fail, err:%v, tenantID:%d", err, userEntity.TenantID)
			return nil, code.GetError(code.AuthLoginError)
		}
		if tenantEntity == nil || tenantEntity.ID == 0 || tenantEntity.Status != iammodel.TenantStatusEnabled {
			continue
		}
		if tenantEntity.OrganizationID != organizationID {
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

func (svc *authSvc) updateLoginInfo(ctx *gin.Context, userEntity *iammodel.UserEntity) {
	now := time.Now()
	updateMap := map[string]any{
		"last_login_at": now,
		"last_login_ip": gincontext.GetClientIP(ctx),
		"login_count":   userEntity.LoginCount + 1,
	}
	if err := iamdao.NewUserDao().UpdateMap(ctx, userEntity.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcauth.updateLoginInfo] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
	}
}

func (svc *authSvc) getUserDeptAndRoles(ctx *gin.Context, userID uint, tenantID uint) (deptID uint, roleIDs []uint) {
	deptID = 0
	roleIDs = []uint{}

	userDeptList, err := iamdao.NewUserDepartmentDao().GetListByCond(ctx, &iamdao.UserDepartmentCond{
		UserID:   userID,
		TenantID: tenantID,
		DeptType: iammodel.UserDeptTypePrimary,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.getUserDeptAndRoles] GetListByCond user department fail, err:%v, userID:%d", err, userID)
	} else if len(userDeptList) > 0 {
		deptID = userDeptList[0].DeptID
	}

	userRoleList, err := iamdao.NewUserRoleDao().GetListByCond(ctx, &iamdao.UserRoleCond{
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

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
