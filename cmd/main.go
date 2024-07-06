package main

import (
	"soc_net/internal/app"
)

func main() {
	app := new(app.Application)

	app.Run()
}
