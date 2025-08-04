package repository_test

import (
	"apparatus/model"
	"apparatus/repository"
	"testing"
)

func TestSaveMemory(t *testing.T) {
	clientId := "rafael.solano"
	history := []string{"command1", "command2", "command3"}

	memory, err := model.NewRedisMemory(clientId)
	if err != nil {
		t.Fatalf("Failed to create new RedisMemory: %v", err)
	}
	memory.History = history

	err = repository.GlobalRedisMemoryRepository.SaveMemory(memory)
	if err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	retrievedMemory, err := repository.GlobalRedisMemoryRepository.GetMemory(clientId)
	if err != nil {
		t.Fatalf("Failed to retrieve memory: %v", err)
	}
	if retrievedMemory.ClientId != clientId {
		t.Fatalf("Expected clientId %s, got %s", clientId, retrievedMemory.ClientId)
	}
	if len(retrievedMemory.History) != len(history) {
		t.Fatalf("Expected history length %d, got %d", len(history), len(retrievedMemory.History))
	}
	for i, cmd := range history {
		if retrievedMemory.History[i] != cmd {
			t.Fatalf("Expected history command %s, got %s", cmd, retrievedMemory.History[i])
		}
	}
}
