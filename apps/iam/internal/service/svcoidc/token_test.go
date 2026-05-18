package svcoidc

import (
	"testing"

	"github.com/morehao/golib/gcrypto"
)

func TestAccessTokenExpireDuration(t *testing.T) {
	if AccessTokenExpireDuration != 1*60*60*1e9 {
		t.Errorf("AccessTokenExpireDuration should be 1 hour, got %v", AccessTokenExpireDuration)
	}
}

func TestRefreshTokenExpireDuration(t *testing.T) {
	if RefreshTokenExpireDuration != 7*24*60*60*1e9 {
		t.Errorf("RefreshTokenExpireDuration should be 7 days, got %v", RefreshTokenExpireDuration)
	}
}

func TestTokenIssuer(t *testing.T) {
	if TokenIssuer != "iam" {
		t.Errorf("TokenIssuer should be 'iam', got %s", TokenIssuer)
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	hash, err := gcrypto.GeneratePasswordHash("password")
	if err != nil {
		t.Fatalf("GeneratePasswordHash failed: %v", err)
	}
	if err := gcrypto.ComparePasswordHash(hash, "password"); err != nil {
		t.Errorf("ComparePasswordHash failed: %v", err)
	}
}
