package svcknowledge

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/client/storage"
	"github.com/morehao/goark/ragforge/config"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtoknowledge"
	"github.com/morehao/goark/ragforge/internal/pipeline"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/dbaccess/gormdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type KnowledgeSvc interface {
	CreateFromFile(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateFileReq) (*dtoknowledge.KnowledgeCreateResp, error)
	CreateFromURL(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateURLReq) (*dtoknowledge.KnowledgeCreateResp, error)
	CreateManual(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateManualReq) (*dtoknowledge.KnowledgeCreateResp, error)
	Update(ctx *gin.Context, req *dtoknowledge.KnowledgeUpdateReq) error
	Delete(ctx *gin.Context, req *dtoknowledge.KnowledgeDeleteReq) error
	Detail(ctx *gin.Context, req *dtoknowledge.KnowledgeDetailReq) (*dtoknowledge.KnowledgeDetailResp, error)
	PageList(ctx *gin.Context, req *dtoknowledge.KnowledgePageListReq) (*dtoknowledge.KnowledgePageListResp, error)
	Reparse(ctx *gin.Context, req *dtoknowledge.KnowledgeReparseReq) error
}

type knowledgeSvc struct {
	fileSvc    storage.FileService
	httpClient *http.Client
}

var _ KnowledgeSvc = (*knowledgeSvc)(nil)

func NewKnowledgeSvc() KnowledgeSvc {
	basePath := "./data/files"
	fileSvc, _, err := storage.NewFileServiceFromStorageConfig("", &config.Conf.Storage, basePath)
	if err != nil {
		fileSvc = storage.NewLocalFileService(basePath)
	}
	return &knowledgeSvc{
		fileSvc:    fileSvc,
		httpClient: http.DefaultClient,
	}
}

func (svc *knowledgeSvc) CreateFromFile(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateFileReq) (*dtoknowledge.KnowledgeCreateResp, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] FormFile fail, err:%v", err)
		return nil, code.GetError(code.KnowledgeCreateError)
	}
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	insertEntity := &model.KnowledgeEntity{
		KbID:        req.KbID,
		TenantID:    tenantID,
		Type:        model.KnowledgeTypeFile,
		Title:       req.Title,
		FileSize:    file.Size,
		ParseStatus: model.ParseStatusPending,
		CreatorID:   userID,
	}
	if err := dao.NewKnowledgeDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	fileURL, err := svc.fileSvc.SaveFile(ctx, file, tenantID, fmt.Sprintf("%d", insertEntity.ID))
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] fileSvc SaveFile fail, err:%v", err)
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	if err := dao.NewKnowledgeDao().UpdateMap(ctx, insertEntity.ID, map[string]interface{}{
		"file_url": fileURL,
	}); err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] dao UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	src, err := file.Open()
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] file Open fail, err:%v", err)
		return nil, code.GetError(code.KnowledgeCreateError)
	}
	contentBytes, err := io.ReadAll(src)
	_ = src.Close()
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromFile] ReadAll fail, err:%v", err)
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	go pipeline.ProcessAndEmbedContent(ctx.Copy(), req.KbID, insertEntity.ID, tenantID, string(contentBytes))

	return &dtoknowledge.KnowledgeCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *knowledgeSvc) CreateFromURL(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateURLReq) (*dtoknowledge.KnowledgeCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	resp, err := svc.httpClient.Get(req.SourceURL)
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromURL] http Get fail, err:%v, url:%s", err, req.SourceURL)
		return nil, code.GetError(code.KnowledgeCreateError)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		glog.Errorf(ctx, "[svcknowledge.CreateFromURL] http status not ok, code:%d, url:%s", resp.StatusCode, req.SourceURL)
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromURL] ReadAll fail, err:%v, url:%s", err, req.SourceURL)
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	insertEntity := &model.KnowledgeEntity{
		KbID:        req.KbID,
		TenantID:    tenantID,
		Type:        model.KnowledgeTypeURL,
		Title:       req.Title,
		SourceURL:   req.SourceURL,
		Content:     string(bodyBytes),
		ParseStatus: model.ParseStatusCompleted,
		CreatorID:   userID,
	}
	if err := dao.NewKnowledgeDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateFromURL] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeCreateError)
	}

	go pipeline.ProcessAndEmbedContent(ctx.Copy(), req.KbID, insertEntity.ID, tenantID, string(bodyBytes))

	return &dtoknowledge.KnowledgeCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *knowledgeSvc) CreateManual(ctx *gin.Context, req *dtoknowledge.KnowledgeCreateManualReq) (*dtoknowledge.KnowledgeCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)
	insertEntity := &model.KnowledgeEntity{
		KbID:        req.KbID,
		TenantID:    tenantID,
		Type:        model.KnowledgeTypeManual,
		Title:       req.Title,
		Content:     req.Content,
		ParseStatus: model.ParseStatusCompleted,
		CreatorID:   userID,
	}
	if err := dao.NewKnowledgeDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcknowledge.CreateManual] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeCreateError)
	}
	go pipeline.ProcessAndEmbedContent(ctx.Copy(), req.KbID, insertEntity.ID, tenantID, req.Content)
	return &dtoknowledge.KnowledgeCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *knowledgeSvc) Update(ctx *gin.Context, req *dtoknowledge.KnowledgeUpdateReq) error {
	updateEntity := &model.KnowledgeEntity{
		Title:   req.Title,
		Content: req.Content,
	}
	if err := dao.NewKnowledgeDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcknowledge.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.KnowledgeUpdateError)
	}
	return nil
}

