package service

import (
	"apparatus/config"
	"apparatus/model"
	"apparatus/repository"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
)

type ApparatusServiceStruct struct {
	OllamaService         *OllamaServiceStruct
	RedisMemoryRepository *repository.RedisMemoryRepositoryStruct
	PgEmbeddingRepository *repository.PgEmbeddingRepositoryStruct
}

var ApparatusService *ApparatusServiceStruct
var (
	FALSE = false
	TRUE  = true
)

func (svc *ApparatusServiceStruct) Query(question *ApparatusQuestion) (*ApparatusAnswer, error) {
	redisMemory := svc.GetMemory(question.UserID)
	if redisMemory == nil {
		return nil, fmt.Errorf("no memory found for user %s", question.UserID)
	}

	messages, err := svc.NewChatMessages(question.Prompt, redisMemory)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat messages: %v", err)
	}

	chatRequest := svc.NewChatRequest(messages)
	response, err := svc.OllamaService.Send(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat response: %v", err)
	}

	answer := &ApparatusAnswer{
		Answer: response,
		Score:  1.0, // Default score, can be adjusted based on response quality
	}

	redisMemory.AddToHistory(response)
	svc.RedisMemoryRepository.SaveMemory(redisMemory)
	return answer, nil
}

func (svc *ApparatusServiceStruct) NewChatMessages(prompt string, redisMemory *model.RedisMemory) ([]api.Message, error) {
	context, err := svc.Context(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %v", err)
	}

	messages := []api.Message{
		{Role: "system", Content: "Responde únicamente a preguntas relevantes al contexto proporcionado y que estén escritas en español. Rechaza cortésmente cualquier otra pregunta."},
		{Role: "system", Content: "Las respuestas deben generarse solo a partir del contexto proporcionado; no uses información externa."},
		{Role: "system", Content: context},
	}

	for _, history := range redisMemory.History {
		messages = append(messages, api.Message{
			Role:    "assistant",
			Content: history,
		})
	}
	messages = append(messages, api.Message{
		Role:    "user",
		Content: prompt,
	})
	return messages, nil
}

func (svc *ApparatusServiceStruct) NewChatRequest(messages []api.Message) *api.ChatRequest {
	chatRequest := &api.ChatRequest{
		Model:    config.Config.OllamaChatModel,
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":    0.0,
			"repeat_last_n":  2,
			"repeat_penalty": 1.8,
			"top_k":          5,
			"top_p":          0.5,
		},
		Stream: &FALSE,
	}

	return chatRequest
}

func (svc *ApparatusServiceStruct) Context(prompt string) (string, error) {
	stringsSet := make(map[string]struct{})

	stringBuilder := strings.Builder{}
	stringBuilder.WriteString("Context for the question (in Spanish): ")

	embeddings, err := svc.OllamaService.EmbedText(prompt)
	if err != nil {
		return "", err
	}

	pgCosEmbeddings, err := svc.PgEmbeddingRepository.CosSimilarity(embeddings)
	if err != nil {
		return "", err
	}

	for _, embedding := range pgCosEmbeddings {
		log.Println("Adding cos embedding content to context: ", embedding.VectorId)
		stringBuilder.WriteString(embedding.Content + "\n")
		stringsSet[embedding.Content] = struct{}{}
	}

	pgEuclideanEmbeddings, err := svc.PgEmbeddingRepository.EuclideanDistance(embeddings)
	if err != nil {
		return "", err
	}

	for _, embedding := range pgEuclideanEmbeddings {
		if _, exists := stringsSet[embedding.Content]; exists {
			log.Println("Skipping existing content in context: ", embedding.VectorId)
			continue // Skip if content already exists
		}
		log.Println("Adding euclidean embedding content to context: ", embedding.VectorId)
		stringBuilder.WriteString(embedding.Content + "\n")
	}

	return stringBuilder.String(), nil
}

func (svc *ApparatusServiceStruct) SaveMemory(redisMemory *model.RedisMemory) error {
	return svc.RedisMemoryRepository.SaveMemory(redisMemory)
}

func (svc *ApparatusServiceStruct) GetMemory(clientId string) *model.RedisMemory {
	memory, _ := svc.RedisMemoryRepository.GetMemory(clientId)
	if memory == nil {
		return &model.RedisMemory{
			ClientId: clientId,
			History:  []string{},
		}
	}
	return memory
}
