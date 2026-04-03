package svcauth

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/internal/middleware"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/cryptoutil"
	"github.com/morehao/goark/pkg/ginext"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/gauth/jwtauth"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type AuthSvc interface {
	Login(ctx *gin.Context, req *dtoauth.LoginReq) (*dtoauth.LoginResp, error)
	SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error)
	SwitchTenant(ctx *gin.Context, req *dtoauth.SwitchTenantReq) (*dtoauth.SwitchTenantResp, error)
	Logout(ctx *gin.Context) error
	GetCurrentUser(ctx *gin.Context) (*dtoauth.CurrentUserResp, error)
}

const defaultTokenExpireHours = 24

type authSvc struct{}

var _ AuthSvc = (*authSvc)(nil)

func NewAuthSvc() AuthSvc {
	return &authSvc{}
}

func (svc *authSvc) Login(ctx *gin.Context, req *dtoauth.LoginReq) (*dtoauth.LoginResp, error) {
	personDao := iamdao.NewPersonDao()

	// Try finding person by mobile first, then email
	var person *iammodel.PersonEntity
	var err error

	person, err = personDao.GetByCond(ctx, &iamdao.PersonCond{
		BaseCond: &genericdao.BaseCond{},
		Mobile:   req.Username,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.Login] personDao GetByCond by mobile fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginFailedError)
	}

	if person == nil || person.ID == 0 {
		person, err = personDao.GetByCond(ctx, &iamdao.PersonCond{
			BaseCond: &genericdao.BaseCond{},
			Email:    req.Username,
		})
		if err != nil {
			glog.Errorf(ctx, "[svcauth.Login] personDao GetByCond by email fail, err:%v, req:%s", err, gutil.ToJsonString(req))
			return nil, code.GetError(code.LoginFailedError)
		}
	}

	if person == nil || person.ID == 0 {
		return nil, code.GetError(code.PersonNotFoundError)
	}

	// Verify password
	if err := cryptoutil.ComparePasswordHash(person.PasswordHash, req.Password); err != nil {
		return nil, code.GetError(code.PasswordIncorrectError)
	}

	// Find all active user accounts for this person
	userDao := iamdao.NewUserDao()
	users, err := userDao.GetListByCond(ctx, &iamdao.UserCond{
		BaseCond: &genericdao.BaseCond{},
		PersonID: person.ID,
		Status:   "active",
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.Login] userDao GetListByCond fail, err:%v, personID:%d", err, person.ID)
		return nil, code.GetError(code.LoginFailedError)
	}

	if len(users) == 0 {
		return nil, code.GetError(code.PersonNotFoundError)
	}

	if len(users) == 1 {
		user := users[0]

		if user.Status == "disabled" {
			return nil, code.GetError(code.UserDisabledError)
		}
		if user.Status == "locked" {
			return nil, code.GetError(code.UserLockedError)
		}

		token, err := generateToken(config.Conf.JWT.SignKey, config.Conf.JWT.ExpireHour, user.ID, person.ID, user.TenantID, 0, user.UserType)
		if err != nil {
			glog.Errorf(ctx, "[svcauth.Login] generateToken fail, err:%v, userID:%d", err, user.ID)
			return nil, code.GetError(code.LoginFailedError)
		}

		// Update login info
		now := time.Now()
		_ = userDao.UpdateMap(ctx, user.ID, map[string]any{
			"last_login_at": now,
			"login_count":   user.LoginCount + 1,
		})

		return &dtoauth.LoginResp{
			NeedSelectTenant: false,
			Token:            token,
			PersonID:         person.ID,
			UserInfo: &dtoauth.LoginUserInfo{
				UserID:   user.ID,
				PersonID: person.ID,
				TenantID: user.TenantID,
				Username: user.Username,
				RealName: person.RealName,
				UserType: user.UserType,
			},
		}, nil
	}

	// Multiple tenants - return tenant list
	tenantIDs := make([]uint, 0, len(users))
	for _, u := range users {
		tenantIDs = append(tenantIDs, u.TenantID)
	}

	tenantDao := iamdao.NewTenantDao()
	tenants, err := tenantDao.GetListByCond(ctx, &iamdao.TenantCond{
		BaseCond: &genericdao.BaseCond{
			IDs: tenantIDs,
		},
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.Login] tenantDao GetListByCond fail, err:%v, tenantIDs:%v", err, tenantIDs)
		return nil, code.GetError(code.LoginFailedError)
	}

	tenantMap := tenants.ToMap()
	tenantItems := make([]dtoauth.TenantItem, 0, len(users))
	for _, u := range users {
		t, ok := tenantMap[u.TenantID]
		if !ok {
			continue
		}
		tenantItems = append(tenantItems, dtoauth.TenantItem{
			TenantID:   t.ID,
			TenantName: t.TenantName,
			TenantCode: t.TenantCode,
		})
	}

	return &dtoauth.LoginResp{
		NeedSelectTenant: true,
		PersonID:         person.ID,
		Tenants:          tenantItems,
	}, nil
}

