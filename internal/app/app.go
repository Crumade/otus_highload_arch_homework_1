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

	db, err := storage.NewDB()
	if err != nil {
		app.logger.Fatal("DB connection error", zap.Error(err))
	}
	defer db.Close()

	// db, err := storage.NewConnection(storage.PORT_MASTER)
	// if err != nil {
	// 	app.logger.Fatal("DB connection error", zap.Error(err))
	// }
	// defer db.Close()

	cache, err := storage.NewRedisConnection()
	if err != nil {
		app.logger.Fatal("Redis error", zap.Error(err))
	}

	storage.MigrateSchema(db.Master())
	app.logger.Info("MasterDB schema migration success")
	storage.MigrateSchema(db.Slave1())
	app.logger.Info("Slave1DB schema migration success")
	storage.MigrateSchema(db.Slave2())
	app.logger.Info("Slave2DB schema migration success")

	app.MigrateData(db.Master())

	storage.Warming(cache)
	userRepo := storage.NewUserRepo(db)
	postsRepo := storage.NewPostsRepo(db.Master())

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

func (app *Application) MigrateData(db *sqlx.DB) {

	err := storage.MigrateUsers(db)
	if err != nil {
		app.logger.Fatal("User migration error", zap.Error(err))
	}
	err = storage.MigratePosts(db)
	if err != nil {
		app.logger.Fatal("Posts migration error", zap.Error(err))
	}
}
