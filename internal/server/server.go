package server

import (
	"log"
	"net/http"
	"social_network/internal/pkg/storage"
	"time"
)

func RunServer() {

	connDB, err := storage.NewConnection()
	if err != nil {
		log.Fatal("Ошибка: не удалось подключиться к БД")
	}
	defer connDB.Close()

	storage.MigrateSchema(connDB)
	log.Println("DB connection success")

	err = storage.MigrateUsers(connDB)
	if err != nil {
		log.Fatal(err)
	}

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
