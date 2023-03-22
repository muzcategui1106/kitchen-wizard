package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
)

// api paths
const (
	ingredientCRUDPath = "ingredient"
)

// @BasePath /api/v1
// V1CreateIngredient godoc
// @Summary creates a new ingredient to be used for recipes
// @Schemes
// @Description creates a new ingredient to be used for recipes
// @Tags ingredient
// @Produce json
// @Param request body model.Ingredient true "Ingredient to be created"
// @Success 200 {object}  IngredientCrudResponse
// @Router /ingredient [post]
func (s *ApiServer) V1CreateIngredient() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lg := logger.Log
		request := &model.Ingredient{}
		if err := ctx.ShouldBind(request); err != nil {
			lg.Sugar().Errorf("could not deserialize into model.Ingredient due to %s", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusNotAcceptable, Error: fmt.Sprintf("failed to deseiralize object due to ... %s", err.Error())})
			return
		}
		err := s.ingredientRepository.Create(ctx, request)
		if err != nil {
			lg.Sugar().Errorf("could not create entry in db due to ... %v", err)
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, ErrorResponse{Code: http.StatusNotAcceptable, Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, &IngredientCrudResponse{IsCreated: true})
	}
}
