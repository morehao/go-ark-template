package objuser

type UserBaseInfo struct {
	// CompanyID 公司id
	CompanyID uint `json:"companyID" form:"companyID"`
	// DepartmentID 部门id
	DepartmentID uint `json:"departmentID" form:"departmentID"`
	// Name 姓名
	Name string `json:"name" form:"name"`
}
