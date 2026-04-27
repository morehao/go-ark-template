package svcoidc

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	svc := &logoutSvc{}

	hash1 := svc.hashToken("test-token")
	if hash1 == "" {
		t.Error("hashToken should return non-empty string")
	}

	hash2 := svc.hashToken("test-token")
	if hash1 != hash2 {
		t.Error("hashToken should return same hash for same input")
	}

	hash3 := svc.hashToken("different-token")
	if hash1 == hash3 {
		t.Error("hashToken should return different hash for different input")
	}
}