package user

import (
	"github.com/morehao/goark/apps/iam/model"
)

type CreatePersonResult struct {
	PersonID uint
	UserID   uint
}

type CreatePersonParams struct {
	Mobile       string
	Email        string
	RealName     string
	OperatorID   uint
	TenantID     uint
	DeptID       uint
	Username     string
	UserType     model.UserType
	Status       model.UserStatus
	EmployeeNo   string
	JobLevel     string
	Position     string
	LastLoginIp  string
	LoginCount   int
	PasswordHash string
}
