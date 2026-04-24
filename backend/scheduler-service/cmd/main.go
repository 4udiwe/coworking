package main

import (
	"os"

	"github.com/4udiwe/cowoking/scheduler-service/internal/app"
)

func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
