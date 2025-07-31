package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	OllamaHost            string `env:"OLLAMA_HOST"`
	OllamaEmbeddingsModel string `env:"OLLAMA_EMBEDDINGS_MODEL"`
	OllamaChatModel       string `env:"OLLAMA_CHAT_MODEL"`
	OllamaEmbeddingsLog   bool   `env:"OLLAMA_EMBEDDINGS_LOG"`
	PostgresHost          string `env:"POSTGRES_HOST"`
	PostgresUser          string `env:"POSTGRES_USER"`
	PostgresPassword      string `env:"POSTGRES_PASSWORD"`
	PostgresDatabase      string `env:"POSTGRES_DATABASE"`
	PostgresSSLMode       string `env:"POSTGRES_SSLMODE"`
}

var Config ConfigStruct

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment variables loaded successfully")

	// Load environment variables into the Config struct
	Config.OllamaChatModel = os.Getenv("OLLAMA_CHAT_MODEL")
	Config.OllamaEmbeddingsModel = os.Getenv("OLLAMA_EMBEDDINGS_MODEL")
	Config.OllamaHost = os.Getenv("OLLAMA_HOST")
	Config.PostgresHost = os.Getenv("POSTGRES_HOST")
	Config.PostgresUser = os.Getenv("POSTGRES_USER")
	Config.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	Config.PostgresDatabase = os.Getenv("POSTGRES_DATABASE")
	Config.PostgresSSLMode = os.Getenv("POSTGRES_SSLMODE")
	Config.OllamaEmbeddingsLog = os.Getenv("OLLAMA_EMBEDDINGS_LOG") == "true"

	if Config.PostgresHost == "" || Config.PostgresUser == "" || Config.PostgresPassword == "" || Config.PostgresDatabase == "" {
		panic("One or more Postgres environment variables are not set")
	}

	if Config.OllamaHost == "" || Config.OllamaEmbeddingsModel == "" || Config.OllamaChatModel == "" {
		panic("One or more Ollama environment variables are not set")
	}

	log.Printf("Postgres configuration loaded: %+v", Config)
}
