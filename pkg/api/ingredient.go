package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
)

// api paths
const (
	ingredientCRUDURI  = "ingredient"
	listIngredientsURI = "ingredients"
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

// @BasePath /api/v1
// V1ListIngredients godoc
// @Summary list ingredients
// @Schemes
// @Description list ingredients. By default it list the first 10 ingredients. Maximum number of ingredients to list is 100
// @Tags ingredient
// @Produce json
// @Param	numItems	query	int  true  "number of elements to list. max 100"
// @Success 200 {array}  model.Ingredient
// @Router /ingredients [get]
func (s *ApiServer) V1ListIngredients() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lg := logger.Log
		numItems, err := strconv.Atoi(ctx.DefaultQuery("numItems", "10"))
		if err != nil {
			lg.Sugar().Debug("incorrect query paremeter passed for numItems defaulting to 10")
			numItems = 10
		}

		if numItems > 100 || numItems < 0 {
			numItems = 10
		}

		ingredients, err := s.ingredientRepository.First(ctx, numItems)
		if err != nil {
			lg.Sugar().Errorf("could not create entry in db due to ... %v", err)
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, ErrorResponse{Code: http.StatusNotAcceptable, Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, ingredients)
	}
}
