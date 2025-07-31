package model

import (
	"fmt"
	"regexp"

	"github.com/pgvector/pgvector-go"
)

type PgEmbedding struct {
	tableName struct{}        `pg:"vector_store"`
	VectorId  string          `json:"vector_id" pg:"vector_id,pk,type:VARCHAR(128),notnull"`
	Content   string          `json:"content" pg:"type:TEXT,notnull"`
	Embedding pgvector.Vector `json:"embedding" pg:"type:vector(1024),notnull"`
}

type RedisMemory struct {
	ClientId string
	History  []string
}

var clientIdValidator = regexp.MustCompile(`^\w+\.\w+$`)

func NewRedisMemory(clientId string) (*RedisMemory, error) {

	if !IsValidClientId(clientId) {
		return nil, fmt.Errorf("invalid clientId: %s must be two words separated by '.', max 64 characters", clientId)
	}

	return &RedisMemory{
		ClientId: clientId,
		History:  []string{},
	}, nil
}

func IsValidClientId(clientId string) bool {
	if len(clientId) == 0 || len(clientId) > 64 {
		return false
	}

	return clientIdValidator.MatchString(clientId)
}

func (history *RedisMemory) AddToHistory(item string) {
	if len(history.History) > 10 {
		history.History = history.History[1:] // Remove the oldest item
	}
	history.History = append(history.History, item)
}
