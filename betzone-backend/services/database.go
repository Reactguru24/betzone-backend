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

	// Auto migrate Bet model
	if err := ds.DB.AutoMigrate(&models.Bet{}); err != nil {
		return fmt.Errorf("migration failed for Bet: %v", err)
	}

	// Auto migrate Transaction model
	if err := ds.DB.AutoMigrate(&models.Transaction{}); err != nil {
		return fmt.Errorf("migration failed for Transaction: %v", err)
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

// User Operations

// GetUserByID retrieves a user by ID
func (ds *DatabaseService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := ds.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}

// UpdateUserBalance updates a user's balance
func (ds *DatabaseService) UpdateUserBalance(userID string, newBalance float64) error {
	if err := ds.DB.Model(&models.User{}).Where("id = ?", userID).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("error updating user balance: %v", err)
	}
	return nil
}

// Bet Operations

// CreateBet creates a new bet record
func (ds *DatabaseService) CreateBet(bet *models.Bet) error {
	if err := ds.DB.Create(bet).Error; err != nil {
		return fmt.Errorf("error creating bet: %v", err)
	}
	return nil
}

// GetBetByID retrieves a bet by ID
func (ds *DatabaseService) GetBetByID(betID string) (*models.Bet, error) {
	var bet models.Bet
	if err := ds.DB.First(&bet, "id = ?", betID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bet not found")
		}
		return nil, fmt.Errorf("error fetching bet: %v", err)
	}
	return &bet, nil
}

// UpdateBetStatus updates a bet's status
func (ds *DatabaseService) UpdateBetStatus(betID string, status string) error {
	if err := ds.DB.Model(&models.Bet{}).Where("id = ?", betID).Update("status", status).Error; err != nil {
		return fmt.Errorf("error updating bet status: %v", err)
	}
	return nil
}

// Transaction Operations

// CreateTransaction creates a new transaction record
func (ds *DatabaseService) CreateTransaction(txn *models.Transaction) error {
	if err := ds.DB.Create(txn).Error; err != nil {
		return fmt.Errorf("error creating transaction: %v", err)
	}
	return nil
}

// GetUserTransactions retrieves all transactions for a user
func (ds *DatabaseService) GetUserTransactions(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := ds.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("error fetching transactions: %v", err)
	}
	return transactions, nil
}

// GetTransactionByBetID retrieves a transaction by bet ID
func (ds *DatabaseService) GetTransactionByBetID(betID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := ds.DB.Where("bet_id = ?", betID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error, just return nil
		}
		return nil, fmt.Errorf("error fetching transaction: %v", err)
	}
	return &transaction, nil
}
