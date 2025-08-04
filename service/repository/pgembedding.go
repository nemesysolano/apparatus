package repository

import (
	"apparatus/config"
	"apparatus/model"
	"log"

	"github.com/go-pg/pg"
	"github.com/pgvector/pgvector-go"
)

type PgEmbeddingRepositoryStruct struct {
	user     string
	password string
	database string
	addr     string
}

func NewPgEmbeddingRepository() *PgEmbeddingRepositoryStruct {
	return &PgEmbeddingRepositoryStruct{
		user:     config.Config.PostgresUser,
		password: config.Config.PostgresPassword,
		database: config.Config.PostgresDatabase,
		addr:     config.Config.PostgresHost,
	}
}

func (repository *PgEmbeddingRepositoryStruct) CosSimilarity(vector []float32) ([]model.PgEmbedding, error) {
	var embeddings []model.PgEmbedding
	db := pg.Connect(&pg.Options{
		User:     repository.user,
		Password: repository.password,
		Database: repository.database,
		Addr:     repository.addr,
	})
	defer db.Close()

	_, err := db.Query(&embeddings, "SELECT vector_id, content, embedding FROM vector_store ORDER BY embedding <=> ? LIMIT 5", pgvector.NewVector(vector))
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func (repository *PgEmbeddingRepositoryStruct) EuclideanDistance(vector []float32) ([]model.PgEmbedding, error) {
	var embeddings []model.PgEmbedding
	db := pg.Connect(&pg.Options{
		User:     repository.user,
		Password: repository.password,
		Database: repository.database,
		Addr:     repository.addr,
	})
	defer db.Close()

	_, err := db.Query(&embeddings, "SELECT vector_id, content, embedding FROM vector_store ORDER BY embedding <-> ? LIMIT 5", pgvector.NewVector(vector))
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func (repository *PgEmbeddingRepositoryStruct) SaveEmbedding(pgEmbedding model.PgEmbedding) error {
	db := pg.Connect(&pg.Options{
		User:     config.Config.PostgresUser,
		Password: config.Config.PostgresPassword,
		Database: config.Config.PostgresDatabase,
		Addr:     config.Config.PostgresHost,
	})
	defer db.Close()

	// Use explicit table name in the query
	_, err := db.Query(&pgEmbedding, "DELETE FROM vector_store WHERE vector_id = ?", pgEmbedding.VectorId)
	if err != nil {
		log.Printf("failed to delete existing embedding %s: %v", pgEmbedding.VectorId, err)
		return err
	}

	_, err = db.Query(&pgEmbedding, "INSERT INTO vector_store (vector_id, content, embedding) VALUES (?, ?, ?)", pgEmbedding.VectorId, pgEmbedding.Content, pgEmbedding.Embedding)
	if err != nil {
		log.Printf("failed to save embedding %s: %v", pgEmbedding.VectorId, err)
		return err
	}

	return nil
}

var GlobalPgEmbeddingRepository *PgEmbeddingRepositoryStruct
