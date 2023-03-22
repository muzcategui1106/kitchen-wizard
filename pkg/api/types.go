package api

type HealthzResponse struct {
	OK bool `json:"ok"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
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
