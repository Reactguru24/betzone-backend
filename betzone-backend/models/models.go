package models

import "time"

// User represents a user account
type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"uniqueIndex;varchar(20)" json:"phone"`
	Password  string    `gorm:"varchar(255)" json:"-"` // Never return password in JSON
	FirstName string    `gorm:"varchar(100)" json:"first_name"`
	LastName  string    `gorm:"varchar(100)" json:"last_name"`
	Balance   float64   `gorm:"default:0" json:"balance"`
	Currency  string    `gorm:"varchar(10);default:'KES'" json:"currency"`
	Status    string    `gorm:"varchar(50);default:'active'" json:"status"` // active, inactive, suspended
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SignupRequest is the payload for user registration
type SignupRequest struct {
	Phone     string `json:"phone" binding:"required,len=10"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// SigninRequest is the payload for user login
type SigninRequest struct {
	Phone    string `json:"phone" binding:"required,len=10"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after successful sign up or sign in
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Game struct {
	ID          string    `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	GameType    string    `json:"game_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Odds struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	Team      string    `json:"team"`
	OddsValue float64   `json:"odds_value"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bet struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `json:"user_id"`
	GameID    string    `json:"game_id"`
	Amount    float64   `json:"amount"`
	OddsValue float64   `json:"odds_value"`
	Status    string    `gorm:"varchar(50)" json:"status"` // pending, processing, won, lost, rolled_back, cancelled
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Transaction represents a balance transaction (debit or credit)
type Transaction struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	UserID        string    `json:"user_id"`
	BetID         string    `json:"bet_id"`
	Type          string    `gorm:"varchar(50)" json:"type"` // bet_placed, bet_won, bet_lost, rollback
	Amount        float64   `json:"amount"`                  // Positive for credit, negative for debit
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `gorm:"varchar(255)" json:"description"`
	Status        string    `gorm:"varchar(50)" json:"status"` // completed, failed, pending
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// Betkraft API Responses

type BetkraftGameResponse struct {
	StatusCode        int              `json:"status_code"`
	StatusDescription string           `json:"status_description"`
	Data              BetkraftGameData `json:"data"`
}

type BetkraftGameData struct {
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	Total   int            `json:"total"`
	Data    []BetkraftGame `json:"data"`
}

type BetkraftGame struct {
	ID           int         `json:"id"`
	GameID       int         `json:"game_id"`
	GameUUID     string      `json:"game_uuid"`
	GameName     string      `json:"game_name"`
	Thumbnail    string      `json:"thumbnail"`
	MinimumStake string      `json:"minimum_stake"`
	MaximumStake string      `json:"maximum_stake"`
	MaximumWin   string      `json:"maximum_win"`
	Currency     string      `json:"currency"`
	CurrencyList string      `json:"currency_list"`
	Denomination interface{} `json:"denomination"`
	Logo         string      `json:"logo"`
	Status       int         `json:"status"`
	PartnerID    int         `json:"partner_id"`
	Date         string      `json:"date"`
}

// Launch Game Models

type LaunchGameRequest struct {
	PlayerID    string  `json:"player_id" binding:"required"`
	PlayerName  string  `json:"player_name" binding:"required"`
	PlayerToken string  `json:"player_token" binding:"required"`
	GameUUID    string  `json:"game_uuid" binding:"required"`
	Currency    string  `json:"currency" binding:"required"`
	Balance     float64 `json:"balance" binding:"required"`
	Demo        int     `json:"demo"`
}

type LaunchGameResponse struct {
	StatusCode        int                    `json:"status_code"`
	StatusDescription string                 `json:"status_description"`
	Data              LaunchGameResponseData `json:"data,omitempty"`
}

type LaunchGameResponseData struct {
	URL string `json:"url"`
}
