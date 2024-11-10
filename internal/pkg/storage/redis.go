package storage

import (
	"context"
	"log/slog"
	"net"

	"github.com/redis/go-redis/v9"
)

const (
	REDIS_HOST = "host.docker.internal"
	REDIS_PORT = "6379"
)

func NewRedisConnection() (*redis.Client, error) {

	cache := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})

	status := cache.Ping(context.Background())
	result, err := status.Result()
	if err != nil {
		slog.Error("redis status", "ping error", err.Error())
		return nil, err
	}
	slog.Info("redis status", "ping result", result)
	return cache, nil
}

func Warming(cache *redis.Client) {
	slog.Debug("process of warcming cache")
}
