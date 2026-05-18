package objapplication

type ApplicationBaseInfo struct {
	AppCode     string `json:"appCode" form:"appCode"`         // 应用编码
	AppName     string `json:"appName" form:"appName"`         // 应用名称
	AppType     string `json:"appType" form:"appType"`         // 应用类型: web-网页 app-移动端 mini-小程序
	CallbackUrl string `json:"callbackUrl" form:"callbackUrl"` // 回调URL
	Description string `json:"description" form:"description"` // 应用描述
	HomepageUrl string `json:"homepageUrl" form:"homepageUrl"` // 应用首页URL
	Logo        string `json:"logo" form:"logo"`               // 应用Logo
	Sequence   int32  `json:"sequence" form:"sequence"`     // 排序
	Status      string `json:"status" form:"status"`           // 状态: enabled-启用 disabled-停用
}
