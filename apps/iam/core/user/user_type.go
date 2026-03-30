package user

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
	UserType    string
	Status      string
	EmployeeNo  string
	JobLevel    string
	Position    string
	LastLoginIp string
	LoginCount  int32
}
