package app

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (app *Application) RunServer() {
	s := &http.Server{
		Addr:           ":80",
		Handler:        app.NewRouter( /* cache.Conn*/ ),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	app.logger.Info("HTTP Server starting on port" + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		app.logger.Fatal("HTTTP Server error", zap.Error(err))
	}
}
