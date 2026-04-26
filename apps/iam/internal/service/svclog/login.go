package svclog

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type LoginLogSvc interface {
	Create(ctx *gin.Context, req *LoginLogCreateReq) (*LoginLogCreateResp, error)
	PageList(ctx *gin.Context, req *LoginLogPageListReq) (*LoginLogPageListResp, error)
}

type loginLogSvc struct {
}

var _ LoginLogSvc = (*loginLogSvc)(nil)

func NewLoginLogSvc() LoginLogSvc {
	return &loginLogSvc{}
}

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
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
	UserID      uint   `json:"userID" form:"userID"`
	Username    string `json:"username" form:"username"`
	LoginType   string `json:"loginType" form:"loginType"`
	LoginStatus string `json:"loginStatus" form:"loginStatus"`
}

type LoginLogPageListResp struct {
	List  []LoginLogPageListItem `json:"list"`
	Total int64                `json:"total"`
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

func (svc *loginLogSvc) Create(ctx *gin.Context, req *LoginLogCreateReq) (*LoginLogCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	insertEntity := &model.LoginLogEntity{
		TenantID:     tenantID,
		UserID:       req.UserID,
		Username:     req.Username,
		LoginType:    req.LoginType,
		LoginStatus:  req.LoginStatus,
		LoginMessage: req.LoginMessage,
		IPAddress:    req.IPAddress,
		Location:     req.Location,
		Browser:      req.Browser,
		OS:           req.OS,
		CreatedBy:    operatorID,
		UpdatedBy:    operatorID,
	}

	if err := dao.NewLoginLogDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svclog.LoginLogCreate] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogCreateError)
	}

	return &LoginLogCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *loginLogSvc) PageList(ctx *gin.Context, req *LoginLogPageListReq) (*LoginLogPageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &dao.LoginLogCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID:    tenantID,
		UserID:      req.UserID,
		Username:    req.Username,
		LoginType:   req.LoginType,
		LoginStatus: req.LoginStatus,
	}

	logList, total, err := dao.NewLoginLogDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svclog.LoginLogPageList] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.LoginLogGetPageListError)
	}

	list := make([]LoginLogPageListItem, 0, len(logList))
	for _, v := range logList {
		list = append(list, LoginLogPageListItem{
			ID:           v.ID,
			UserID:       v.UserID,
			Username:     v.Username,
			LoginType:    v.LoginType,
			LoginStatus:  v.LoginStatus,
			LoginMessage: v.LoginMessage,
			IPAddress:    v.IPAddress,
			Location:     v.Location,
			Browser:      v.Browser,
			OS:           v.OS,
		})
	}

	return &LoginLogPageListResp{
		List:  list,
		Total: total,
	}, nil
}
