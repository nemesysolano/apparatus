package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	RedisHost             string `env:"REDIS_HOST"`  // Assuming Redis host is also needed
	ServerPort            string `env:"SERVER_PORT"` // Added for server port configuration
}

var Config ConfigStruct

func init() {
	// Check if we're running in test mode
	if strings.Contains(os.Args[0], ".test") || strings.Contains(os.Args[0], "go-build") {
		// Running in test mode - try .env.test first
		err := godotenv.Load(".env.test")
		if err != nil {
			// If .env.test doesn't exist, try ../.env.test
			err = godotenv.Load("../.env.test")
			if err != nil {
				log.Printf("Warning: Could not load test .env file: %v", err)
				log.Println("Falling back to regular .env file...")
				loadRegularEnv()
			} else {
				log.Println("Test environment variables loaded successfully from ../.env.test")
			}
		} else {
			log.Println("Test environment variables loaded successfully from .env.test")
		}
	} else {
		// Regular mode - load .env
		loadRegularEnv()
	}

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
	Config.RedisHost = os.Getenv("REDIS_HOST")
	Config.ServerPort = os.Getenv("SERVER_PORT")

	if Config.PostgresHost == "" || Config.PostgresUser == "" || Config.PostgresPassword == "" || Config.PostgresDatabase == "" {
		panic("One or more Postgres environment variables are not set")
	}

	if Config.OllamaHost == "" || Config.OllamaEmbeddingsModel == "" || Config.OllamaChatModel == "" {
		panic(fmt.Sprintf("One or more Ollama environment variables are not set %s, %s, %s", Config.OllamaHost, Config.OllamaEmbeddingsModel, Config.OllamaChatModel))
	}

	if Config.ServerPort != "" {
		Config.ServerPort = func(portStr string) string {
			if _, err := strconv.Atoi(portStr); err != nil {
				panic("SERVER_PORT must be an integer")
			}
			return portStr
		}(Config.ServerPort)
	}

	log.Printf("Postgres configuration loaded: %+v", Config)
}

func loadRegularEnv() {
	// Try to load .env file from multiple locations
	err := godotenv.Load()
	if err != nil {
		// If .env not found in current directory, try parent directory
		err = godotenv.Load("../.env")
		if err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
			log.Println("Continuing with environment variables from system...")
		} else {
			log.Println("Environment variables loaded successfully from ../.env")
		}
	} else {
		log.Println("Environment variables loaded successfully from .env")
	}
}
