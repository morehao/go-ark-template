package svcoidc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	svc := &authorizeSvc{}

	code, err := svc.generateCode()
	assert.Nil(t, err)
	assert.NotEmpty(t, code)
	assert.Len(t, code, 64)

	code2, err := svc.generateCode()
	assert.Nil(t, err)
	assert.NotEqual(t, code, code2)
}

func TestValidatePKCE_EmptyParams(t *testing.T) {
	svc := &authorizeSvc{}
	ctx := context.Background()

	err := svc.ValidatePKCE(ctx, "", "", "")
	assert.Error(t, err)
	assert.Equal(t, "pkce parameters required", err.Error())
}

func TestValidatePKCE_ValidParams(t *testing.T) {
	svc := &authorizeSvc{}
	ctx := context.Background()

	err := svc.ValidatePKCE(ctx, "challenge", "verifier", "S256")
	assert.Nil(t, err)
}