package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	models "social_network/internal/model"
	utils "social_network/internal/pkg/utils"
	"social_network/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func NewRouter(db *sqlx.DB, cache *redis.Client) *chi.Mux {

	rateLimit := RateLimit{
		cache: cache,
	}

	mux := chi.NewRouter()
	mux.Use(LogRequest)
	mux.Use(rateLimit.Handle)

	mux.Post("/login", login(db))
	mux.Route("/user", func(r chi.Router) {
		r.Post("/register", register(db))
		r.Get("/get/{id}", getUserByID(db))
		r.Get("/search", searchUser(db))
	})

	mux.Route("/post", func(r chi.Router) {
		r.Get("/feed", getPostFeed(db, cache))
		r.Get("/get/{id}", getPostByID(db))
		r.Post("/create", createPost(db))
		r.Put("/update", updatePost(db))
		r.Put("/delete/{id}", deletePost(db))
	})

	mux.Route("/friend", func(r chi.Router) {
		r.Put("/set/{user_id}", setFriend(db))
		r.Put("/delete/{user_id}", deleteFriend(db))

	})
	return mux
}

func login(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginData := new(models.LoginRequest)

		err := utils.ParseJSON(r, loginData)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}

		token, err := service.Login(db, loginData)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
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
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}
		if user.Password == "" {
			utils.WriteError(w, http.StatusBadRequest, *r, errors.New("не передан пароль"))
			return
		}
		result, err := service.Register(db, user)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
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
			utils.WriteError(w, http.StatusBadRequest, *r, err)
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
			utils.WriteError(w, http.StatusBadRequest, *r, errors.New("отсутствуют GET параметры"))
			return
		}

		result, err := service.SearchUser(db, firstName, lastName)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
	}
}

func getPostFeed(db *sqlx.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			offset int = 0
			limit  int = 10
			err    error
		)

		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}

		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}

		if offset < 0 || limit < 1 {
			utils.WriteError(w, http.StatusBadRequest, *r, errors.New("not valid params"))
			return
		}

		result, err := service.GetPostFeed(db, cache, offset, limit)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
	}
}

func getPostByID(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		result, err := service.GetPost(db, id)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)
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
		id := r.PathValue("id")
		result, err := service.DeletePost(db, id)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, *r, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, result)

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
