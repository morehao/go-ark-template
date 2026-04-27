package svcoidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSessionID(t *testing.T) {
	svc := &ssoSessionSvc{}

	sessionID, err := svc.generateSessionID()
	assert.Nil(t, err)
	assert.NotEmpty(t, sessionID)
	assert.Len(t, sessionID, 64)

	sessionID2, err := svc.generateSessionID()
	assert.Nil(t, err)
	assert.NotEqual(t, sessionID, sessionID2)
}
