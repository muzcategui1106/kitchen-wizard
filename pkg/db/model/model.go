package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeletedAt sql.NullTime

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey" swaggerignore:"true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
