package dtooidc

type TokenReq struct {
	GrantType    string `form:"grant_type" binding:"required"` // 授权类型
	Code         string `form:"code"`                          // 授权码
	RedirectURI  string `form:"redirect_uri"`                  // 重定向URI
	ClientID     string `form:"client_id"`                    // 客户端ID
	ClientSecret string `form:"client_secret"`                // 客户端密钥
	CodeVerifier string `form:"code_verifier"`                // PKCE验证码
	RefreshToken string `form:"refresh_token"`                // 刷新令牌
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`   // 访问令牌
	TokenType    string `json:"token_type"`     // 令牌类型
	ExpiresIn    int    `json:"expires_in"`     // 过期时间
	RefreshToken string `json:"refresh_token,omitempty"` // 刷新令牌
	IDToken      string `json:"id_token,omitempty"`       // ID令牌
	Scope        string `json:"scope,omitempty"`          // 作用域
}

type TokenRefreshReq struct {
	RefreshToken string `form:"refresh_token" binding:"required"` // 刷新令牌
	ClientID     string `form:"client_id" binding:"required"`     // 客户端ID
	ClientSecret string `form:"client_secret" binding:"required"` // 客户端密钥
}
