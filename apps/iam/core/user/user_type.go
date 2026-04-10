package user

import (
	"github.com/morehao/goark/apps/iam/iammodel"
)

type CreatePersonResult struct {
	PersonID uint
	UserID   uint
}

type CreatePersonParams struct {
	Mobile      string
	Email       string
	RealName    string
	OperatorID  uint
	TenantID    uint
	DeptID      uint
	Username    string
	UserType    iammodel.UserType
	Status      iammodel.UserStatus
	EmployeeNo  string
	JobLevel    string
	Position    string
	LastLoginIp string
	LoginCount  int32
}
