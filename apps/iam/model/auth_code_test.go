package model

import (
	"testing"
	"time"
)

func TestAuthCodeEntity_TableName(t *testing.T) {
	entity := AuthCodeEntity{}
	if entity.TableName() != TableNameAuthCode {
		t.Errorf("TableName() = %s, want %s", entity.TableName(), TableNameAuthCode)
	}
}

func TestAuthCodeEntity_Fields(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(5 * time.Minute)
	usedAt := now.Add(2 * time.Minute)

	entity := AuthCodeEntity{
		ID:                  1,
		Code:                "test-auth-code",
		ClientID:            "test-client-id",
		PersonID:            100,
		TenantID:            10,
		OrgID:               200,
		RedirectURI:         "https://example.com/callback",
		Scope:               "openid,profile",
		State:               "test-state",
		CodeChallenge:       "test-code-challenge",
		CodeChallengeMethod: "S256",
		ExpiresAt:           expiresAt,
		Used:                false,
		UsedAt:              &usedAt,
	}

	if entity.ID != 1 {
		t.Errorf("ID = %d, want 1", entity.ID)
	}
	if entity.Code != "test-auth-code" {
		t.Errorf("Code = %s, want test-auth-code", entity.Code)
	}
	if entity.ClientID != "test-client-id" {
		t.Errorf("ClientID = %s, want test-client-id", entity.ClientID)
	}
	if entity.PersonID != 100 {
		t.Errorf("PersonID = %d, want 100", entity.PersonID)
	}
	if entity.TenantID != 10 {
		t.Errorf("TenantID = %d, want 10", entity.TenantID)
	}
	if entity.OrgID != 200 {
		t.Errorf("OrgID = %d, want 200", entity.OrgID)
	}
	if entity.RedirectURI != "https://example.com/callback" {
		t.Errorf("RedirectURI = %s, want https://example.com/callback", entity.RedirectURI)
	}
	if entity.Scope != "openid,profile" {
		t.Errorf("Scope = %s, want openid,profile", entity.Scope)
	}
	if entity.State != "test-state" {
		t.Errorf("State = %s, want test-state", entity.State)
	}
	if entity.CodeChallenge != "test-code-challenge" {
		t.Errorf("CodeChallenge = %s, want test-code-challenge", entity.CodeChallenge)
	}
	if entity.CodeChallengeMethod != "S256" {
		t.Errorf("CodeChallengeMethod = %s, want S256", entity.CodeChallengeMethod)
	}
	if !entity.ExpiresAt.Equal(expiresAt) {
		t.Errorf("ExpiresAt = %v, want %v", entity.ExpiresAt, expiresAt)
	}
	if entity.Used != false {
		t.Errorf("Used = %v, want false", entity.Used)
	}
	if !entity.UsedAt.Equal(usedAt) {
		t.Errorf("UsedAt = %v, want %v", entity.UsedAt, usedAt)
	}
}
