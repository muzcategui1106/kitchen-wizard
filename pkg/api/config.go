package api

import (
	gooidc "github.com/coreos/go-oidc/v3/oidc"
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
