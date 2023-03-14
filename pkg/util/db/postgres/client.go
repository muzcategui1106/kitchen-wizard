package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
)

const (
	defaultDBName = "kitchenwizard"
)

// NewClient creates a new client that connects to postgress DB
func NewClient(dbHost, dbPort, username, password string) (*sql.DB, error) {
	logger := logger.Log

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, username, password, defaultDBName)

	fmt.Println(connStr)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Sugar().Errorf("unable to open sql connection due to %s", err)
		return nil, err
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		logger.Sugar().Errorf("could not ping db due to %s", err)
		return nil, err
	}

	logger.Sugar().Infof("succefuly established connection to %s db", defaultDBName)
	return db, nil
}
