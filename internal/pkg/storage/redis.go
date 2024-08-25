package storage

import (
	"context"
	"log/slog"
	"net"
	models "social_network/internal/model"

	"github.com/redis/go-redis/v9"
)

const (
	REDIS_HOST = "host.docker.internal"
	REDIS_PORT = "6379"
)

type CacheDB struct {
	Conn *redis.Client
}

func (cache *CacheDB) NewRedisConnection() error {

	cache.Conn = redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})

	status := cache.Conn.Ping(context.Background())
	result, err := status.Result()
	if err != nil {
		slog.Error("redis status", "ping error", err.Error())
		return err
	}
	slog.Info("redis status", "ping result", result)
	return nil
}

func (cache *CacheDB) Warming() {
	slog.Debug("process of warcming cache")
}

func (cache *CacheDB) GetPosts(limit int, offset int) (*[]models.Post, error) {

	return nil, nil
}
