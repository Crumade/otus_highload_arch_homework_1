package app

import (
	"log/slog"
	"net/http"
	"os"
	"social_network/internal/server"
	"time"
)

type Application struct {
	Logger *slog.Logger
	Server *http.Server
}

func (app *Application) Run() {
	app.NewLogger()
	app.NewServer()
	server.RunServer(app.Server)
}

func (app *Application) NewLogger() {
	app.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(app.Logger)
}

func (app *Application) NewServer() {
	app.Server = &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
