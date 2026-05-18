package dtouser

import "github.com/morehao/golib/biz/gobject"

type LoginLogCreateReq struct {
	UserID       uint   `json:"userID" comment:"用户ID"`
	Username     string `json:"username" comment:"用户名"`
	LoginType    string `json:"loginType" comment:"登录类型: password/sms/wechat"`
	LoginStatus  string `json:"loginStatus" comment:"登录状态: success/failed"`
	LoginMessage string `json:"loginMessage" comment:"登录消息"`
	IPAddress    string `json:"ipAddress" comment:"IP地址"`
	Location     string `json:"location" comment:"登录地点"`
	Browser      string `json:"browser" comment:"浏览器"`
	OS           string `json:"os" comment:"操作系统"`
}

type LoginLogCreateResp struct {
	ID uint `json:"id"`
}

type LoginLogPageListReq struct {
	gobject.PageQuery
	UserID      uint   `json:"userID" form:"userID"`
	Username    string `json:"username" form:"username"`
	LoginType   string `json:"loginType" form:"loginType"`
	LoginStatus string `json:"loginStatus" form:"loginStatus"`
}

type LoginLogPageListResp struct {
	List  []LoginLogPageListItem `json:"list"`
	Total int64                  `json:"total"`
}

type LoginLogPageListItem struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"userID"`
	Username     string `json:"username"`
	LoginType    string `json:"loginType"`
	LoginStatus  string `json:"loginStatus"`
	LoginMessage string `json:"loginMessage"`
	IPAddress    string `json:"ipAddress"`
	Location     string `json:"location"`
	Browser      string `json:"browser"`
	OS           string `json:"os"`
	CreatedAt    int64  `json:"createdAt"`
}