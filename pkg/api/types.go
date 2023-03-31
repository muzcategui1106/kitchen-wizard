package api

import "github.com/muzcategui1106/kitchen-wizard/pkg/db/model"

type HealthzResponse struct {
	OK bool `json:"ok"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUpdateIngredientRequest struct {
	Ingredient *model.Ingredient `json:"ingredient"`
	Image      []byte            `json:"image"`
}

type IngredientCrudResponse struct {
	IsCreated bool `json:"isCreated"`
	IsUpdated bool `json:"isUpdated"`
	IsDeleted bool `json:"isDeleted"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
