package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/muzcategui1106/kitchen-wizard/pkg/api"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/tracing"
)

func main() {
	var logLevel int
	var logTimeFormat string
	var tracingCollectorAddress string

	flag.IntVar(&logLevel, "log-level", -1, "Global log level")
	flag.StringVar(&logTimeFormat, "log-time-format", "2006-01-02T15:04:05Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.StringVar(&tracingCollectorAddress, "otp-collector-address", "http://collector-collector.observability.svc:14268/api/traces", "open tracing collector address")
	flag.Parse()

	mainContext := context.Background()

	if err := logger.Init(logLevel, logTimeFormat); err != nil {
		log.Fatalf("failed to initialize logging: %v", err)
	}

	if err := tracing.InitJaegerTracer(mainContext, tracingCollectorAddress); err != nil {
		logger.Log.Sugar().Warnf("could not setup tracing, erro was %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9443))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// start grpc server
	grpcServer, err := api.NewApiGRPCServer(mainContext, lis, api.Config{})
	if err != nil {
		log.Fatalf("could not initialize grpc server: %v", err)
	}
	go grpcServer.Serve(lis)

	// start http server
	httpServer, err := api.NewApiHTTPServer(mainContext)
	if err != nil {
		log.Fatalf("could not initialize http server: %v", err)
	}
	go httpServer.ListenAndServe()

	// run forerver
	stop := make(chan struct{}, 1)
	<-stop

}
