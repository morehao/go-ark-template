package model

import (
	"testing"
)

func TestTokenEntity_TableName(t *testing.T) {
	if TableNameToken != "iam_token" {
		t.Errorf("TableNameToken = %v, want %v", TableNameToken, "iam_token")
	}
}

func TestTokenType_Values(t *testing.T) {
	if TokenTypeAccess != "access" {
		t.Errorf("TokenTypeAccess = %v, want %v", TokenTypeAccess, "access")
	}
	if TokenTypeRefresh != "refresh" {
		t.Errorf("TokenTypeRefresh = %v, want %v", TokenTypeRefresh, "refresh")
	}
	if TokenTypeID != "id" {
		t.Errorf("TokenTypeID = %v, want %v", TokenTypeID, "id")
	}
}