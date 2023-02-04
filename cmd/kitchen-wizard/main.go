package main

import (
	"fmt"
	"log"
	"net"

	"github.com/muzcategui1106/kitchen-wizard/pkg/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8443))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server.RegisterKitchenwizardServer(grpcServer, server.NewKitchenWizardServer())
	grpcServer.Serve(lis)
}
