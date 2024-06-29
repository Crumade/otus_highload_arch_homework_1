package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/login", login)
	router.HandleFunc("/user/register", register)
	router.HandleFunc("/user/get/{id}", getUserByID)

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

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is login"))
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is registration"))
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is user by id"))
}
