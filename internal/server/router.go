package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB /*, cache *redis.Client*/) *chi.Mux {

	// rateLimit := RateLimit{
	// 	cache: cache,
	// }

	mux := chi.NewRouter()
	mux.Use(LogRequest)
	//mux.Use(rateLimit.Handle)

	mux.Post("/login", login(db))
	mux.Route("/user", func(r chi.Router) {
		r.Post("/register", register(db))
		r.Get("/get/{id}", getUserByID(db))
		r.Get("/search", searchUser(db))
	})

	mux.Route("/post", func(r chi.Router) {
		//r.Get("/feed", getPostFeed(db, cache))
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
