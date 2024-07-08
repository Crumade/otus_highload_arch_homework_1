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
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/Crumade/Documents/projects/otus_highload_arch_homework_1/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
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
