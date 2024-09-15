package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"tender-service/internal/app"
	"tender-service/internal/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")

	ctx := context.TODO()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal("cannot setup server:", err.Error())
	}
	go func() {
		if err = a.Run(); err != nil {
			log.Fatal("stop server:", err.Error())
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := a.Stop(); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}
}
