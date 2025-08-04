package embeddings

import (
	"context"
	"fmt"
	"vectorizer/config"

	"github.com/ollama/ollama/api"
)

func defaultPromptLogger(prompt string) {
	if config.Config.OllamaEmbeddingsLog {
		fmt.Println(prompt)
	}
}

func EmbedParsedDocument(ctx context.Context, client *api.Client, parsedDocument *ParsedDocument) ([]Embedding, error) {
	var count int = 0
	var embeddingRequest api.EmbeddingRequest
	var embeddings []Embedding
	embeddingRequest.Model = config.Config.OllamaEmbeddingsModel

	for _, segment := range parsedDocument.Segments {
		if segment.Content == "" {
			continue // Skip empty segments
		}

		vectorId := fmt.Sprintf("%s %s", parsedDocument.FileName, segment.Title)
		prompt := fmt.Sprintf("%s. %s", vectorId, segment.Content)
		defaultPromptLogger(prompt)
		embeddingRequest.Prompt = prompt
		embeddingResponse, err := client.Embeddings(ctx, &embeddingRequest)

		if err != nil {
			return nil, fmt.Errorf("failed to get embeddings for segment %d: %w", count, err)
		}

		embeddings = append(embeddings, Embedding{
			VectorId:  vectorId,
			Content:   segment.Content,
			Embedding: embeddingResponse.Embedding,
		})

		count++
	}

	return embeddings, nil
}
