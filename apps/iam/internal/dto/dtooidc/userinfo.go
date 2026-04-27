package dtooidc

type UserInfoResp struct {
	Subject string `json:"sub"`     // 主题(用户ID)
	Name    string `json:"name,omitempty"`    // 名称
	Email   string `json:"email,omitempty"`   // 邮箱
	Phone   string `json:"phone,omitempty"`   // 手机号
}
