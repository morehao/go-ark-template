package dtooidc

type LogoutReq struct {
	RefreshToken string `form:"refresh_token"` // 刷新令牌
	State        string `form:"state"`         // 状态
}

type LogoutResp struct {
	RedirectURI string `json:"redirect_uri,omitempty"` // 重定向URI
}
