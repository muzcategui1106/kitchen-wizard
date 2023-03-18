package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
