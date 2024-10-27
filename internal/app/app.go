package app

import (
	"log"
	"log/slog"
	"os"
	"social_network/internal/pkg/storage"
	"social_network/internal/server"
)

type Application struct {
	Logger  *slog.Logger
	CacheDB *storage.CacheDB
	PG      *storage.PostgresDB
}

func (app *Application) Run() {
	app.NewLogger()
	app.InitDB()
	server.RunServer(app.PG, app.CacheDB)
}

func (app *Application) NewLogger() {
	app.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(app.Logger)
}

func (app *Application) InitDB() {
	app.PG = new(storage.PostgresDB)
	err := app.PG.NewConnection()
	if err != nil {
		log.Fatal("DB connection failure")
	}
	//defer app.PG.Conn.Close()

	cache := new(storage.CacheDB)
	err = cache.NewRedisConnection()
	if err != nil {
		log.Fatalf("Redis error: " + err.Error())
	}

	app.PG.MigrateSchema()
	slog.Info("DB connection success")

	err = app.PG.MigrateUsers()
	if err != nil {
		log.Fatal(err)
	}
	err = app.PG.MigratePosts()
	if err != nil {
		log.Fatal(err)
	}

	cache.Warming()
}
