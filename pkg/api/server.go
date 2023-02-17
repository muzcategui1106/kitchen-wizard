package api

import (
	context "context"
	"net"
	"net/http"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	v1 "github.com/muzcategui1106/kitchen-wizard/pkg/proto/v1"
	grpc_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/grpc/middleware"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type kitchenWizardService struct{}

// NewApiGRPCServer creates a GRPC server for the API
func NewApiGRPCServer(ctx context.Context, listener net.Listener, cfg Config) (*grpc.Server, error) {
	opts := []grpc.ServerOption{}
	lg := logger.Log
	opts = grpc_middleware.AddUnaryInterceptors(opts, lg)
	opts = grpc_middleware.AddStreamInterceptors(opts, lg)
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterApiServer(grpcServer, newKitchenWizardServer())
	return grpcServer, nil
}

func NewApiHTTPServer(ctx context.Context) (*http.Server, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	mux := runtime.NewServeMux()
	rest_middleware.AddLogger(logger.Log, mux)
	err := v1.RegisterApiHandlerFromEndpoint(ctx, mux, "localhost:9443", opts)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr: "0.0.0.0:8443",
		// add handler with middleware
		Handler: rest_middleware.AddRequestID(
			rest_middleware.AddLogger(logger.Log, mux)),
	}

	return srv, nil
}

func newKitchenWizardServer() *kitchenWizardService {
	return &kitchenWizardService{}
}

func (service *kitchenWizardService) Healthz(ctx context.Context, in *empty.Empty) (*v1.HealthzResponse, error) {
	// Add fields the ctxtags of the request which will be added to all extracted loggers.
	grpc_ctxtags.Extract(ctx).Set("custom_tags.string", "something").Set("custom_tags.int", 1337)
	lg := ctxzap.Extract(ctx)
	lg.Debug("health is ok")
	return &v1.HealthzResponse{
		Result: "ok",
	}, nil
}

func (service *kitchenWizardService) mustEmbedUnimplementedKitchenwizardServer() {}
