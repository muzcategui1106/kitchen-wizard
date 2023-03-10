package api

import (
	context "context"
	"mime"
	"net"
	"net/http"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	v1 "github.com/muzcategui1106/kitchen-wizard/pkg/proto/v1"
	grpc_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/grpc/middleware"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

func NewApiHTTPServer(ctx context.Context, cfg Config) (*http.Server, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	gwMux := runtime.NewServeMux()
	mux := http.NewServeMux()

	// creating oidc client and verifier
	oauth2Config, verifier, err := oidc.CreateOIDCClient(ctx, cfg.OidcProviderConfig)
	if err != nil {
		return nil, err
	}

	mime.AddExtensionType(".svg", "image/svg+xml")
	// TODO use a proper session key
	authHandler := rest_middleware.NewAuthHandler(oauth2Config, *verifier, []byte("my-dummy-key"))
	mux.Handle(rest_middleware.BaseAuthPathV1, authHandler)
	mux.Handle("/api/", gwMux)
	mux.Handle(swagger.UIPrefix, http.StripPrefix(swagger.UIPrefix, swagger.Handler))

	err = v1.RegisterApiHandlerFromEndpoint(ctx, gwMux, "localhost:9443", opts)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr: "0.0.0.0:8443",
		Handler: rest_middleware.AddRequestID(
			rest_middleware.AddLogger(logger.Log,
				authHandler.AuthenticationInterceptor(mux))),
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

// V1GetLoggedUser get the user from the current sessions
func (service *kitchenWizardService) V1GetLoggedUser(ctx context.Context, in *empty.Empty) (*v1.V1UserInfoResponse, error) {
	lg := ctxzap.Extract(ctx)
	idTokenInterface := ctx.Value(rest_middleware.IDTokenKey)
	if idTokenInterface == nil {
		lg.Sugar().Error("id token not found in context")
		return nil, grpc.Errorf(codes.Unauthenticated, "no authenticated user present")
	}

	idToken, ok := idTokenInterface.(*gooidc.IDToken)
	if !ok {
		lg.Sugar().Error("could not unmarshall id token from context")
		return nil, status.Error(http.StatusInternalServerError, "could not extract id token from context")
	}

	return &v1.V1UserInfoResponse{
		Name:     "",
		Email:    "",
		Username: idToken.Subject,
	}, nil

}

func (service *kitchenWizardService) mustEmbedUnimplementedKitchenwizardServer() {}
