package code

import (
	"fmt"

	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/gerror"
)

var errorMap = gerror.ErrorMap{}

func registerError(codeMsgMap gerror.CodeMsgMap) {
	for code, msg := range codeMsgMap {

		if _, ok := errorMap[code]; ok {
			panic(fmt.Sprintf("error code %d already exists", code))
		}
		errorMap[code] = gerror.Error{
			Code: code,
			Msg:  msg,
		}
	}
}

func GetError(code int) *gerror.Error {
	err := errorMap[code]
	return &err
}

func init() {
	// 业务错误码从110100开始
	registerError(genericdao.DBErrorMsgMap)
	registerError(gconstant.SystemErrorMsgMap)
	registerError(gconstant.AuthErrorMsgMap)
	registerError(tenantErrorMsgMap)
	registerError(companyErrorMsgMap)
	registerError(departmentErrorMsgMap)
	registerError(userErrorMsgMap)
	registerError(menuErrorMsgMap)
	registerError(roleErrorMsgMap)
}
