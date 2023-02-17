package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/muzcategui1106/kitchen-wizard/pkg/api"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
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

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8443))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serverContext := context.Background()

	// start grpc server
	grpcServer, err := api.NewApiGRPCServer(serverContext, lis, api.Config{})
	if err != nil {
		log.Fatalf("could not initialize grpc server: %v", err)
	}
	go grpcServer.Serve(lis)

	// start http server
	httpServer, err := api.NewApiHTTPServer(serverContext)
	if err != nil {
		log.Fatalf("could not initialize http server: %v", err)
	}
	go httpServer.ListenAndServe()

	// run forerver
	stop := make(chan struct{}, 1)
	<-stop

}
