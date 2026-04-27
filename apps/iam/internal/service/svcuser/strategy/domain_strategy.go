package strategy

import (
	"github.com/gin-gonic/gin"
)

type domainStrategy struct {
}

func NewDomainStrategy() RegisterStrategy {
	return &domainStrategy{}
}

func (s *domainStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	return &RegisterResult{}, nil
}

func (s *domainStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint) error {
	return nil
}

func (s *domainStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyDomain
}
