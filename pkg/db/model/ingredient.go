package model

import (
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