func (svc *knowledgeSvc) Delete(ctx *gin.Context, req *dtoknowledge.KnowledgeDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewKnowledgeDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcknowledge.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.KnowledgeDeleteError)
	}
	return nil
}

func (svc *knowledgeSvc) Detail(ctx *gin.Context, req *dtoknowledge.KnowledgeDetailReq) (*dtoknowledge.KnowledgeDetailResp, error) {
	detailEntity, err := dao.NewKnowledgeDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.KnowledgeNotExistError)
	}
	resp := &dtoknowledge.KnowledgeDetailResp{
		ID:          detailEntity.ID,
		KbID:        detailEntity.KbID,
		TenantID:    detailEntity.TenantID,
		Type:        detailEntity.Type,
		Title:       detailEntity.Title,
		Content:     detailEntity.Content,
		FileURL:     detailEntity.FileURL,
		SourceURL:   detailEntity.SourceURL,
		ParseStatus: detailEntity.ParseStatus,
		FileSize:    detailEntity.FileSize,
		CreatorID:   detailEntity.CreatorID,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *knowledgeSvc) PageList(ctx *gin.Context, req *dtoknowledge.KnowledgePageListReq) (*dtoknowledge.KnowledgePageListResp, error) {
	cond := &dao.KnowledgeCond{
		BaseCond: &gormdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		KbID:          req.KbID,
		KnowledgeType: req.Type,
		ParseStatus:   req.ParseStatus,
	}
	dataList, total, err := dao.NewKnowledgeDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcknowledge.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KnowledgeGetPageListError)
	}
	list := make([]dtoknowledge.KnowledgePageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtoknowledge.KnowledgePageListItem{
			ID:          v.ID,
			KbID:        v.KbID,
			TenantID:    v.TenantID,
			Type:        v.Type,
			Title:       v.Title,
			ParseStatus: v.ParseStatus,
			FileSize:    v.FileSize,
			CreatorID:   v.CreatorID,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtoknowledge.KnowledgePageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *knowledgeSvc) Reparse(ctx *gin.Context, req *dtoknowledge.KnowledgeReparseReq) error {
	updateMap := map[string]any{
		"parse_status": model.ParseStatusPending,
	}
	if err := dao.NewKnowledgeDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcknowledge.Reparse] dao UpdateMap fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.KnowledgeUpdateError)
	}
	return nil
}
