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
	"github.com/morehao/golib/gauth/jwtauth"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
)

const (
	tokenExpireDuration     = 24 * time.Hour
	tempTokenExpireDuration = 10 * time.Minute
	tokenBlacklistKeyPrefix = "iam:token:blacklist:"
	tokenIssuer             = "iam"
)

// JWTCustomData JWT自定义数据，字段名需与ginmiddleware.JWTAuth解析后设置到gin.Context的key一致
// 注意：在多租户登录的临时token中，UserID字段存储的是personID，TenantID为0，UserType为"temp"
type JWTCustomData struct {
	UserID   uint   `json:"userId"`
	UserType string `json:"userType"`
	TenantID uint   `json:"tenantId"`
	OrgID    uint   `json:"orgId"`
}

type AuthSvc interface {
	Login(ctx *gin.Context, req *dtoauth.LoginReq) (*dtoauth.LoginResp, error)
	SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error)
	Logout(ctx *gin.Context) error
}

type authSvc struct {
}

var _ AuthSvc = (*authSvc)(nil)

func NewAuthSvc() AuthSvc {
	return &authSvc{}
}

// Login 登录
func (svc *authSvc) Login(ctx *gin.Context, req *dtoauth.LoginReq) (*dtoauth.LoginResp, error) {
	account := strings.TrimSpace(req.Account)

	// 查找自然人(通过手机号或邮箱)
	personEntity, err := svc.findPersonByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := gcrypto.ComparePasswordHash(personEntity.PasswordHash, req.Password); err != nil {
		glog.Errorf(ctx, "[svcauth.Login] password mismatch, account:%s", account)
		return nil, code.GetError(code.AuthPasswordError)
	}

	// 查询该自然人关联的所有用户账号
	userList, err := iamdao.NewUserDao().GetListByCond(ctx, &iamdao.UserCond{
		PersonID: personEntity.ID,
		Status:   "active",
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.Login] GetListByCond fail, err:%v, personID:%d", err, personEntity.ID)
		return nil, code.GetError(code.AuthLoginError)
	}
	if len(userList) == 0 {
		return nil, code.GetError(code.AuthNoTenantError)
	}

	// 单租户：直接返回完整token
	if len(userList) == 1 {
		userEntity := userList[0]
		token, err := svc.generateToken(ctx, userEntity, personEntity.ID)
		if err != nil {
			return nil, err
		}
		tenantEntity, _ := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		tenantName := ""
		if tenantEntity != nil {
			tenantName = tenantEntity.TenantName
		}
		svc.updateLoginInfo(ctx, &userEntity)
		return &dtoauth.LoginResp{
			Token:            token,
			NeedSelectTenant: false,
			Tenants:          []dtoauth.TenantItem{},
			UserInfo: &dtoauth.LoginUserInfo{
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

	// 多租户：返回临时token + 租户列表
	tempToken, err := svc.generateTempToken(personEntity.ID)
	if err != nil {
		return nil, err
	}

	tenants, err := svc.buildTenantList(ctx, userList)
	if err != nil {
		return nil, err
	}

	return &dtoauth.LoginResp{
		Token:            tempToken,
		NeedSelectTenant: true,
		Tenants:          tenants,
		UserInfo:         nil,
	}, nil
}

// SelectTenant 选择租户
func (svc *authSvc) SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error) {
	// 从context获取当前用户信息(临时token中的personID通过UserID字段传递)
	currentUserID := gincontext.GetUserID(ctx)
	if currentUserID == 0 {
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	// 尝试通过UserID获取用户(如果是临时token，UserID实际存的是personID)
	// 先尝试作为personID查找
	userEntity, err := iamdao.NewUserDao().GetByCond(ctx, &iamdao.UserCond{
		PersonID: currentUserID,
		TenantID: req.TenantID,
		Status:   "active",
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByCond fail, err:%v, personID:%d, tenantID:%d", err, currentUserID, req.TenantID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	// 如果personID没找到,尝试通过user获取person后再查找
	if userEntity == nil || userEntity.ID == 0 {
		existingUser, err := iamdao.NewUserDao().GetByID(ctx, currentUserID)
		if err != nil || existingUser == nil || existingUser.ID == 0 {
			return nil, code.GetError(code.AuthTenantSelectError)
		}
		userEntity, err = iamdao.NewUserDao().GetByCond(ctx, &iamdao.UserCond{
			PersonID: existingUser.PersonID,
			TenantID: req.TenantID,
			Status:   "active",
		})
		if err != nil || userEntity == nil || userEntity.ID == 0 {
			return nil, code.GetError(code.AuthTenantSelectError)
		}
	}

	personEntity, err := iamdao.NewPersonDao().GetByID(ctx, userEntity.PersonID)
	if err != nil || personEntity == nil || personEntity.ID == 0 {
		glog.Errorf(ctx, "[svcauth.SelectTenant] GetByID person fail, err:%v, personID:%d", err, userEntity.PersonID)
		return nil, code.GetError(code.AuthTenantSelectError)
	}

	token, err := svc.generateToken(ctx, *userEntity, personEntity.ID)
	if err != nil {
		return nil, err
	}

	tenantEntity, _ := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
	tenantName := ""
	if tenantEntity != nil {
		tenantName = tenantEntity.TenantName
	}

	svc.updateLoginInfo(ctx, userEntity)

	return &dtoauth.SelectTenantResp{
		Token: token,
		UserInfo: &dtoauth.LoginUserInfo{
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
func (svc *authSvc) Logout(ctx *gin.Context) error {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return nil
	}
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil
	}

	key := tokenBlacklistKeyPrefix + hashToken(token)
	if err := dbclient.RedisCli.Set(ctx.Request.Context(), key, "1", tokenExpireDuration).Err(); err != nil {
		glog.Errorf(ctx, "[svcauth.Logout] Redis Set fail, err:%v", err)
		return code.GetError(code.AuthLogoutError)
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
	jwtAuth, err := jwtauth.New[JWTCustomData](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	// 获取用户对应的orgID
	var orgID uint
	if userEntity.TenantID > 0 {
		tenantEntity, _ := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
		if tenantEntity != nil {
			orgID = tenantEntity.OrganizationID
		}
	}

	customData := JWTCustomData{
		UserID:   userEntity.ID,
		UserType: userEntity.UserType,
		TenantID: userEntity.TenantID,
		OrgID:    orgID,
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

func (svc *authSvc) generateTempToken(personID uint) (string, error) {
	jwtAuth, err := jwtauth.New[JWTCustomData](config.Conf.JWT.SignKey)
	if err != nil {
		return "", code.GetError(code.AuthTokenGenerateError)
	}

	customData := JWTCustomData{
		UserID:   personID, // 临时token中使用personID作为UserID
		UserType: "temp",
		TenantID: 0,
		OrgID:    0,
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

func (svc *authSvc) buildTenantList(ctx *gin.Context, userList iammodel.UserEntityList) ([]dtoauth.TenantItem, error) {
	tenants := make([]dtoauth.TenantItem, 0, len(userList))
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
		tenants = append(tenants, dtoauth.TenantItem{
			TenantID:   tenantEntity.ID,
			TenantName: tenantEntity.TenantName,
			OrgID:      tenantEntity.OrganizationID,
			OrgName:    orgName,
		})
	}
	return tenants, nil
}

func (svc *authSvc) updateLoginInfo(ctx *gin.Context, userEntity *iammodel.UserEntity) {
	now := time.Now()
	updateMap := map[string]any{
		"last_login_at": now,
		"last_login_ip": gincontext.GetClientIp(ctx),
		"login_count":   userEntity.LoginCount + 1,
	}
	if err := iamdao.NewUserDao().UpdateMap(ctx, userEntity.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcauth.updateLoginInfo] UpdateMap fail, err:%v, userID:%d", err, userEntity.ID)
	}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
