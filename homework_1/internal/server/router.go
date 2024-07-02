package server

import "net/http"

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/login", login)
	router.HandleFunc("/user/register", register)
	router.HandleFunc("/user/get/{id}", getUserByID)

	return router
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
