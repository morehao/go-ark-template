package objuser

type PersonBaseInfo struct {
	AvatarUrl    string `json:"avatarUrl" form:"avatarUrl"`       // 头像URL
	BirthDate    int64  `json:"birthDate" form:"birthDate"`       // 出生日期
	Email        string `json:"email" form:"email"`               // 邮箱
	Gender       string `json:"gender" form:"gender"`             // 性别: male-男 female-女 unknown-未知
	Mobile       string `json:"mobile" form:"mobile"`             // 手机号
	PasswordHash string `json:"passwordHash" form:"passwordHash"` // 密码哈希(不存储盐值,盐值在应用层生成)
	RealName     string `json:"realName" form:"realName"`         // 真实姓名
	Remark       string `json:"remark" form:"remark"`             // 备注
	Wechat       string `json:"wechat" form:"wechat"`             // 微信号
}
