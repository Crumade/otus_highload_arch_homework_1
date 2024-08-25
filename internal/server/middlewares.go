package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			slog.Error("dump request", "Error", err.Error())
		}
		RequestID := uuid.New().String()
		r.Header.Add("X-Request-ID", RequestID)
		slog.Info("log entry request", "request_uri", r.RequestURI, "request_id", RequestID, "request_body", string(dump))
		next.ServeHTTP(w, r)
	})
}

type RateLimit struct {
	cache *redis.Client
}

func (r RateLimit) Limit(ctx context.Context, now time.Time, clientID string) (map[string]string, bool, error) {
	var num int64
	err := r.cache.Get(ctx, clientID).Err()
	slog.Debug("redis vals", "res", strconv.Itoa(int(num)))
	if err == redis.Nil {
		num, err = r.cache.Incr(ctx, clientID).Result()
		if err != nil {
			panic(err)
		}
		slog.Debug("redis vals 1", "num", strconv.Itoa(int(num)))
		if num > 5 {
			slog.Debug("redis vals 2", "num", strconv.Itoa(int(num)))
			num = 0
			result := map[string]string{"header": "rate"}
			return result, false, nil
		}
		_, err = r.cache.Expire(ctx, clientID, time.Second*30).Result()
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		num, err = r.cache.Incr(ctx, clientID).Result()
		if err != nil {
			panic(err)
		}
		slog.Debug("redis vals 3", "num", strconv.Itoa(int(num)))
		if num > 5 {
			slog.Debug("redis vals 4", "num", strconv.Itoa(int(num)))
			num = 0
			result := map[string]string{"header": "rate"}
			return result, false, nil
		}
	}
	num = 0
	result := map[string]string{"header": "rate"}
	return result, true, nil
}

func (r RateLimit) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
		headers, ok, err := r.Limit(rr.Context(), time.Now().UTC(), "localhost")
		if err != nil {
			log.Println(err)
			slog.Error("rate limiter", "status", err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		for k, v := range headers {
			w.Header().Set(k, v)
		}

		if !ok {
			slog.Warn("rate limiter", "too many requests", "sorry")
			w.WriteHeader(http.StatusTooManyRequests)

			return
		}

		next.ServeHTTP(w, rr)
	})
}
