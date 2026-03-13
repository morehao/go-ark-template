package code

import "github.com/morehao/golib/gerror"

const (
	CompanyCreateError      = 100300
	CompanyDeleteError      = 100301
	CompanyUpdateError      = 100302
	CompanyGetDetailError   = 100303
	CompanyGetPageListError = 100304
	CompanyNotExistError    = 100305
)

var companyErrorMsgMap = gerror.CodeMsgMap{
	CompanyCreateError:      "创建公司管理失败",
	CompanyDeleteError:      "删除公司管理失败",
	CompanyUpdateError:      "修改公司管理失败",
	CompanyGetDetailError:   "查看公司管理失败",
	CompanyGetPageListError: "查看公司管理列表失败",
	CompanyNotExistError:    "公司管理不存在",
}

const (
	DepartmentCreateError      = 100400
	DepartmentDeleteError      = 100401
	DepartmentUpdateError      = 100402
	DepartmentGetDetailError   = 100403
	DepartmentGetPageListError = 100404
	DepartmentNotExistError    = 100405
)

var departmentErrorMsgMap = gerror.CodeMsgMap{
	DepartmentCreateError:      "创建部门管理失败",
	DepartmentDeleteError:      "删除部门管理失败",
	DepartmentUpdateError:      "修改部门管理失败",
	DepartmentGetDetailError:   "查看部门管理失败",
	DepartmentGetPageListError: "查看部门管理列表失败",
	DepartmentNotExistError:    "部门管理不存在",
}
