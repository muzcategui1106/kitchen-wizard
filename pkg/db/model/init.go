package model

import (
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"gorm.io/gorm"
)

// AutoMigrateSchemas ensures all schemas are migrated to the DB
func AutoMigrateSchemas(db *gorm.DB) error {
	schemas := []interface{}{
		&Ingredient{},
	}
	if err := db.AutoMigrate(schemas...); err != nil {
		logger.Log.Sugar().Errorf("could not automigrate schema due to ... %s", err.Error())
		return err
	}

	return nil
}
