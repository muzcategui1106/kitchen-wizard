package api

type HealthzResponse struct {
	OK bool `json:"ok"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
