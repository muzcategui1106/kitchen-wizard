package api

import (
	context "context"
	"log"

	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"

	"github.com/gin-gonic/gin"
)

// constants for api group
const (
	apiBasePath = "/api"
)

// constants for versions
const (
	versionV1 = "/v1"
)

// constants for api paths
const (
	healthz = "/healthz"
)

// constants for server configuraiton
const (
	defaultAddress = "0.0.0.0:8443"
)

func NewApiServer(ctx context.Context, cfg ApiServerConfig, opts ...ApiServerOption) (*ApiServer, error) {
	r := gin.Default()
	s := &ApiServer{
		engine: r,
		dbConn: cfg.DBConn,
	}

	for _, opt := range opts {
		opt(s)
	}

	swagger.AddSwagger(r)
	apiV1 := r.Group(apiBasePath + versionV1)
	{
		apiV1.GET(healthz, V1Healthz())
		apiV1.GET(loggedUserPath, s.V1GetLoggedUser())
	}

	return s, nil
}

// Run runs the http server until context stops. it blocks
// TODO find a way to use ctx to stop the server when its channel is closed
func (s *ApiServer) Run(ctx context.Context) {
	if s.address == "" {
		s.address = defaultAddress
	}

	if serverErr := s.engine.Run(s.address); serverErr != nil {
		log.Fatal(serverErr)
	}
}
