package app

import (
	server "social_network/internal/server"
)

type Application struct{}

func (*Application) Run() {

	server.RunServer()

}
