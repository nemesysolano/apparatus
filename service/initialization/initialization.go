package initialization

import (
	"apparatus/config"
	"apparatus/repository"
	"apparatus/service"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-redis/redis"
	"github.com/ollama/ollama/api"
)

func init() {
	Init()
}

func Init() {
	initPgEmbeddingRepository()
	initRedisMemoryRepository()
	initOllamaService()
	initApparatusService()
}

func initPgEmbeddingRepository() {
	repository.GlobalPgEmbeddingRepository = repository.NewPgEmbeddingRepository()
}

func initRedisMemoryRepository() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	repository.GlobalRedisMemoryRepository = &repository.RedisMemoryRepositoryStruct{
		RedisHost: config.Config.RedisHost, // Assuming RedisHost is used for Redis connection
		Client:    client,
	}
}

func initOllamaService() {
	url, err := url.Parse(config.Config.OllamaHost)
	if err != nil {
		log.Printf("Invalid OLLAMA_HOST URL: %s", config.Config.OllamaHost)
		os.Exit(4)
	}
	client := api.NewClient(url, http.DefaultClient)
	service.OllamaService = &service.OllamaServiceStruct{Client: client}
}

func initApparatusService() {
	service.ApparatusService = &service.ApparatusServiceStruct{
		OllamaService:         service.OllamaService,
		RedisMemoryRepository: repository.GlobalRedisMemoryRepository,
		PgEmbeddingRepository: repository.GlobalPgEmbeddingRepository,
	}

}
