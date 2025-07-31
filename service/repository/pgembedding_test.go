package repository_test

import (
	"apparatus/initialization"
	"apparatus/model"
	"apparatus/repository"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestScalarMultipleByCosineSimilarity(test *testing.T) {
	initialization.Init()
	// Create test vectors with 1024 dimensions as expected by the database
	vector1 := make([]float32, 1024)
	vector2 := make([]float32, 1024)

	// Set some sample values (keeping it simple for testing)
	for i := 0; i < 5; i++ {
		vector1[i] = float32(i + 1)       // [1, 2, 3, 4, 5, 0, 0, ...]
		vector2[i] = float32((i + 1) * 2) // [2, 4, 6, 8, 10, 0, 0, ...] - scalar multiple
	}

	firstEmbedding := model.PgEmbedding{
		VectorId:  "test_vector_1",
		Content:   "Test vector 1",
		Embedding: pgvector.NewVector(vector1),
	}
	err := repository.GlobalPgEmbeddingRepository.SaveEmbedding(firstEmbedding)
	if err != nil {
		test.Fatalf("Failed to save first embedding: %v", err)
	}

	secondEmbedding := model.PgEmbedding{
		VectorId:  "test_vector_2",
		Content:   "Test vector 2",
		Embedding: pgvector.NewVector(vector2),
	}
	err = repository.GlobalPgEmbeddingRepository.SaveEmbedding(secondEmbedding)
	if err != nil {
		test.Fatalf("Failed to save second embedding: %v", err)
	}

	queryVector := make([]float32, 1024)
	for i := 0; i < 5; i++ {
		queryVector[i] = float32(i + 1) // Same as vector1
	}

	embeddings, err := repository.GlobalPgEmbeddingRepository.CosSimilarity(queryVector)
	if err != nil {
		test.Fatalf("Failed to retrieve embeddings: %v", err)
	}
	if len(embeddings) == 0 {
		test.Fatalf("Expected at least some embeddings, got %d", len(embeddings))
	}

	// Check that our test embeddings are present in the results
	foundFirst := embeddings[0].VectorId == firstEmbedding.VectorId && embeddings[0].Content == firstEmbedding.Content
	foundSecond := embeddings[1].VectorId == secondEmbedding.VectorId && embeddings[1].Content == secondEmbedding.Content

	if !foundFirst {
		test.Errorf("Expected to find first test embedding (test_vector_1) in results")
	}
	if !foundSecond {
		test.Errorf("Expected to find second test embedding (test_vector_2) in results")
	}
}

func Init() {
	panic("unimplemented")
}

func TestOrthogonalByCosineSimilarity(test *testing.T) {
	initialization.Init()
	// Create test vectors with 1024 dimensions as expected by the database
	v := []float32{2, 3, 5, 7, -3}
	u := []float32{11, 13, 17, 19, 93}

	vector1 := make([]float32, 1024)
	vector2 := make([]float32, 1024)

	// Set some sample values (keeping it simple for testing)
	for i := 0; i < len(v); i++ {
		vector1[i] = v[i]
		vector2[i] = u[i]
	}

	firstEmbedding := model.PgEmbedding{
		VectorId:  "test_vector_1",
		Content:   "Test vector 1",
		Embedding: pgvector.NewVector(vector1),
	}
	err := repository.GlobalPgEmbeddingRepository.SaveEmbedding(firstEmbedding)
	if err != nil {
		test.Fatalf("Failed to save first embedding: %v", err)
	}

	secondEmbedding := model.PgEmbedding{
		VectorId:  "test_vector_2",
		Content:   "Test vector 2",
		Embedding: pgvector.NewVector(vector2),
	}
	err = repository.GlobalPgEmbeddingRepository.SaveEmbedding(secondEmbedding)
	if err != nil {
		test.Fatalf("Failed to save second embedding: %v", err)
	}

	queryVector := make([]float32, 1024)
	copy(queryVector, u)

	embeddings, err := repository.GlobalPgEmbeddingRepository.CosSimilarity(queryVector)
	if err != nil {
		test.Fatalf("Failed to retrieve embeddings: %v", err)
	}
	if len(embeddings) == 0 {
		test.Fatalf("Expected at least some embeddings, got %d", len(embeddings))
	}

	// Check that our test embeddings are present in the results
	foundSecond := embeddings[0].VectorId == secondEmbedding.VectorId && embeddings[0].Content == secondEmbedding.Content
	foundFirst := embeddings[1].VectorId == firstEmbedding.VectorId && embeddings[1].Content == firstEmbedding.Content

	if !foundFirst {
		test.Errorf("Expected to find first test embedding (test_vector_1) in results")
	}
	if !foundSecond {
		test.Errorf("Expected to find second test embedding (test_vector_2) in results")
	}
}
