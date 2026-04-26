package dtolog

type OperationLogCreateReq struct {
	UserID        uint   `json:"userID" comment:"操作人ID"`
	Username      string `json:"username" comment:"操作人账号"`
	Module        string `json:"module" comment:"操作模块"`
	Operation     string `json:"operation" comment:"操作类型: create/update/delete/query"`
	Method        string `json:"method" comment:"请求方法: GET/POST/PUT/DELETE等"`
	RequestID     string `json:"requestID" comment:"请求ID"`
	RequestURL    string `json:"requestURL" comment:"请求URL"`
	RequestParams string `json:"requestParams" comment:"请求参数"`
	ResponseResult string `json:"responseResult" comment:"返回结果"`
	IPAddress     string `json:"ipAddress" comment:"IP地址"`
	UserAgent     string `json:"userAgent" comment:"用户代理"`
	Status        string `json:"status" comment:"操作状态: success/failed"`
	ErrorMsg      string `json:"errorMsg" comment:"错误信息"`
	ExecuteTime   int    `json:"executeTime" comment:"执行时长(ms)"`
}

type OperationLogPageListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	UserID   uint   `json:"userID" form:"userID"`
	Username string `json:"username" form:"username"`
	Module   string `json:"module" form:"module"`
	Operation string `json:"operation" form:"operation"`
	Status   string `json:"status" form:"status"`
}

type OperationLogPageListResp struct {
	List  []OperationLogPageListItem `json:"list"`
	Total int64                     `json:"total"`
}

type OperationLogPageListItem struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"userID"`
	Username      string `json:"username"`
	Module        string `json:"module"`
	Operation     string `json:"operation"`
	Method        string `json:"method"`
	RequestURL    string `json:"requestURL"`
	IPAddress     string `json:"ipAddress"`
	Status        string `json:"status"`
	ErrorMsg      string `json:"errorMsg"`
	ExecuteTime   int    `json:"executeTime"`
	CreatedAt     int64  `json:"createdAt"`
}
