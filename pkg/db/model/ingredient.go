package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/db/format"
	"gorm.io/gorm"
)

type Ingredient struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (u *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {

	// create a uuid for the ingredient
	u.ID = uuid.New()

	// ensure the names are consistent upper case at the beginning of each words with no extra spaces between words
	u.Name = format.ComplyAsTitle(u.Name)

	similarIngredient := &Ingredient{}
	result := tx.Where(Ingredient{Name: u.Name}).First(similarIngredient)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		logger.Log.Sugar().Errorf("precondition error: could query preexistent ingredients due to ... %s", result.Error.Error())
		return result.Error
	}

	return errors.New("ingredient already exists")

}
