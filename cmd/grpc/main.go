package main

import (
	"gin_training/cmd/grpc/configGRPC"
	"gin_training/internal/myGRPC/server"
	pb "gin_training/internal/proto"
	storage "gin_training/internal/storage/postgreSQL"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := configGRPC.SetConfig()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	db, err := storage.NewPDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPsw, cfg.PostgresDB, cfg.PostgresSSL)
	if err != nil {
		log.Fatalf("couldn't connect to db: %v\n", err)
	}

	log.Println(db.Pdb.Ping())

	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, server.NewGRPCStorage(db))

	log.Printf("GRPC server started on port: %v\n", cfg.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
