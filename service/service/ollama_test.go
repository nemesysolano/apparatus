package service_test

import (
	"apparatus/initialization"
	"apparatus/service"
	"testing"
)

func TestEmbedText(test *testing.T) {
	initialization.Init()
	embedding, err := service.OllamaService.EmbedText("Test prompt")
	if err != nil {
		test.Fatalf("Failed to embed text: %v", err)
	}
	if len(embedding) == 0 {
		test.Fatalf("Expected non-empty embedding, got %d elements", len(embedding))
	}
	test.Logf("Embedding length: %d", len(embedding))

	// TODO: Check that 1st 10 elements are not all zero
	for i := 0; i < 10 && i < len(embedding); i++ {
		if embedding[i] == 0 {
			test.Errorf("Embedding element %d is zero", i)
		}
	}
}
