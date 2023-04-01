package repository

import (
	"context"

	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/storage/object"
	"gorm.io/gorm"
)

// constants for object storage
const (
	ingredientsBasePath = "ingredients/"
	defaultImageFormat  = ".jpg"
)

// UploadImageFunc  function that allows a series of bytes to be uploaded to a storage. This is to allow decoupling
// of
type UploadImageFunc func([]byte) error

// IIngredient represents interface for allowed actions against ingredentients
type IIngredient interface {
	// Create an ingredient to the DB
	Create(context.Context, *model.Ingredient) error

	// returns the first X ingredients sort by alphabetical order
	First(context.Context, int) ([]model.Ingredient, error)
}

type Ingredient struct {
	db             *gorm.DB
	iObjectStorage object.Storage
}

func NewIngredientRepository(db *gorm.DB, iObjectStorage object.Storage) IIngredient {
	return &Ingredient{
		db:             db,
		iObjectStorage: iObjectStorage,
	}
}

func (repo *Ingredient) Create(ctx context.Context, i *model.Ingredient) error {
	tx := repo.db.Begin()
	if err := tx.WithContext(ctx).Create(i).Error; err != nil {
		tx.Rollback()
		return err
	}

	res := tx.Commit()
	return res.Error
}

func (repo *Ingredient) First(ctx context.Context, limit int) ([]model.Ingredient, error) {
	var ingredients []model.Ingredient
	result := repo.db.WithContext(ctx).Limit(limit).Find(&ingredients)
	return ingredients, result.Error
}

func imageKey(ingredient *model.Ingredient) string {
	return ingredientsBasePath + ingredient.Name + defaultImageFormat
}
