package strategy

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/iam/dao"
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/code"
)

type RegisterStrategyType string

const (
	RegisterStrategyOpen   RegisterStrategyType = "open"
	RegisterStrategyDomain RegisterStrategyType = "domain"
	RegisterStrategyInvite RegisterStrategyType = "invite"
	RegisterStrategySSO    RegisterStrategyType = "sso"
)

type RegisterStrategy interface {
	PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error)
	PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint, result *RegisterResult) error
	GetStrategyType() RegisterStrategyType
}

type RegisterRequest struct {
	Username   string
	Password   string
	Mobile     string
	Email      string
	RealName   string
	InviteCode string
	SSOType    string
	OpenID     string
}

type RegisterResult struct {
	TenantID     uint
	PersonID     uint
	PasswordHash string
	Status       model.UserStatus
	PersonExists bool
	Message      string
	InviteID     uint
}

type strategySelector struct {
}

func NewStrategySelector() *strategySelector {
	return &strategySelector{}
}

func (s *strategySelector) SelectStrategy(ctx *gin.Context, req *RegisterRequest) (RegisterStrategy, error) {
	orgEntity, err := getCurrentOrg(ctx)
	if err != nil {
		return nil, err
	}

	registerWay, err := s.getOrgConfigString(ctx, orgEntity.ID, model.OrgConfigKeyRegisterWay)
	if err != nil {
		return nil, code.GetError(code.AuthRegisterDisabled)
	}

	switch RegisterStrategyType(registerWay) {
	case RegisterStrategyOpen:
		return NewOpenStrategy(), nil
	case RegisterStrategyDomain:
		return NewDomainStrategy(), nil
	case RegisterStrategyInvite:
		return NewInviteStrategy(), nil
	case RegisterStrategySSO:
		return NewSSOStrategy(), nil
	default:
		return nil, code.GetError(code.AuthRegisterDisabled)
	}
}

func getCurrentOrg(ctx *gin.Context) (*model.OrganizationEntity, error) {
	domain := resolveDomain(ctx)
	if domain == "" {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	tenant, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
		Domain: domain,
		Status: model.TenantStatusEnabled,
	})
	if err != nil {
		return nil, code.GetError(code.AuthLoginError)
	}
	if tenant == nil || tenant.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}

	orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, tenant.OrgID)
	if err != nil {
		return nil, code.GetError(code.AuthLoginError)
	}
	if orgEntity == nil || orgEntity.ID == 0 {
		return nil, code.GetError(code.AuthOrgNotFoundError)
	}
	return orgEntity, nil
}

func (s *strategySelector) getOrgConfigString(ctx *gin.Context, orgID uint, configKey string) (string, error) {
	configEntity, err := dao.NewOrganizationConfigDao().GetByCond(ctx, &dao.OrganizationConfigCond{
		OrgID:     orgID,
		ConfigKey: configKey,
	})
	if err != nil {
		return "", err
	}
	if configEntity == nil || configEntity.ID == 0 {
		return "", nil
	}
	return configEntity.ConfigValue, nil
}

func resolveDomain(ctx *gin.Context) string {
	host := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Host"))
	if host == "" && ctx.Request != nil {
		host = strings.TrimSpace(ctx.Request.Host)
	}
	host = strings.TrimSpace(strings.Split(host, ",")[0])
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.Split(host, "/")[0]
	if strings.Count(host, ":") == 1 {
		host = strings.Split(host, ":")[0]
	}
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	return host
}