func (svc *authSvc) SelectTenant(ctx *gin.Context, req *dtoauth.SelectTenantReq) (*dtoauth.SelectTenantResp, error) {
	userDao := iamdao.NewUserDao()
	user, err := userDao.GetByCond(ctx, &iamdao.UserCond{
		BaseCond: &genericdao.BaseCond{},
		PersonID: req.PersonID,
		TenantID: req.TenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] userDao GetByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginFailedError)
	}
	if user == nil || user.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}

	if user.Status == "disabled" {
		return nil, code.GetError(code.UserDisabledError)
	}
	if user.Status == "locked" {
		return nil, code.GetError(code.UserLockedError)
	}

	personDao := iamdao.NewPersonDao()
	person, err := personDao.GetByID(ctx, req.PersonID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] personDao GetByID fail, err:%v, personID:%d", err, req.PersonID)
		return nil, code.GetError(code.LoginFailedError)
	}

	token, err := generateToken(config.Conf.JWT.SignKey, config.Conf.JWT.ExpireHour, user.ID, user.PersonID, user.TenantID, 0, user.UserType)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SelectTenant] generateToken fail, err:%v, userID:%d", err, user.ID)
		return nil, code.GetError(code.LoginFailedError)
	}

	now := time.Now()
	_ = userDao.UpdateMap(ctx, user.ID, map[string]any{
		"last_login_at": now,
		"login_count":   user.LoginCount + 1,
	})

	return &dtoauth.SelectTenantResp{
		Token: token,
		UserInfo: &dtoauth.LoginUserInfo{
			UserID:   user.ID,
			PersonID: user.PersonID,
			TenantID: user.TenantID,
			Username: user.Username,
			RealName: person.RealName,
			UserType: user.UserType,
		},
	}, nil
}

func (svc *authSvc) SwitchTenant(ctx *gin.Context, req *dtoauth.SwitchTenantReq) (*dtoauth.SwitchTenantResp, error) {
	personID := ginext.GetPersonID(ctx)
	if personID == 0 {
		return nil, code.GetError(code.LoginFailedError)
	}

	userDao := iamdao.NewUserDao()
	user, err := userDao.GetByCond(ctx, &iamdao.UserCond{
		BaseCond: &genericdao.BaseCond{},
		PersonID: personID,
		TenantID: req.TenantID,
	})
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SwitchTenant] userDao GetByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginFailedError)
	}
	if user == nil || user.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}

	if user.Status == "disabled" {
		return nil, code.GetError(code.UserDisabledError)
	}
	if user.Status == "locked" {
		return nil, code.GetError(code.UserLockedError)
	}

	personDao := iamdao.NewPersonDao()
	person, err := personDao.GetByID(ctx, personID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SwitchTenant] personDao GetByID fail, err:%v, personID:%d", err, personID)
		return nil, code.GetError(code.LoginFailedError)
	}

	token, err := generateToken(config.Conf.JWT.SignKey, config.Conf.JWT.ExpireHour, user.ID, user.PersonID, user.TenantID, 0, user.UserType)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.SwitchTenant] generateToken fail, err:%v, userID:%d", err, user.ID)
		return nil, code.GetError(code.LoginFailedError)
	}

	// Blacklist old token
	oldToken := extractToken(ctx)
	if oldToken != "" {
		_ = middleware.AddToBlacklist(ctx.Request.Context(), oldToken, tokenExpiry())
	}

	now := time.Now()
	_ = userDao.UpdateMap(ctx, user.ID, map[string]any{
		"last_login_at": now,
	})

	return &dtoauth.SwitchTenantResp{
		Token: token,
		UserInfo: &dtoauth.LoginUserInfo{
			UserID:   user.ID,
			PersonID: user.PersonID,
			TenantID: user.TenantID,
			Username: user.Username,
			RealName: person.RealName,
			UserType: user.UserType,
		},
	}, nil
}

func (svc *authSvc) Logout(ctx *gin.Context) error {
	tokenStr := extractToken(ctx)
	if tokenStr == "" {
		return nil
	}

	if err := middleware.AddToBlacklist(ctx.Request.Context(), tokenStr, tokenExpiry()); err != nil {
		glog.Errorf(ctx, "[svcauth.Logout] AddToBlacklist fail, err:%v", err)
		return code.GetError(code.LogoutFailedError)
	}
	return nil
}

func (svc *authSvc) GetCurrentUser(ctx *gin.Context) (*dtoauth.CurrentUserResp, error) {
	userID := gincontext.GetUserID(ctx)
	if userID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}

	userDao := iamdao.NewUserDao()
	user, err := userDao.GetByID(ctx, userID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.GetCurrentUser] userDao GetByID fail, err:%v, userID:%d", err, userID)
		return nil, code.GetError(code.UserGetDetailError)
	}
	if user == nil || user.ID == 0 {
		return nil, code.GetError(code.UserNotExistError)
	}

	personDao := iamdao.NewPersonDao()
	person, err := personDao.GetByID(ctx, user.PersonID)
	if err != nil {
		glog.Errorf(ctx, "[svcauth.GetCurrentUser] personDao GetByID fail, err:%v, personID:%d", err, user.PersonID)
		return nil, code.GetError(code.UserGetDetailError)
	}

	return &dtoauth.CurrentUserResp{
		UserID:   user.ID,
		PersonID: user.PersonID,
		TenantID: user.TenantID,
		Username: user.Username,
		RealName: person.RealName,
		UserType: user.UserType,
	}, nil
}

func generateToken(signKey string, expireHour int, userID, personID, tenantID, orgID uint, userType string) (string, error) {
	if expireHour <= 0 {
		expireHour = defaultTokenExpireHours
	}
	expiresAt := time.Now().Add(time.Duration(expireHour) * time.Hour)
	customData := middleware.JWTCustomData{
		UserID:   userID,
		PersonID: personID,
		TenantID: tenantID,
		OrgID:    orgID,
		UserType: userType,
	}
	claims := jwtauth.NewClaims(strconv.FormatUint(uint64(userID), 10), expiresAt, customData)
	return jwtauth.CreateToken(signKey, claims)
}

func extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func tokenExpiry() time.Time {
	expireHour := config.Conf.JWT.ExpireHour
	if expireHour <= 0 {
		expireHour = defaultTokenExpireHours
	}
	return time.Now().Add(time.Duration(expireHour) * time.Hour)
}
