package server

import (
	context "context"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"github.com/muzcategui1106/kitchen-wizard/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

type kitchenWizardService struct{}

func NewKitchenWizardServer(listener net.Listener, cfg Config) (*grpc.Server, error) {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	middleware.AddLogging(logger.Log, opts)
	RegisterKitchenwizardServer(grpcServer, newKitchenWizardServer())
	return grpcServer, nil
}

func newKitchenWizardServer() *kitchenWizardService {
	return &kitchenWizardService{}
}

func (service *kitchenWizardService) Healthz(ctx context.Context, in *empty.Empty) (*HealthzResponse, error) {
	logger := middleware.Extract(ctx)
	logger.Info("health is ok")
	return &HealthzResponse{
		Result: "ok",
	}, nil
}

func (service *kitchenWizardService) mustEmbedUnimplementedKitchenwizardServer() {}
