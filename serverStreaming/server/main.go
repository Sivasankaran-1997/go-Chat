package main

import (
	"log"
	"net"
	pb "serverstream/proto"
	"serverstream/server/controller"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen. %v", err)
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	pb.RegisterUserServiceServer(srv, controller.NewUserControllerServer())

	if err := srv.Serve(listen); err != nil {
		log.Fatalf("Failed to serve. %v", err)

	}

}
