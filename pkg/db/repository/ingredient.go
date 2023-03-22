package repository

import (
	"context"

	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"gorm.io/gorm"
)

// IIngredient represents interface for allowed actions against ingredentients
type IIngredient interface {
	// Create an ingredient to the DB
	Create(context.Context, *model.Ingredient) error

	// returns the first X ingredients sort by alphabetical order
	First(context.Context, int) ([]Ingredient, error)
}

type Ingredient struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IIngredient {
	return &Ingredient{
		db: db,
	}
}

func (repo *Ingredient) Create(ctx context.Context, i *model.Ingredient) error {
	return repo.db.WithContext(ctx).Create(i).Error
}

func (repo *Ingredient) First(ctx context.Context, limit int) ([]Ingredient, error) {
	var ingredients []Ingredient
	result := repo.db.WithContext(ctx).Limit(limit).Find(&ingredients).Order("name")
	return ingredients, result.Error
}
