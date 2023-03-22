package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
)

// constants for api URL
const (
	loggedUserPath = "logged-user"
)

// @BasePath /api/v1
// V1GetLoggedUser godoc
// @Summary get the user from the current sessions
// @Schemes
// @Description get the user from the current sessions by looking into the cookies
// @Tags user
// @Produce json
// @Success 200 {object} User
// @Router /logged-user [get]
// V1GetLoggedUser get the user from the current sessions
func (s *ApiServer) V1GetLoggedUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := rest_middleware.LoggerFromContext(ctx)
		// do not do login if a session ID has been extracted
		session := sessions.Default(ctx)

		// store the id token in the session
		email := session.Get(oidc.EmailKey)
		if email == nil {
			logger.Error("could not extact email from session")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not extract email from session"})
			return
		}

		user := &User{
			Email: email.(string),
		}
		ctx.JSON(http.StatusOK, user)

	}
}
