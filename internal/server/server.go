package server

import (
	"log"
	"log/slog"
	"net/http"
	"social_network/internal/pkg/storage"
	"time"
)

func RunServer(pg *storage.PostgresDB, cache *storage.CacheDB) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        NewRouter(pg.Conn /*, cache.Conn*/),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	slog.Info("server running on port 8080")
	log.Fatal(s.ListenAndServe())
}
