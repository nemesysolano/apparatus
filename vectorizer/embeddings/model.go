package embeddings

import (
	"github.com/pgvector/pgvector-go"
)

type SourceDocument struct {
	FileName        string          `json:"file_name"`
	InstitutionCode InstitutionCode `json:"institution_code"`
	Content         string          `json:"content"`
}

type Segment struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ParsedDocument struct {
	FileName        string          `json:"file_name"`
	InstitutionCode InstitutionCode `json:"institution_code"`
	Segments        []Segment       `json:"segments"`
}

type Embedding struct {
	VectorId  string
	Content   string
	Embedding []float64
}

type PgEmbedding struct {
	tableName struct{}        `pg:"vector_store"`
	VectorId  string          `json:"vector_id" pg:"vector_id,pk,type:VARCHAR(128),notnull"`
	Content   string          `json:"content" pg:"type:TEXT,notnull"`
	Embedding pgvector.Vector `json:"embedding" pg:"type:vector(1024),notnull"`
}

func (embedding Embedding) ToPgEmbedding() PgEmbedding {
	// Convert []float64 to []float32
	float32Embedding := make([]float32, len(embedding.Embedding))
	for i, v := range embedding.Embedding {
		float32Embedding[i] = float32(v)
	}
	return PgEmbedding{
		VectorId:  embedding.VectorId,
		Content:   embedding.Content,
		Embedding: pgvector.NewVector(float32Embedding),
	}
}
