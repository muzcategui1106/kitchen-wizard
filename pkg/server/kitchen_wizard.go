package server

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
)

type kitchenWizardService struct{}

func NewKitchenWizardServer() *kitchenWizardService {
	return &kitchenWizardService{}
}

func (service *kitchenWizardService) Healthz(ctx context.Context, in *empty.Empty) (*HealthzResponse, error) {
	return &HealthzResponse{
		Result: "ok",
	}, nil
}

func (service *kitchenWizardService) mustEmbedUnimplementedKitchenwizardServer() {}
