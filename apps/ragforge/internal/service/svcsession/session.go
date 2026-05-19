package svcsession

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtosession"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type SessionSvc interface {
	Create(ctx *gin.Context, req *dtosession.SessionCreateReq) (*dtosession.SessionCreateResp, error)
	Update(ctx *gin.Context, req *dtosession.SessionUpdateReq) error
	Delete(ctx *gin.Context, req *dtosession.SessionDeleteReq) error
	Detail(ctx *gin.Context, req *dtosession.SessionDetailReq) (*dtosession.SessionDetailResp, error)
	PageList(ctx *gin.Context, req *dtosession.SessionPageListReq) (*dtosession.SessionPageListResp, error)
	GenerateTitle(ctx *gin.Context, req *dtosession.SessionGenerateTitleReq) (*dtosession.SessionGenerateTitleResp, error)
	Stop(ctx *gin.Context, req *dtosession.SessionStopReq) error
}

type sessionSvc struct {
}

var _ SessionSvc = (*sessionSvc)(nil)

func NewSessionSvc() SessionSvc {
	return &sessionSvc{}
}

func (svc *sessionSvc) generateTitle(ctx *gin.Context, content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	llm := engine.GetGlobalFactory().GetLLM()
	if llm == nil {
		return ""
	}
	prompt := fmt.Sprintf("基于以下内容，生成一个简洁的标题（不超过20个字）：\n\n%s", content)
	resp, err := llm.ChatCompletion(ctx, &engine.ChatCompletionRequest{
		Messages: []engine.ChatMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		MaxTokens:   50,
	})
	if err != nil || resp == nil || strings.TrimSpace(resp.Content) == "" {
		glog.Errorf(ctx, "[svcsession.generateTitle] LLM fail, err:%v", err)
		return ""
	}
	title := strings.TrimSpace(resp.Content)
	if len([]rune(title)) > 100 {
		title = string([]rune(title)[:100])
	}
	return title
}

func (svc *sessionSvc) Create(ctx *gin.Context, req *dtosession.SessionCreateReq) (*dtosession.SessionCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	title := req.Title
	if title == "" && req.Content != "" {
		if generatedTitle := svc.generateTitle(ctx, req.Content); generatedTitle != "" {
			title = generatedTitle
		} else {
			title = truncateContent(req.Content, 100)
		}
	}

	insertEntity := &model.SessionEntity{
		TenantID: tenantID,
		UserID:   userID,
		KbID:     req.KbID,
		Title:    title,
	}
	if err := dao.NewSessionDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcsession.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.SessionCreateError)
	}
	return &dtosession.SessionCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func truncateContent(content string, maxLen int) string {
	runes := []rune(strings.TrimSpace(content))
	if len(runes) <= maxLen {
		return string(runes)
	}
	return string(runes[:maxLen])
}

func (svc *sessionSvc) Update(ctx *gin.Context, req *dtosession.SessionUpdateReq) error {
	updateEntity := &model.SessionEntity{
		Title: req.Title,
	}
	if err := dao.NewSessionDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcsession.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.SessionUpdateError)
	}
	return nil
}

func (svc *sessionSvc) Delete(ctx *gin.Context, req *dtosession.SessionDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewSessionDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcsession.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.SessionDeleteError)
	}
	return nil
}

func (svc *sessionSvc) Detail(ctx *gin.Context, req *dtosession.SessionDetailReq) (*dtosession.SessionDetailResp, error) {
	detailEntity, err := dao.NewSessionDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcsession.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.SessionGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.SessionGetDetailError)
	}
	resp := &dtosession.SessionDetailResp{
		ID:          detailEntity.ID,
		TenantID:    detailEntity.TenantID,
		UserID:      detailEntity.UserID,
		KbID:        detailEntity.KbID,
		Title:       detailEntity.Title,
		Description: detailEntity.Description,
		IsPinned:    detailEntity.IsPinned,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *sessionSvc) GenerateTitle(ctx *gin.Context, req *dtosession.SessionGenerateTitleReq) (*dtosession.SessionGenerateTitleResp, error) {
	sessionEntity, err := dao.NewSessionDao().GetByID(ctx, req.SessionID)
	if err != nil {
		glog.Errorf(ctx, "[svcsession.GenerateTitle] dao GetByID fail, err:%v, req:%d", err, req.SessionID)
		return nil, code.GetError(code.SessionGenerateTitleError)
	}
	if sessionEntity == nil || sessionEntity.ID == 0 {
		return nil, code.GetError(code.SessionGenerateTitleError)
	}

	if sessionEntity.Title != "" {
		return &dtosession.SessionGenerateTitleResp{Title: sessionEntity.Title}, nil
	}

	messageList, err := dao.NewMessageDao().GetListByCond(ctx, &dao.MessageCond{SessionID: req.SessionID})
	if err != nil {
		glog.Errorf(ctx, "[svcsession.GenerateTitle] GetListByCond fail, err:%v", err)
		return nil, code.GetError(code.SessionGenerateTitleError)
	}

	firstContent := ""
	if len(messageList) > 0 {
		firstContent = messageList[0].Content
	}

	title := ""
	if firstContent != "" {
		if generatedTitle := svc.generateTitle(ctx, firstContent); generatedTitle != "" {
			title = generatedTitle
		} else {
			runes := []rune(strings.TrimSpace(firstContent))
			if len(runes) > 100 {
				title = string(runes[:100])
			} else {
				title = string(runes)
			}
		}
	}

	updateEntity := &model.SessionEntity{Title: title}
	if err := dao.NewSessionDao().UpdateByID(ctx, req.SessionID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcsession.GenerateTitle] UpdateByID fail, err:%v", err)
	}

	return &dtosession.SessionGenerateTitleResp{Title: title}, nil
}

func (svc *sessionSvc) Stop(ctx *gin.Context, req *dtosession.SessionStopReq) error {
	return nil
}

func (svc *sessionSvc) PageList(ctx *gin.Context, req *dtosession.SessionPageListReq) (*dtosession.SessionPageListResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)
	cond := &dao.SessionCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID: tenantID,
		UserID:   userID,
		KbID:     req.KbID,
	}
	dataList, total, err := dao.NewSessionDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcsession.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.SessionGetPageListError)
	}
	list := make([]dtosession.SessionPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtosession.SessionPageListItem{
			ID:          v.ID,
			TenantID:    v.TenantID,
			UserID:      v.UserID,
			KbID:        v.KbID,
			Title:       v.Title,
			Description: v.Description,
			IsPinned:    v.IsPinned,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtosession.SessionPageListResp{
		List:  list,
		Total: total,
	}, nil
}
