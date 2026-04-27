package strategy

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
)

type strategyCommon struct {
}

func newStrategyCommon() *strategyCommon {
	return &strategyCommon{}
}

func (sc *strategyCommon) getOrgConfigBool(ctx *gin.Context, orgID uint, configKey string) (bool, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[strategyCommon.getOrgConfigBool] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return false, err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return false, nil
	}
	return strings.ToLower(configEntity.ConfigValue) == "true", nil
}

func (sc *strategyCommon) getOrgConfigString(ctx *gin.Context, orgID uint, configKey string) (string, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[strategyCommon.getOrgConfigString] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return "", err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return "", nil
	}
	return configEntity.ConfigValue, nil
}

func (sc *strategyCommon) getOrgConfigInt(ctx *gin.Context, orgID uint, configKey string) (int, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		glog.Errorf(ctx, "[strategyCommon.getOrgConfigInt] GetByCond fail, err:%v, orgID:%d, key:%s", err, orgID, configKey)
		return 0, err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return 0, nil
	}
	var val int
	_, err = fmt.Sscanf(configEntity.ConfigValue, "%d", &val)
	return val, err
}

func (sc *strategyCommon) createRegisterResult(ctx *gin.Context, orgID, tenantID uint, req *RegisterRequest) (*RegisterResult, error) {
	registerEnabled, err := sc.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyRegisterEnabled)
	if err != nil {
		glog.Errorf(ctx, "[createRegisterResult] GetBool registerEnabled fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterDisabled)
	}
	if !registerEnabled {
		return nil, code.GetError(code.AuthRegisterDisabled)
	}

	identityType, err := sc.getOrgConfigString(ctx, orgID, model.OrgConfigKeyRegisterIdentityType)
	if err != nil {
		glog.Errorf(ctx, "[createRegisterResult] GetString identityType fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}
	if identityType == "" {
		identityType = string(model.RegisterIdentityTypeEmail)
	}
	if err := sc.validateIdentity(ctx, req, model.RegisterIdentityType(identityType)); err != nil {
		return nil, err
	}

	passwordHash, err := gcrypto.GeneratePasswordHash(req.Password)
	if err != nil {
		glog.Errorf(ctx, "[createRegisterResult] GeneratePasswordHash fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}

	requireApproval, err := sc.getOrgConfigBool(ctx, orgID, model.OrgConfigKeyRegisterRequireApproval)
	if err != nil {
		glog.Errorf(ctx, "[createRegisterResult] GetBool requireApproval fail, err:%v", err)
		return nil, code.GetError(code.AuthRegisterError)
	}
	userStatus := model.UserStatusEnabled
	message := "注册成功"
	if requireApproval {
		userStatus = model.UserStatusPending
		message = "注册成功，等待管理员审核"
	}

	return &RegisterResult{
		TenantID:     tenantID,
		PasswordHash: passwordHash,
		Status:       userStatus,
		Message:      message,
	}, nil
}

func (sc *strategyCommon) validateIdentity(ctx *gin.Context, req *RegisterRequest, identityType model.RegisterIdentityType) error {
	switch identityType {
	case model.RegisterIdentityTypeEmail:
		email := strings.TrimSpace(req.Email)
		if email == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
		person, err := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{Email: email})
		if err != nil {
			glog.Errorf(ctx, "[validateIdentity] GetByCond email fail, err:%v", err)
			return code.GetError(code.AuthRegisterError)
		}
		if person != nil && person.ID > 0 {
			return code.GetError(code.AuthRegisterError)
		}
	case model.RegisterIdentityTypeMobile:
		mobile := strings.TrimSpace(req.Mobile)
		if mobile == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
		person, err := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{Mobile: mobile})
		if err != nil {
			glog.Errorf(ctx, "[validateIdentity] GetByCond mobile fail, err:%v", err)
			return code.GetError(code.AuthRegisterError)
		}
		if person != nil && person.ID > 0 {
			return code.GetError(code.AuthRegisterError)
		}
	case model.RegisterIdentityTypeBoth:
		email := strings.TrimSpace(req.Email)
		mobile := strings.TrimSpace(req.Mobile)
		if email == "" || mobile == "" {
			return code.GetError(code.AuthRegisterIdentityRequired)
		}
		person, err := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{Email: email, Mobile: mobile})
		if err != nil {
			glog.Errorf(ctx, "[validateIdentity] GetByCond both fail, err:%v", err)
			return code.GetError(code.AuthRegisterError)
		}
		if person != nil && person.ID > 0 {
			return code.GetError(code.AuthRegisterError)
		}
	}
	return nil
}

func (sc *strategyCommon) generateEmployeeNo(ctx *gin.Context, tenantCode string) (string, error) {
	if len(tenantCode) < 2 {
		tenantCode = fmt.Sprintf("%-2s", tenantCode)
	}
	today := time.Now().Format("20060102")
	key := fmt.Sprintf("employee_no:%s:%s", tenantCode, today)

	seq, err := dbclient.RedisCli.Incr(ctx, key).Result()
	if err != nil {
		glog.Errorf(ctx, "[generateEmployeeNo] Redis Incr fail, err:%v", err)
		return "", code.GetError(code.AuthRegisterError)
	}

	expiry := time.Now().AddDate(0, 0, 1)
	if _, err := dbclient.RedisCli.ExpireAt(ctx, key, expiry).Result(); err != nil {
		glog.Errorf(ctx, "[generateEmployeeNo] Redis ExpireAt fail, err:%v", err)
	}

	return fmt.Sprintf("%s%s%04d", tenantCode[:2], today, seq), nil
}

func (sc *strategyCommon) assignDefaultRolesAndDepts(ctx *gin.Context, orgID, userID uint) error {
	defaultRoleIDsStr, err := sc.getOrgConfigString(ctx, orgID, model.OrgConfigKeyRegisterDefaultRoles)
	if err != nil {
		glog.Errorf(ctx, "[assignDefaultRolesAndDepts] GetConfig roles fail, err:%v", err)
		return code.GetError(code.AuthRegisterError)
	}
	if defaultRoleIDsStr != "" && defaultRoleIDsStr != "[]" {
		var roleIDs []uint
		if err := json.Unmarshal([]byte(defaultRoleIDsStr), &roleIDs); err == nil {
			for _, roleID := range roleIDs {
				userRole := &model.UserRoleEntity{
					UserID: userID,
					RoleID: roleID,
				}
				if err := dao.NewUserRoleDao().Insert(ctx, userRole); err != nil {
					glog.Errorf(ctx, "[assignDefaultRolesAndDepts] Insert role fail, userID:%d, roleID:%d, err:%v", userID, roleID, err)
				}
			}
		}
	}

	defaultDeptIDsStr, err := sc.getOrgConfigString(ctx, orgID, model.OrgConfigKeyRegisterDefaultDepts)
	if err != nil {
		glog.Errorf(ctx, "[assignDefaultRolesAndDepts] GetConfig depts fail, err:%v", err)
		return code.GetError(code.AuthRegisterError)
	}
	if defaultDeptIDsStr != "" && defaultDeptIDsStr != "[]" {
		var deptIDs []uint
		if err := json.Unmarshal([]byte(defaultDeptIDsStr), &deptIDs); err == nil {
			for _, deptID := range deptIDs {
				userDept := &model.UserDepartmentEntity{
					UserID: userID,
					DeptID: deptID,
				}
				if err := dao.NewUserDepartmentDao().Insert(ctx, userDept); err != nil {
					glog.Errorf(ctx, "[assignDefaultRolesAndDepts] Insert dept fail, userID:%d, deptID:%d, err:%v", userID, deptID, err)
				}
			}
		}
	}
	return nil
}