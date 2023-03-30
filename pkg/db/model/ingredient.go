package model

import (
	"github.com/google/uuid"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/db/format"
	"gorm.io/gorm"
)

type Ingredient struct {
	Model
	Name        string `json:"name" gorm:"uniqueIndex"`
	Description string `json:"description"`
}

func (u *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {

	// create a uuid for the ingredient
	u.ID = uuid.New()

	// ensure the names are consistent upper case at the beginning of each words with no extra spaces between words
	u.Name = format.ComplyAsTitle(u.Name)
	return nil
}
