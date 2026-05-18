package model

import (
	"testing"
)

func TestApplicationEntity_OAuthFields(t *testing.T) {
	entity := ApplicationEntity{
		ClientID:          "test-client-id",
		ClientSecret:      "test-client-secret",
		ClientType:        "web",
		PkceRequired:      true,
		AllowedScopes:     "openid,profile,email",
		AllowedCallbacks: `["https://callback.example.com/oauth"]`,
	}

	if entity.ClientID != "test-client-id" {
		t.Errorf("ClientID = %s, want test-client-id", entity.ClientID)
	}
	if entity.ClientSecret != "test-client-secret" {
		t.Errorf("ClientSecret = %s, want test-client-secret", entity.ClientSecret)
	}
	if entity.ClientType != "web" {
		t.Errorf("ClientType = %s, want web", entity.ClientType)
	}
	if !entity.PkceRequired {
		t.Error("PkceRequired should be true")
	}
	if entity.AllowedScopes != "openid,profile,email" {
		t.Errorf("AllowedScopes = %s, want openid,profile,email", entity.AllowedScopes)
	}
	if entity.AllowedCallbacks != `["https://callback.example.com/oauth"]` {
		t.Errorf("AllowedCallbacks = %s, want JSON array", entity.AllowedCallbacks)
	}
}

func TestClientType_Constants(t *testing.T) {
	if ClientTypeWeb != "web" {
		t.Errorf("ClientTypeWeb = %s, want web", ClientTypeWeb)
	}
	if ClientTypeApp != "app" {
		t.Errorf("ClientTypeApp = %s, want app", ClientTypeApp)
	}
	if ClientTypeSpa != "spa" {
		t.Errorf("ClientTypeSpa = %s, want spa", ClientTypeSpa)
	}
	if ClientTypeMini != "mini" {
		t.Errorf("ClientTypeMini = %s, want mini", ClientTypeMini)
	}
}