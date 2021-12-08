package main

import (
	"context"
	"gin_training/internal/myGRPC/clientGRPC"
	pb "gin_training/internal/proto"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"

	"gin_training/cmd/config"
	"gin_training/internal/controller"
)

func main() {
	cfg := config.SetConfig()

	conn, err := grpc.Dial(cfg.Grpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %v", err)
	}
	defer conn.Close()

	store := clientGRPC.New(pb.NewBookServiceClient(conn))

	router := controller.NewController(store)

	r := router.Routes()

	srv := http.Server{
		Addr:    cfg.HTTPPort,
		Handler: r,
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
