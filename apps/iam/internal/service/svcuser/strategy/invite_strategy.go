package strategy

import (
	"github.com/gin-gonic/gin"
)

type inviteStrategy struct {
}

func NewInviteStrategy() RegisterStrategy {
	return &inviteStrategy{}
}

func (s *inviteStrategy) PreRegister(ctx *gin.Context, req *RegisterRequest) (*RegisterResult, error) {
	return &RegisterResult{}, nil
}

func (s *inviteStrategy) PostRegister(ctx *gin.Context, req *RegisterRequest, userID uint) error {
	return nil
}

func (s *inviteStrategy) GetStrategyType() RegisterStrategyType {
	return RegisterStrategyInvite
}
