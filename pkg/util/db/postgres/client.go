package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	defaultDBName = "kitchenwizard"
)

// NewClient creates a new client that connects to postgress DB
func NewClient(dbHost, dbPort, username, password string) (*gorm.DB, error) {
	logger := logger.Log

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, username, password, defaultDBName)

	fmt.Println(connStr)

	// Open a database connection
	dbLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,           // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		logger.Sugar().Errorf("unable to open sql connection due to %s", err)
		return nil, err
	}

	logger.Sugar().Infof("succefuly established connection to %s db", defaultDBName)
	return db, nil
}
