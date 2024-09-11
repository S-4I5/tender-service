package main

import (
	"context"
	"log"
	"tender-service/internal/app"
	"tender-service/internal/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")

	ctx := context.TODO()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal("cannot setup server:", err.Error())
	}

	if err = a.Run(); err != nil {
		log.Fatal("stop server:", err.Error())
	}
}
