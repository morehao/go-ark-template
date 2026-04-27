package model

import (
	"testing"
	"time"
)

func TestSsoSessionEntity_TableName(t *testing.T) {
	entity := SsoSessionEntity{}
	if entity.TableName() != TableNameSsoSession {
		t.Errorf("TableName() = %s, want %s", entity.TableName(), TableNameSsoSession)
	}
}

func TestSsoSessionEntity_Fields(t *testing.T) {
	now := time.Now()
	entity := SsoSessionEntity{
		ID:             1,
		SessionID:      "test-session-id",
		PersonID:       100,
		OrgID:          200,
		LoginTime:      now,
		LastActiveTime: now,
		ExpiresAt:      now.Add(24 * time.Hour),
	}

	if entity.ID != 1 {
		t.Errorf("ID = %d, want 1", entity.ID)
	}
	if entity.SessionID != "test-session-id" {
		t.Errorf("SessionID = %s, want test-session-id", entity.SessionID)
	}
	if entity.PersonID != 100 {
		t.Errorf("PersonID = %d, want 100", entity.PersonID)
	}
	if entity.OrgID != 200 {
		t.Errorf("OrgID = %d, want 200", entity.OrgID)
	}
	if entity.LoginTime != now {
		t.Errorf("LoginTime = %v, want %v", entity.LoginTime, now)
	}
	if entity.LastActiveTime != now {
		t.Errorf("LastActiveTime = %v, want %v", entity.LastActiveTime, now)
	}
	if entity.ExpiresAt != now.Add(24*time.Hour) {
		t.Errorf("ExpiresAt = %v, want %v", entity.ExpiresAt, now.Add(24*time.Hour))
	}
}