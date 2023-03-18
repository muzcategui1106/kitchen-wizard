package api

import (
	"database/sql"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"golang.org/x/oauth2"
)

type ApiServerOption func(s *ApiServer)

type ApiServerConfig struct {
	DBConn *sql.DB
}

// ApiServer represents a kitchenwizard api server
type ApiServer struct {
	engine  *gin.Engine
	dbConn  *sql.DB
	address string
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
