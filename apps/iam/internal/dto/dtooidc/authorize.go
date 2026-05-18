package dtooidc

type AuthorizeReq struct {
	ResponseType        string `form:"response_type" binding:"required"`         // 响应类型
	ClientID            string `form:"client_id" binding:"required"`             // 客户端ID
	RedirectURI         string `form:"redirect_uri" binding:"required"`           // 重定向URI
	Scope               string `form:"scope"`                                    // 作用域
	State               string `form:"state"`                                    // 状态
	CodeChallenge       string `form:"code_challenge"`                           // PKCE挑战
	CodeChallengeMethod string `form:"code_challenge_method"`                     // PKCE挑战方法
}

type AuthorizeResp struct {
	Code  string `json:"code"`  // 授权码
	State string `json:"state,omitempty"` // 状态
}
