package svcuser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type PasswordSvc struct {
}

func NewPasswordSvc() *PasswordSvc {
	return &PasswordSvc{}
}

type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumber    bool
	RequireSpecial   bool
	MaxFailCount     int
	LockDuration     int
}

func (s *PasswordSvc) GetPasswordPolicy(ctx *gin.Context, orgID uint) (*PasswordPolicy, error) {
	policy := &PasswordPolicy{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   false,
		MaxFailCount:     5,
		LockDuration:     300,
	}

	policy.MinLength = s.getOrgConfigInt(ctx, orgID, model.OrgConfigKeyPasswordMinLength, 8)
	policy.RequireUppercase = s.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyPasswordRequireUppercase, true)
	policy.RequireLowercase = s.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyPasswordRequireLowercase, true)
	policy.RequireNumber = s.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyPasswordRequireNumber, true)
	policy.RequireSpecial = s.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyPasswordRequireSpecial, false)
	policy.MaxFailCount = s.getOrgConfigInt(ctx, orgID, model.OrgConfigKeyLoginMaxFailCount, 5)
	policy.LockDuration = s.getOrgConfigInt(ctx, orgID, model.OrgConfigKeyLoginLockDuration, 300)

	return policy, nil
}

func (s *PasswordSvc) ValidatePasswordComplexity(ctx *gin.Context, orgID uint, password string) error {
	policy, err := s.GetPasswordPolicy(ctx, orgID)
	if err != nil {
		return err
	}

	if len(password) < policy.MinLength {
		return code.GetError(code.PasswordComplexityError)
	}

	if policy.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return code.GetError(code.PasswordComplexityError)
	}

	if policy.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return code.GetError(code.PasswordComplexityError)
	}

	if policy.RequireNumber && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return code.GetError(code.PasswordComplexityError)
	}

	if policy.RequireSpecial && !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return code.GetError(code.PasswordComplexityError)
	}

	return nil
}

func (s *PasswordSvc) CheckUserLockStatus(ctx *gin.Context, userID uint) error {
	userEntity, err := dao.NewUserDao().GetByID(ctx, userID)
	if err != nil || userEntity == nil || userEntity.ID == 0 {
		return code.GetError(code.UserNotExistError)
	}

	if userEntity.Status == model.UserStatusLocked {
		if userEntity.LockedUntil != nil && time.Now().After(*userEntity.LockedUntil) {
			dao.NewUserDao().UpdateMap(ctx, userID, map[string]any{
				"status":           model.UserStatusEnabled,
				"login_fail_count": 0,
				"locked_until":     nil,
			})
		} else {
			return code.GetError(code.AuthUserLockedError)
		}
	}

	if userEntity.Status == model.UserStatusDisabled {
		return code.GetError(code.AuthAccountDisabledError)
	}

	return nil
}

func (s *PasswordSvc) RecordLoginFail(ctx *gin.Context, userID uint, orgID uint) error {
	policy, err := s.GetPasswordPolicy(ctx, orgID)
	if err != nil {
		return err
	}

	userEntity, err := dao.NewUserDao().GetByID(ctx, userID)
	if err != nil || userEntity == nil {
		return err
	}

	newFailCount := userEntity.LoginFailCount + 1
	updateMap := map[string]any{
		"login_fail_count": newFailCount,
	}

	if newFailCount >= policy.MaxFailCount {
		updateMap["status"] = model.UserStatusLocked
		updateMap["locked_until"] = time.Now().Add(time.Duration(policy.LockDuration) * time.Second)
	}

	return dao.NewUserDao().UpdateMap(ctx, userID, updateMap)
}

func (s *PasswordSvc) ClearLoginFail(ctx *gin.Context, userID uint) error {
	return dao.NewUserDao().UpdateMap(ctx, userID, map[string]any{
		"login_fail_count": 0,
	})
}

func (s *PasswordSvc) getOrgConfigInt(ctx *gin.Context, orgID uint, key string, defaultVal int) int {
	if orgID == 0 {
		orgID = gincontext.GetOrgID(ctx)
	}
	if orgID == 0 {
		return defaultVal
	}
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: key,
	})
	if err != nil || configEntity == nil || configEntity.ID == 0 {
		return defaultVal
	}
	var val int
	if _, err := fmt.Sscanf(configEntity.ConfigValue, "%d", &val); err != nil {
		return defaultVal
	}
	return val
}

func (s *PasswordSvc) getOrgConfigBool(ctx *gin.Context, orgID uint, key string, defaultVal bool) bool {
	if orgID == 0 {
		orgID = gincontext.GetOrgID(ctx)
	}
	if orgID == 0 {
		return defaultVal
	}
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: key,
	})
	if err != nil || configEntity == nil || configEntity.ID == 0 {
		return defaultVal
	}
	return strings.ToLower(configEntity.ConfigValue) == "true"
}