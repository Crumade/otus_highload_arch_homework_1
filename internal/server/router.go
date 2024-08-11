package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("GET /user/search", searchUser(db))
	router.HandleFunc("GET /post/feed", getPostFeed(db))
	router.HandleFunc("GET /post/get/{id}", getPostByID(db))
	router.HandleFunc("POST /post/create", createPost(db))
	router.HandleFunc("PUT /post/update", updatePost(db))
	router.HandleFunc("DELETE /post/delete/{id}", deletePost(db))
	router.HandleFunc("PUT /friend/set/{user_id}", setFriend(db))
	router.HandleFunc("PUT /friend/delete/{user_id}", deleteFriend(db))

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
		w.Header().Add("Content-Type", "application/json")
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
		w.Header().Add("Content-Type", "application/json")
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
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
	}
}

func searchUser(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		if firstName == "" || lastName == "" {
			utils.WriteError(w, http.StatusBadRequest, errors.New("отсутствуют GET параметры"))
			return
		}

		result, err := service.SearchUser(db, firstName, lastName)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
	}
}

func getPostFeed(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			offset int = 0
			limit  int = 10
			err    error
		)

		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		if offset < 0 || limit < 1 {
			utils.WriteError(w, http.StatusBadRequest, errors.New("not valid params"))
			return
		}

		result, err := service.GetPostFeed(db, offset, limit)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
	}
}

func getPostByID(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//id := r.PathValue("id")
		log.Println(db)
		w.Write([]byte("метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}

func createPost(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(db)
		w.Write([]byte("поcт получен, метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}

func updatePost(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(db)
		w.Write([]byte("поcт получен, метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}

func deletePost(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//id := r.PathValue("id")
		log.Println(db)
		w.Write([]byte("поcт на удаление  получен, метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}

func setFriend(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//id := r.PathValue("user_id")
		log.Println(db)
		w.Write([]byte("запрос в друзья поулчен, метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}

func deleteFriend(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//id := r.PathValue("user_id")
		log.Println(db)
		w.Write([]byte("запрос на удаление получен, метод не реализован"))
		// result, err := service.GetUser(db, id)
		// if err != nil {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		// w.Header().Add("Content-Type", "application/json")
		// utils.WriteJSON(w, http.StatusOK, result)

	}
}
