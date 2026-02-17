package services

import (
	"fmt"
	"log"

	"github.com/betzone/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseService handles database connection and operations
type DatabaseService struct {
	DB *gorm.DB
}

// NewDatabaseService creates a new database service and establishes a connection
func NewDatabaseService(dsn string) (*DatabaseService, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("✓ Connected to MySQL database successfully")

	return &DatabaseService{DB: db}, nil
}

// Migrate runs database migrations for all models
func (ds *DatabaseService) Migrate() error {
	// Auto migrate User model
	if err := ds.DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("migration failed for User: %v", err)
	}

	log.Println("✓ Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks the database connection
func (ds *DatabaseService) Health() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
