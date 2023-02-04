package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"github.com/muzcategui1106/kitchen-wizard/pkg/server"
)

func main() {
	var logLevel int
	var logTimeFormat string

	flag.IntVar(&logLevel, "log-level", 0, "Global log level")
	flag.StringVar(&logTimeFormat, "log-time-format", "2006-01-02T15:04:05Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	if err := logger.Init(logLevel, logTimeFormat); err != nil {
		log.Fatalf("failed to initialize logging: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8443))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, err := server.NewKitchenWizardServer(lis, server.Config{})
	if err != nil {
		log.Fatalf("could not initialize grpc server: %v", err)
	}
	grpcServer.Serve(lis)
}
