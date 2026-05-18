package svcoidc

import (
	"testing"
)

func TestNewUserInfoSvc(t *testing.T) {
	svc := NewUserInfoSvc()
	if svc == nil {
		t.Fatal("NewUserInfoSvc should not return nil")
	}
	if _, ok := svc.(UserInfoSvc); !ok {
		t.Error("NewUserInfoSvc should return UserInfoSvc interface")
	}
}
