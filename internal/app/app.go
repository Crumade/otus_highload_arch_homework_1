package app

import (
	"log"
	"social_network/internal/pkg/storage"
	"social_network/internal/service"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	service Service
	logger  *zap.Logger
}

func (app *Application) Run() {
	app.NewZapLogger()

	db, err := storage.NewConnection()
	if err != nil {
		app.logger.Fatal("DB connection error", zap.Error(err))
	}
	defer db.Close()

	cache, err := storage.NewRedisConnection()
	if err != nil {
		app.logger.Fatal("Redis error", zap.Error(err))
	}

	app.InitDB(db)
	storage.Warming(cache)
	userRepo := storage.NewUserRepo(db)
	postsRepo := storage.NewPostsRepo(db)

	app.service = service.NewService(userRepo, postsRepo, app.logger)
	app.RunServer()
}

func (app *Application) NewZapLogger() {
	var err error
	app.logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("logger failure")
	}
}

func (app *Application) InitDB(db *sqlx.DB) {

	storage.MigrateSchema(db)
	app.logger.Info("DB connection success")

	err := storage.MigrateUsers(db)
	if err != nil {
		app.logger.Fatal("User migration error", zap.Error(err))
	}
	err = storage.MigratePosts(db)
	if err != nil {
		app.logger.Fatal("Posts migration error", zap.Error(err))
	}
}
