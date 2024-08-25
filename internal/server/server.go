package server

import (
	"log"
	"log/slog"
	"net/http"
	"social_network/internal/pkg/storage"
)

func RunServer(s *http.Server) {

	pg := new(storage.PostgresDB)
	err := pg.NewConnection()
	if err != nil {
		log.Fatal("DB connection failure")
	}
	defer pg.Conn.Close()

	cache := new(storage.CacheDB)
	err = cache.NewRedisConnection()
	if err != nil {
		log.Fatalf("Redis error: " + err.Error())
	}

	storage.MigrateSchema(pg.Conn)
	slog.Info("DB connection success")

	err = storage.MigrateUsers(pg.Conn)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.MigratePosts(pg.Conn)
	if err != nil {
		log.Fatal(err)
	}

	cache.Warming()

	s.Handler = NewRouter(pg.Conn, cache.Conn)

	slog.Info("server running on port 8080")
	log.Fatal(s.ListenAndServe())
}
