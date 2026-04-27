package dtooidc

import (
	"encoding/json"
	"testing"
)

func TestAuthorizeReq(t *testing.T) {
	req := AuthorizeReq{
		ResponseType:        "code",
		ClientID:            "client123",
		RedirectURI:         "https://example.com/callback",
		Scope:               "openid profile",
		State:               "state123",
		CodeChallenge:       "challenge",
		CodeChallengeMethod: "S256",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed AuthorizeReq
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.ResponseType != req.ResponseType {
		t.Errorf("ResponseType mismatch: got %s, want %s", parsed.ResponseType, req.ResponseType)
	}
	if parsed.ClientID != req.ClientID {
		t.Errorf("ClientID mismatch: got %s, want %s", parsed.ClientID, req.ClientID)
	}
}

func TestAuthorizeResp(t *testing.T) {
	resp := AuthorizeResp{
		Code:  "authcode123",
		State: "state456",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed AuthorizeResp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.Code != resp.Code {
		t.Errorf("Code mismatch: got %s, want %s", parsed.Code, resp.Code)
	}
	if parsed.State != resp.State {
		t.Errorf("State mismatch: got %s, want %s", parsed.State, resp.State)
	}
}

func TestTokenReq(t *testing.T) {
	req := TokenReq{
		GrantType:    "authorization_code",
		Code:         "code123",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client123",
		ClientSecret: "secret",
		CodeVerifier: "verifier",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed TokenReq
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.GrantType != req.GrantType {
		t.Errorf("GrantType mismatch: got %s, want %s", parsed.GrantType, req.GrantType)
	}
}

func TestTokenResp(t *testing.T) {
	resp := TokenResp{
		AccessToken:  "access123",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "refresh123",
		IDToken:      "id123",
		Scope:        "openid profile",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed TokenResp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.AccessToken != resp.AccessToken {
		t.Errorf("AccessToken mismatch: got %s, want %s", parsed.AccessToken, resp.AccessToken)
	}
	if parsed.ExpiresIn != resp.ExpiresIn {
		t.Errorf("ExpiresIn mismatch: got %d, want %d", parsed.ExpiresIn, resp.ExpiresIn)
	}
}

func TestUserInfoResp(t *testing.T) {
	resp := UserInfoResp{
		Subject: "user123",
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "1234567890",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed UserInfoResp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.Subject != resp.Subject {
		t.Errorf("Subject mismatch: got %s, want %s", parsed.Subject, resp.Subject)
	}
	if parsed.Email != resp.Email {
		t.Errorf("Email mismatch: got %s, want %s", parsed.Email, resp.Email)
	}
}

func TestLogoutReq(t *testing.T) {
	req := LogoutReq{
		RefreshToken: "refresh123",
		State:        "state456",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed LogoutReq
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.RefreshToken != req.RefreshToken {
		t.Errorf("RefreshToken mismatch: got %s, want %s", parsed.RefreshToken, req.RefreshToken)
	}
}

func TestLogoutResp(t *testing.T) {
	resp := LogoutResp{
		RedirectURI: "https://example.com/loggedout",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed LogoutResp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed.RedirectURI != resp.RedirectURI {
		t.Errorf("RedirectURI mismatch: got %s, want %s", parsed.RedirectURI, resp.RedirectURI)
	}
}
