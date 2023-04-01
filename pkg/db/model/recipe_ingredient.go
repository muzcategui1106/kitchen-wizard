package model

type IngredientRecipe struct {
	Model

	Ingredient *Ingredient `gorm:"foreignKey:ID" json:"ingredient"`
	Recipe     *Recipe     `gorm:"foreignKey:ID" json:"recipe"`
	Amount     string      `json:"amount"`
}
