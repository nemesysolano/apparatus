package model

import "testing"

func TestNewRedisMemory(t *testing.T) {
	clientId := "00000"
	_, err := NewRedisMemory(clientId)
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	clientId = "invalid_client_id"
	_, err = NewRedisMemory(clientId)
	if err == nil {
		t.Fatalf("Expected error for invalid clientId, got nil")
	}

	clientId = ""
	_, err = NewRedisMemory(clientId)
	if err == nil {
		t.Fatalf("Expected error for invalid clientId, got nil")
	}

	clientId = "valid.client.id"
	_, err = NewRedisMemory(clientId)
	if err == nil {
		t.Fatalf("Expected no error for valid clientId, got %v", err)
	}

	clientId = "rafael.solano"
	_, err = NewRedisMemory(clientId)
	if err != nil {
		t.Fatalf("Expected no error for valid clientId, got %v", err)
	}
}
