package redis

import (
	"auth-api/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type Conn struct {
	client    *redis.Client
	keyPrefix string
}

func InitConnection(cfg *models.Config, log *zap.Logger) *Conn {
	log = log.Named("[REDIS]")
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host,
		DB:   cfg.Redis.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(fmt.Sprintf("connection: %v", err))
	}

	return &Conn{
		client:    client,
		keyPrefix: "",
	}
}

func (r *Conn) Get(ctx context.Context, key string) ([]byte, error) {
	key = r.keyPrefix + key
	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []byte{}, nil
		}
		return []byte{}, err
	}
	return result, nil
}

func (r *Conn) Set(ctx context.Context, key string, value []byte) error {
	key = r.keyPrefix + key
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *Conn) SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	key = r.keyPrefix + key
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Conn) Delete(ctx context.Context, key string) error {
	key = r.keyPrefix + key
	return r.client.Del(ctx, key).Err()
}

func (r *Conn) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *Conn) Close() error {
	return r.client.Close()
}
