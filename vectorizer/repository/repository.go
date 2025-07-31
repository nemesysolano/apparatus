package repository

import (
	"fmt"
	"log"
	"vectorizer/config"
	"vectorizer/embeddings"

	"github.com/go-pg/pg/v10"
)

func SaveAllEmbeddings(embeddings []embeddings.Embedding) error {
	db := pg.Connect(&pg.Options{
		User:     config.Config.PostgresUser,
		Password: config.Config.PostgresPassword,
		Database: config.Config.PostgresDatabase,
		Addr:     config.Config.PostgresHost,
	})
	defer db.Close()

	for _, embedding := range embeddings {
		pgEmbedding := embedding.ToPgEmbedding()
		_, err := db.Model(&pgEmbedding).Where("vector_id = ?", pgEmbedding.VectorId).Delete()
		if err != nil {
			log.Printf("failed to delete existing embedding %s: %v", embedding.VectorId, err)
			return err
		}
		_, err = db.Model(&pgEmbedding).Insert()
		if err != nil {
			log.Printf("failed to save embedding %s: %v", embedding.VectorId, err)
			return err
		}
		log.Println(fmt.Sprintf("Saved embedding for vector ID: %s", embedding.VectorId))
	}

	return nil
}
