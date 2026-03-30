package user

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
	"gorm.io/gorm"
)

func CreatePersonWithUser(ctx *gin.Context, tx *gorm.DB, params *CreatePersonParams) (*CreatePersonResult, error) {
	mobile := params.Mobile
	email := params.Email

	passwordHash, err := GeneratePassword(mobile, email)
	if err != nil {
		glog.Errorf(ctx, "[user.CreatePersonWithUser] GeneratePassword fail, err:%v, mobile:%s, email:%s", err, mobile, email)
		return nil, code.GetError(code.PersonCreateError)
	}

	personCreateEntity := &iammodel.PersonEntity{
		Mobile:       mobile,
		Email:        email,
		RealName:     params.RealName,
		PasswordHash: passwordHash,
		CreatedBy:    params.OperatorID,
		UpdatedBy:    params.OperatorID,
	}

	personEntity, err := iamdao.NewPersonDao().GetByCond(ctx, &iamdao.PersonCond{
		Mobile: mobile,
		Email:  email,
	})
	if err != nil {
		glog.Errorf(ctx, "[user.CreatePersonWithUser] GetByCond fail, err:%v, mobile:%s, email:%s", err, mobile, email)
		return nil, code.GetError(code.PersonCreateError)
	}

	var personID uint
	if personEntity != nil && personEntity.ID != 0 {
		personID = personEntity.ID
	} else {
		if err := iamdao.NewPersonDao().WithTx(tx).Insert(ctx, personCreateEntity); err != nil {
			glog.Errorf(ctx, "[user.CreatePersonWithUser] Insert person fail, err:%v, mobile:%s, email:%s", err, mobile, email)
			return nil, code.GetError(code.PersonCreateError)
		}
		personID = personCreateEntity.ID
	}

	userEntity := &iammodel.UserEntity{
		TenantID:    params.TenantID,
		DeptID:      params.DeptID,
		PersonID:    personID,
		Username:    params.Username,
		UserType:    params.UserType,
		Status:      params.Status,
		EmployeeNo:  params.EmployeeNo,
		JobLevel:    params.JobLevel,
		Position:    params.Position,
		LastLoginIp: params.LastLoginIp,
		LoginCount:  params.LoginCount,
		CreatedBy:   params.OperatorID,
		UpdatedBy:   params.OperatorID,
	}
	if err := iamdao.NewUserDao().WithTx(tx).Insert(ctx, userEntity); err != nil {
		glog.Errorf(ctx, "[user.CreatePersonWithUser] Insert user fail, err:%v, mobile:%s, email:%s", err, mobile, email)
		return nil, code.GetError(code.PersonCreateError)
	}

	return &CreatePersonResult{
		PersonID: personID,
		UserID:   userEntity.ID,
	}, nil
}

// GeneratePassword 根据用户身份标识生成密码哈希
// identities: 可变参数，优先使用第一个非空的标识（通常是手机号或邮箱）
// prefix: 从配置读取的密码前缀，默认为 "pwd"
// 最终生成: hash(prefix + 第一个非空identity)
func GeneratePassword(identities ...string) (string, error) {
	prefix := "pwd"
	if config.Conf != nil && config.Conf.Password.Prefix != "" {
		prefix = config.Conf.Password.Prefix
	}
	var suffix string
	for _, identity := range identities {
		if identity != "" {
			suffix = identity
			break
		}
	}
	if suffix == "" {
		return "", code.GetError(code.PersonCreateError)
	}
	hash, err := gcrypto.GeneratePasswordHash(prefix + suffix)
	if err != nil {
		return "", code.GetError(code.PersonCreateError)
	}
	return hash, nil
}
