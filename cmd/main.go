package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"gin_training/cmd/config"
	"gin_training/internal/controller"
	"gin_training/internal/storage/postgreSQL"
)

func main() {
	cfg := config.SetConfig()

	db, err := storage.NewPDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPsw, cfg.PostgresDB, cfg.PostgresSSL)
	if err != nil {
		log.Fatalf("couldn't connect to db: %v\n", err)
	}

	log.Println(db.Pdb.Ping())

	controller := controller.NewController(db)

	router := controller.Routes()

	srv := http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		srv.Shutdown(context.Background())
	}()

	log.Printf("HTTP server started on port: %v\n", cfg.HTTPPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
