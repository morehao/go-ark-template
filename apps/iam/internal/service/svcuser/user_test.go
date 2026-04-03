package svcuser

import (
	"testing"

	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
	"github.com/morehao/goark/pkg/testsetup"
	"github.com/morehao/golib/gcrypto"
	"github.com/morehao/golib/gutil"
	"github.com/stretchr/testify/assert"
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

func TestUserList(t *testing.T) {
	testsetup.Initialize(testsetup.AppNameIam)
	defer testsetup.Close(testsetup.AppNameIam)

	ctx := testsetup.NewContext()
	svc := NewUserSvc()
	res, err := svc.PageList(ctx, &dtouser.UserPageListReq{})
	assert.Nil(t, err)
	t.Logf("res: %s", gutil.ToJsonString(res))
}
