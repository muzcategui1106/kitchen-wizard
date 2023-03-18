package swagger

import (
	docs "github.com/muzcategui1106/kitchen-wizard/pkg/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const UIPrefix = "/swagger-ui/"

// AddSwagger enables swagger for a gin router under /swagger-ui/swagger endopint
func AddSwagger(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
