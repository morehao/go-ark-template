package strategy

import (
	"github.com/gin-gonic/gin"
)

type openStrategy struct {
}

func NewOpenStrategy() RegisterStrategy {
	return &openStrategy{}
}

func (s *openStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	return &RegisterResult{}, nil
}

func (s *openStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint) error {
	return nil
}

func (s *openStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyOpen
}
