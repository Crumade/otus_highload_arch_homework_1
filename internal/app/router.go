package app

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) NewRouter( /* cache *redis.Client*/ ) *chi.Mux {

	// rateLimit := RateLimit{
	// 	cache: cache,
	// }

	handler := &ServiceHandler{service: app.service}

	mux := chi.NewRouter()
	mux.Use(LogRequest)
	//mux.Use(rateLimit.Handle)

	mux.Post("/login", handler.login())
	mux.Route("/user", func(r chi.Router) {
		r.Post("/register", handler.register())
		r.Get("/get/{id}", handler.getUserByID())
		r.Get("/search", handler.searchUser())
	})

	mux.Route("/post", func(r chi.Router) {
		//r.Get("/feed", getPostFeed(db, cache))
		r.Get("/get/{id}", handler.getPostByID())
		r.Post("/create", handler.createPost())
		r.Put("/update", handler.updatePost())
		r.Put("/delete/{id}", handler.deletePost())
	})

	mux.Route("/friend", func(r chi.Router) {
		r.Put("/set/{user_id}", handler.setFriend())
		r.Put("/delete/{user_id}", handler.deleteFriend())

	})
	return mux
}
