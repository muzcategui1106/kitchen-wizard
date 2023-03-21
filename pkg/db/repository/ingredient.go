package repository

import (
	"context"

	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"gorm.io/gorm"
)

// IngredientI represents interface for allowed actions against ingredentients
type IngredientI interface {
	// Saves an ingredient to the DB
	Save(context.Context, *model.Ingredient) error

	// returns the first X ingredients sort by alphabetical order
	First(context.Context, int) ([]Ingredient, error)
}

type Ingredient struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientI {
	return &Ingredient{
		db: db,
	}
}

func (repo *Ingredient) Save(ctx context.Context, i *model.Ingredient) error {
	return repo.db.WithContext(ctx).Save(i).Error
}

func (repo *Ingredient) First(ctx context.Context, limit int) ([]Ingredient, error) {
	var ingredients []Ingredient
	result := repo.db.WithContext(ctx).Limit(limit).Find(&ingredients).Order("name")
	return ingredients, result.Error
}
