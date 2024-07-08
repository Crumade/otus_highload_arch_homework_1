package server

import (
	"errors"
	"net/http"

	models "social_network/internal/model"
	utils "social_network/internal/pkg/utils"
	"social_network/internal/service"

	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", login(db))
	router.HandleFunc("POST /user/register", register(db))
	router.HandleFunc("GET /user/get/{id}", getUserByID(db))

	return router
}

func login(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		loginData := new(models.LoginRequest)

		err := utils.ParseJSON(r, loginData)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		token, err := service.Login(db, loginData)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusAccepted, token)
	}
}

func register(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)

		err := utils.ParseJSON(r, user)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if user.Password == "" {
			utils.WriteError(w, http.StatusBadRequest, errors.New("не передан пароль"))
			return
		}
		result, err := service.Register(db, user)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, result)
	}
}

func getUserByID(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		result, err := service.GetUser(db, id)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, result)
	}
}
