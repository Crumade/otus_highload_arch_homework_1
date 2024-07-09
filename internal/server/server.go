package server

import (
	"log"
	"net/http"
	"social_network/internal/pkg/storage"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunServer() {

	connDB, err := storage.NewConnection()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД!")
	}
	defer connDB.Close()

	driver, err := postgres.WithInstance(connDB.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Instance error: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("New DB Instance error: " + err.Error())
	}
	if err := m.Up(); err != nil {
		log.Fatal("Up migrations error: " + err.Error())
	}

	log.Println("DB connection success")

	router := NewRouter(connDB)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("starting server on :8080")
	log.Fatal(s.ListenAndServe())
}
