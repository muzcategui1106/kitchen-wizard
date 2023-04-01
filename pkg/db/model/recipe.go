package model

type Recipe struct {
	Model

	Name         string        `gorm:"unique"`
	Summary      string        `json:"sumamry"`
	Instructions []string      `json:"instructions"`
	ImageURL     string        `json:"imageURL"`
	Ingredients  []*Ingredient `gorm:"many2many:ingredients_recipes;"`
}
