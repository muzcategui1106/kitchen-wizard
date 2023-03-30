package api

import (
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/muzcategui1106/kitchen-wizard/pkg/db/repository"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type ApiServerOption func(s *ApiServer)

type ApiServerConfig struct {
	DBConn *gorm.DB
}

// ApiServer represents a kitchenwizard api server
type ApiServer struct {
	engine               *gin.Engine
	ingredientRepository repository.IIngredient
	address              string
}

func WithSessionManagement() ApiServerOption {
	store := cookie.NewStore([]byte("secret"))
	return func(s *ApiServer) {
		s.engine.Use(sessions.Sessions("kitchenwizard", store))
	}

}

func WithMiddleware(h gin.HandlerFunc) ApiServerOption {
	return func(s *ApiServer) {
		s.engine.Use(h)
	}
}

func WithOIDCAuth(oauth2Config oauth2.Config, idTokenVerifier gooidc.IDTokenVerifier) ApiServerOption {
	authHandler := rest_middleware.NewAuthHandler(oauth2Config, idTokenVerifier)

	return func(s *ApiServer) {
		s.engine.Use(authHandler.AuthenticationInterceptor())
		authHandler.AddAuthHandling(s.engine)
	}
}

func WithTracing() ApiServerOption {
	return func(s *ApiServer) {
		s.engine.Use(ginhttp.Middleware(
			opentracing.GlobalTracer(),
			ginhttp.MWComponentName("api"),
		))
	}
}

func WithCors() ApiServerOption {
	return func(s *ApiServer) {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true // TODO for simplicity this needs to be fixed tho

		s.engine.Use(cors.New(config))

		// // TODO for simplicity this needs to be fixed tho
		// responseCors := func(ctx *gin.Context) {
		// 	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		// 	ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// 	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		// 	ctx.Next()
		// }
		// s.engine.Use(responseCors)
	}
}
