package svcapikey

import (
	"crypto/rand"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/config"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/internal/dto/dtoapikey"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type ApiKeySvc interface {
	Create(ctx *gin.Context, req *dtoapikey.ApiKeyCreateReq) (*dtoapikey.ApiKeyCreateResp, error)
	Delete(ctx *gin.Context, req *dtoapikey.ApiKeyDeleteReq) error
	List(ctx *gin.Context, req *dtoapikey.ApiKeyListReq) (*dtoapikey.ApiKeyListResp, error)
	Disable(ctx *gin.Context, req *dtoapikey.ApiKeyDisableReq) error
	Enable(ctx *gin.Context, req *dtoapikey.ApiKeyEnableReq) error
}

type apiKeySvc struct {
}

var _ ApiKeySvc = (*apiKeySvc)(nil)

func NewApiKeySvc() ApiKeySvc {
	return &apiKeySvc{}
}

func (svc *apiKeySvc) Create(ctx *gin.Context, req *dtoapikey.ApiKeyCreateReq) (*dtoapikey.ApiKeyCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	userID := gincontext.GetUserID(ctx)
	operatorID := gincontext.GetUserID(ctx)

	apiKeyPlain := svc.generateApiKey()

	aesEncryptor, err := gcrypto.NewAES(config.Conf.MasterKey)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.Create] NewAES fail, err:%v", err)
		return nil, code.GetError(code.ApiKeyCreateError)
	}

	encryptedApiKey, err := aesEncryptor.EncryptString(apiKeyPlain)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.Create] Encrypt API Key fail, err:%v", err)
		return nil, code.GetError(code.ApiKeyCreateError)
	}

	keyPrefix := "ark_" + apiKeyPlain[4:12]

	insertEntity := &model.ApiKeyEntity{
		TenantID:     tenantID,
		UserID:       userID,
		AppID:        req.AppID,
		KeyName:      req.KeyName,
		KeyPrefix:    keyPrefix,
		ApiKey:       encryptedApiKey,
		AccessPolicy: model.ApiKeyAccessPolicy(req.AccessPolicy),
		AllowedIPs:   req.AllowedIPs,
		Scopes:       req.Scopes,
		ExpiresAt:    req.ExpiresAt,
		Status:       model.ApiKeyStatusEnabled,
		CreatedBy:    operatorID,
		UpdatedBy:    operatorID,
	}

	if err := dao.NewApiKeyDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcapikey.Create] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ApiKeyCreateError)
	}

	return &dtoapikey.ApiKeyCreateResp{
		ID:        insertEntity.ID,
		KeyPrefix: keyPrefix,
		ApiKey:    apiKeyPlain,
		ExpiresAt: req.ExpiresAt,
	}, nil
}

func (svc *apiKeySvc) Delete(ctx *gin.Context, req *dtoapikey.ApiKeyDeleteReq) error {
	operatorID := gincontext.GetUserID(ctx)

	apiKeyEntity, err := dao.NewApiKeyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.Delete] GetByID fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyDeleteError)
	}
	if apiKeyEntity == nil || apiKeyEntity.ID == 0 {
		return code.GetError(code.ApiKeyNotExistError)
	}

	if err := dao.NewApiKeyDao().Delete(ctx, req.ID, operatorID); err != nil {
		glog.Errorf(ctx, "[svcapikey.Delete] Delete fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyDeleteError)
	}
	return nil
}

func (svc *apiKeySvc) List(ctx *gin.Context, req *dtoapikey.ApiKeyListReq) (*dtoapikey.ApiKeyListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)

	cond := &dao.ApiKeyCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID: tenantID,
		AppID:    req.AppID,
		KeyName:  req.KeyName,
		Status:   model.ApiKeyStatus(req.Status),
	}

	apiKeyList, total, err := dao.NewApiKeyDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.List] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ApiKeyGetPageListError)
	}

	aesDecryptor, err := gcrypto.NewAES(config.Conf.MasterKey)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.List] NewAES fail, err:%v", err)
		return nil, code.GetError(code.ApiKeyGetPageListError)
	}

	list := make([]dtoapikey.ApiKeyListItem, 0, len(apiKeyList))
	for _, v := range apiKeyList {
		apiKeyPlain, err := aesDecryptor.DecryptString(v.ApiKey)
		if err != nil {
			glog.Errorf(ctx, "[svcapikey.List] Decrypt API Key fail, err:%v, id:%d", err, v.ID)
			apiKeyPlain = ""
		}

		list = append(list, dtoapikey.ApiKeyListItem{
			ID:           v.ID,
			AppID:        v.AppID,
			KeyName:      v.KeyName,
			KeyPrefix:    v.KeyPrefix,
			ApiKey:       apiKeyPlain,
			Scopes:       v.Scopes,
			AccessPolicy: string(v.AccessPolicy),
			AllowedIPs:   v.AllowedIPs,
			Status:       string(v.Status),
			LastUsedAt:   v.LastUsedAt,
			ExpiresAt:    v.ExpiresAt,
		})
	}

	return &dtoapikey.ApiKeyListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *apiKeySvc) Disable(ctx *gin.Context, req *dtoapikey.ApiKeyDisableReq) error {
	operatorID := gincontext.GetUserID(ctx)

	apiKeyEntity, err := dao.NewApiKeyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.Disable] GetByID fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyUpdateError)
	}
	if apiKeyEntity == nil || apiKeyEntity.ID == 0 {
		return code.GetError(code.ApiKeyNotExistError)
	}

	updateMap := map[string]any{
		"status":     model.ApiKeyStatusDisabled,
		"updated_by": operatorID,
	}
	if err := dao.NewApiKeyDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcapikey.Disable] UpdateMap fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyUpdateError)
	}
	return nil
}

func (svc *apiKeySvc) Enable(ctx *gin.Context, req *dtoapikey.ApiKeyEnableReq) error {
	operatorID := gincontext.GetUserID(ctx)

	apiKeyEntity, err := dao.NewApiKeyDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcapikey.Enable] GetByID fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyUpdateError)
	}
	if apiKeyEntity == nil || apiKeyEntity.ID == 0 {
		return code.GetError(code.ApiKeyNotExistError)
	}

	updateMap := map[string]any{
		"status":     model.ApiKeyStatusEnabled,
		"updated_by": operatorID,
	}
	if err := dao.NewApiKeyDao().UpdateMap(ctx, req.ID, updateMap); err != nil {
		glog.Errorf(ctx, "[svcapikey.Enable] UpdateMap fail, err:%v, id:%d", err, req.ID)
		return code.GetError(code.ApiKeyUpdateError)
	}
	return nil
}

func (svc *apiKeySvc) generateApiKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 16)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			n = big.NewInt(int64(i % len(charset)))
		}
		result[i] = charset[n.Int64()]
	}
	return "ark_" + string(result)
}
