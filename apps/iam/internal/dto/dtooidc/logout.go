package dtooidc

type LogoutReq struct {
	RefreshToken string `form:"refresh_token"`
	State        string `form:"state"`
}

type LogoutResp struct {
	RedirectURI string `json:"redirect_uri,omitempty"`
}
