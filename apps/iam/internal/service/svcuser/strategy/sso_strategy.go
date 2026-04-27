package strategy

import (
	"github.com/gin-gonic/gin"
)

type ssoStrategy struct {
}

func NewSSOStrategy() RegisterStrategy {
	return &ssoStrategy{}
}

func (s *ssoStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	return &RegisterResult{}, nil
}

func (s *ssoStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint, result *RegisterResult) error {
	return nil
}

func (s *ssoStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategySSO
}
