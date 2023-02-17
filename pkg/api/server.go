package api

import (
	context "context"
	"net"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	v1 "github.com/muzcategui1106/kitchen-wizard/pkg/proto/v1"
	"github.com/muzcategui1106/kitchen-wizard/pkg/protocol/grpc/middleware"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type kitchenWizardService struct{}

// NewApiGRPCServer creates a GRPC server for the API
func NewApiGRPCServer(ctx context.Context, listener net.Listener, cfg Config) (*grpc.Server, error) {
	opts := []grpc.ServerOption{}
	middleware.AddLogging(logger.Log, opts)
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterApiServer(grpcServer, newKitchenWizardServer())

	return grpcServer, nil
}

func NewApiHTTPServer(ctx context.Context) (*runtime.ServeMux, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	mux := runtime.NewServeMux()
	err := v1.RegisterApiHandlerFromEndpoint(ctx, mux, "localhost:8443", opts)
	return mux, err
}

func newKitchenWizardServer() *kitchenWizardService {
	return &kitchenWizardService{}
}

func (service *kitchenWizardService) Healthz(ctx context.Context, in *empty.Empty) (*v1.HealthzResponse, error) {
	logger := middleware.Extract(ctx)
	logger.Info("health is ok")
	return &v1.HealthzResponse{
		Result: "ok",
	}, nil
}

func (service *kitchenWizardService) mustEmbedUnimplementedKitchenwizardServer() {}
