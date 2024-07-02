package app

import server "soc_net/homework_1/internal/server"

type Application struct{}

func (*Application) Run() {

	server.RunServer()

}
