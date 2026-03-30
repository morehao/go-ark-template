package svcuser

import (
	"testing"

	"github.com/morehao/golib/gcrypto"
)

func TestGeneratePassword(t *testing.T) {
	email := "admin@platform.com"
	plainPassword := "pwd" + email

	hash, err := gcrypto.GeneratePasswordHash(plainPassword)
	if err != nil {
		t.Fatalf("GeneratePasswordHash failed: %v", err)
	}

	t.Logf("Email: %s", email)
	t.Logf("Plain password: %s", plainPassword)
	t.Logf("Password hash: %s", hash)

	if err := gcrypto.ComparePasswordHash(hash, plainPassword); err != nil {
		t.Errorf("ComparePasswordHash failed: %v", err)
	}
}
