package app

import (
	server "soc_net/internal/server"
)

type Application struct{}

func (*Application) Run() {

	server.RunServer()

}
