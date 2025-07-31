package service

import (
	"apparatus/config"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
)

type OllamaServiceStruct struct {
	Client *api.Client
}

var OllamaService *OllamaServiceStruct

func (svc *OllamaServiceStruct) EmbedText(prompt string) ([]float32, error) {
	var count int = 0
	var embeddingRequest api.EmbeddingRequest
	embeddingRequest.Model = config.Config.OllamaEmbeddingsModel

	embeddingRequest.Prompt = prompt

	embeddingResponse, err := svc.Client.Embeddings(context.Background(), &embeddingRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings for segment %d: %w", count, err)
	}

	embeddingFloat32 := make([]float32, len(embeddingResponse.Embedding))
	for i, v := range embeddingResponse.Embedding {
		embeddingFloat32[i] = float32(v)
	}
	return embeddingFloat32, nil
}

func (svc *OllamaServiceStruct) Send(request *api.ChatRequest) (string, error) {
	client := svc.Client
	responseMessage := strings.Builder{}

	err := client.Chat(context.Background(), request, func(response api.ChatResponse) error {
		if response.Message.Content != "" {
			responseMessage.WriteString(response.Message.Content)
		}
		return nil
	})

	if err != nil {
		log.Printf("failed to send chat request: %v\n", err)
	}
	return responseMessage.String(), err
}
