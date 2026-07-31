package svcmessage

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtomessage"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/dbaccess/gormdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type MessageSvc interface {
	List(ctx *gin.Context, req *dtomessage.MessageListReq) (*dtomessage.MessageListResp, error)
	Delete(ctx *gin.Context, req *dtomessage.MessageDeleteReq) error
	Search(ctx *gin.Context, req *dtomessage.MessageSearchReq) (*dtomessage.MessageSearchResp, error)
}

type messageSvc struct {
}

var _ MessageSvc = (*messageSvc)(nil)

func NewMessageSvc() MessageSvc {
	return &messageSvc{}
}

func (svc *messageSvc) List(ctx *gin.Context, req *dtomessage.MessageListReq) (*dtomessage.MessageListResp, error) {
	cond := &dao.MessageCond{
		SessionID: req.SessionID,
	}
	dataList, err := dao.NewMessageDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcmessage.List] dao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MessageListError)
	}
	list := make([]dtomessage.MessageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtomessage.MessageListItem{
			ID:         v.ID,
			Role:       v.Role,
			Content:    v.Content,
			Metadata:   v.Metadata,
			TokenCount: v.TokenCount,
			CreatedAt:  v.CreatedAt.Unix(),
		})
	}
	return &dtomessage.MessageListResp{
		List: list,
	}, nil
}

func (svc *messageSvc) Search(ctx *gin.Context, req *dtomessage.MessageSearchReq) (*dtomessage.MessageSearchResp, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	cond := &dao.MessageCond{
		BaseCond: &gormdao.BaseCond{
			Page:     page,
			PageSize: pageSize,
		},
		SessionID: req.SessionID,
		Keyword:   req.Keyword,
	}
	dataList, total, err := dao.NewMessageDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcmessage.Search] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.MessageSearchError)
	}

	list := make([]dtomessage.MessageSearchItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtomessage.MessageSearchItem{
			ID:         v.ID,
			SessionID:  v.SessionID,
			Role:       v.Role,
			Content:    v.Content,
			Metadata:   v.Metadata,
			TokenCount: v.TokenCount,
			CreatedAt:  v.CreatedAt.Unix(),
		})
	}
	return &dtomessage.MessageSearchResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (svc *messageSvc) Delete(ctx *gin.Context, req *dtomessage.MessageDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewMessageDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcmessage.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.MessageDeleteError)
	}
	return nil
}
