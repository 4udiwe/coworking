package main

import (
	"os"

	"github.com/4udiwe/coworking/auth-service/internal/app"
)

func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
