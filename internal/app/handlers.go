package app

import (
	"errors"
	"net/http"
	models "social_network/internal/model"
	utils "social_network/internal/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type ServiceHandler struct {
	service Service
}

type Service interface {
	Login(*models.LoginRequest) (*models.LoginResponse, error)
	GetUser(string) (*models.User, error)
	Register(*models.User) (*models.UserRegisterResponse, error)
	SearchUser(string, string) (*[]models.User, error)
	GetPostFeed(*redis.Client, int, int) (*[]models.Post, error)
	DeletePost(string) (bool, error)
	GetPost(string) (*models.Post, error)
	CreatePost(*models.Post) (*models.NewPostResponse, error)
	UpdatePost(*models.Post) error
}

func (srv ServiceHandler) login(w http.ResponseWriter, r *http.Request) {
	loginData := new(models.LoginRequest)

	err := utils.ParseJSON(r, loginData)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}

	token, err := srv.service.Login(loginData)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusAccepted, token)
}

func (srv ServiceHandler) register(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)

	err := utils.ParseJSON(r, user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	if user.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, *r, errors.New("не передан пароль"))
	}
	result, err := srv.service.Register(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusCreated, result)
}

func (srv ServiceHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := srv.service.GetUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, result)
}

func (srv ServiceHandler) searchUser(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	if firstName == "" || lastName == "" {
		utils.WriteError(w, http.StatusBadRequest, *r, errors.New("отсутствуют GET параметры"))
		return
	}

	result, err := srv.service.SearchUser(firstName, lastName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, result)
}

// func getPostFeed(db *sqlx.DB, cache *redis.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var (
// 			offset int = 0
// 			limit  int = 10
// 			err    error
// 		)

// 		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
// 		if err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, *r, err)
// 			return
// 		}

// 		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
// 		if err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, *r, err)
// 			return
// 		}

// 		if offset < 0 || limit < 1 {
// 			utils.WriteError(w, http.StatusBadRequest, *r, errors.New("not valid params"))
// 			return
// 		}

// 		result, err := service.GetPostFeed(db, cache, offset, limit)
// 		if err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, *r, err)
// 			return
// 		}
// 		w.Header().Add("Content-Type", "application/json")
// 		utils.WriteJSON(w, http.StatusOK, result)
// 	}
// }

func (srv ServiceHandler) getPostByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := srv.service.GetPost(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, result)
}

func (srv ServiceHandler) createPost(w http.ResponseWriter, r *http.Request) {
	post := new(models.Post)

	err := utils.ParseJSON(r, post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	result, err := srv.service.CreatePost(post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusCreated, result)

}

func (srv ServiceHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	post := new(models.Post)

	err := utils.ParseJSON(r, post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	err = srv.service.UpdatePost(post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, nil)

}

func (srv ServiceHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := srv.service.DeletePost(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, *r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, result)
}

func (srv ServiceHandler) setFriend(w http.ResponseWriter, r *http.Request) {
	//id := r.PathValue("user_id")
	//log.Println(db)
	w.Write([]byte("запрос в друзья поулчен, метод не реализован"))
	// result, err := service.GetUser(db, id)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }
	// w.Header().Add("Content-Type", "application/json")
	// utils.WriteJSON(w, http.StatusOK, result)
}

func (srv ServiceHandler) deleteFriend(w http.ResponseWriter, r *http.Request) {
	//id := r.PathValue("user_id")
	//log.Println(db)
	w.Write([]byte("запрос на удаление получен, метод не реализован"))
	// result, err := service.GetUser(db, id)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }
	// w.Header().Add("Content-Type", "application/json")
	// utils.WriteJSON(w, http.StatusOK, result)
}
