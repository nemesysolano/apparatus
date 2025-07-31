package repository

import (
	"apparatus/model"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisMemoryRepositoryStruct struct {
	RedisHost string
	Client    *redis.Client
}

var GlobalRedisMemoryRepository *RedisMemoryRepositoryStruct

func (repo *RedisMemoryRepositoryStruct) SaveMemory(memory *model.RedisMemory) error {
	data, err := json.Marshal(memory)
	if err != nil {
		return err
	}

	err = repo.Client.Set(memory.ClientId, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RedisMemoryRepositoryStruct) GetMemory(clientId string) (*model.RedisMemory, error) {
	data, err := repo.Client.Get(clientId).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("memory not found for client ID: %s", clientId)
		}
		return nil, err
	}

	var memory model.RedisMemory
	err = json.Unmarshal([]byte(data), &memory)
	if err != nil {
		return nil, err
	}

	return &memory, nil
}
