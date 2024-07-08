package main

import (
	"social_network/internal/app"
)

func main() {
	app := new(app.Application)

	app.Run()
}
