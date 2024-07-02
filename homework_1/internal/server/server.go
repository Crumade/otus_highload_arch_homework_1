package server

import (
	"log"
	"net/http"
	"time"
)

func RunServer() {
	router := NewRouter()

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
