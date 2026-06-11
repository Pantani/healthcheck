package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	client *redis.Client
}

func Init(ctx context.Context, rawURL string) (*Database, error) {
	options, err := redisOptions(rawURL)
	if err != nil {
		return nil, err
	}
	db := &Database{client: redis.NewClient(options)}
	if err := db.client.Ping(ctx).Err(); err != nil {
		_ = db.client.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	return db, nil
}

func (db *Database) SaveData(ctx context.Context, entity, key string, value interface{}) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode value for %s.%s: %w", entity, key, err)
	}
	if err := db.client.HSet(ctx, entity, key, encoded).Err(); err != nil {
		return fmt.Errorf("save value for %s.%s: %w", entity, key, err)
	}
	return nil
}

func (db *Database) GetData(ctx context.Context, entity, key string, result interface{}) error {
	encoded, err := db.client.HGet(ctx, entity, key).Result()
	if err != nil {
		return fmt.Errorf("get value for %s.%s: %w", entity, key, err)
	}
	if err := json.Unmarshal([]byte(encoded), result); err != nil {
		return fmt.Errorf("decode value for %s.%s: %w", entity, key, err)
	}
	return nil
}

func (db *Database) Close() error {
	return db.client.Close()
}

func redisOptions(rawURL string) (*redis.Options, error) {
	if strings.Contains(rawURL, "://") {
		options, err := redis.ParseURL(rawURL)
		if err != nil {
			return nil, fmt.Errorf("parse redis url: %w", err)
		}
		return options, nil
	}
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("redis url cannot be empty")
	}
	return &redis.Options{Addr: rawURL}, nil
}
