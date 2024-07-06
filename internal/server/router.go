package server

import (
	"log"
	"net/http"
	"strconv"

	httpResponse "soc_net/internal/pkg"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", login)
	router.HandleFunc("POST /user/register", register)
	router.HandleFunc("GET /user/get/{id}", getUserByID)

	return router
}

func login(w http.ResponseWriter, r *http.Request) {

	resp := httpResponse.Response{}
	resp.Code = 0

	log.Println(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is login"))
}

func register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is registration"))
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is user by id " + id))
}
