package dtotag

type TagCreateReq struct {
	KbID  uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Name  string `json:"name" validate:"required" label:"标签名称"`
	Color string `json:"color" label:"标签颜色"`
}

type TagUpdateReq struct {
	ID    uint   `json:"id" validate:"required" label:"标签ID"`
	Name  string `json:"name" validate:"required" label:"标签名称"`
	Color string `json:"color" label:"标签颜色"`
}

type TagDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"标签ID"`
}

type TagListReq struct {
	KbID uint `json:"kbId" form:"kbId" validate:"required" label:"知识库ID"`
}
