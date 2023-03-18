package api

import (
	context "context"
	"net/http"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"

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

func NewApiHTTPServer(ctx context.Context, cfg Config) (*gin.Engine, error) {
	// creating oidc client and verifier
	oauth2Config, verifier, err := oidc.CreateOIDCClient(ctx, cfg.OidcProviderConfig)
	if err != nil {
		return nil, err
	}
	authHandler := rest_middleware.NewAuthHandler(oauth2Config, *verifier)

	r := gin.Default()
	r.Use(rest_middleware.StructuredLogger(logger.Log))
	r.Use(ginhttp.Middleware(opentracing.GlobalTracer()))
	rest_middleware.AddSessionManagement(r)
	authHandler.AddAuthHandling(r)
	swagger.AddSwagger(r)
	apiV1 := r.Group(apiBasePath + versionV1)
	apiV1.GET(healthz, V1Healthz())
	return r, nil
}

// @BasePath /api/v1

// V1Healthz godoc
// @Summary healthz asserts that the server is running
// @Schemes
// @Description do ping
// @Tags example
// @Produce json
// @Success 200
// @Router /healthz [get]
func V1Healthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &HealthzResponse{
			OK: true,
		})
	}
}

// V1GetLoggedUser get the user from the current sessions
// func (service *kitchenWizardService) V1GetLoggedUser(ctx context.Context, in *empty.Empty) (*v1.V1UserInfoResponse, error) {
// 	lg := ctxzap.Extract(ctx)
// 	fmt.Println("hello")
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		lg.Sugar().Error("could not retrieve metadata from context")
// 		return nil, status.Error(http.StatusInternalServerError, "ould not retrieve metadata from context")
// 	}

// 	// we can assume this is correct as the metadata should include from the OIDC token. perhaps there is an easier
// 	// way of doing this by relating an id token to user info in a database or redis cache. I dont know at the moment
// 	emails := md.Get(oidc.EmailKey)

// 	return &v1.V1UserInfoResponse{
// 		// Name:     profile[0],
// 		Email:    emails[0],
// 		Username: "",
// 	}, nil

// }
