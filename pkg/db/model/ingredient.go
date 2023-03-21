package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (u *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
