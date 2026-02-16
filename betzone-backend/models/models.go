package models

import "time"

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
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	GameID    string    `json:"game_id"`
	Amount    float64   `json:"amount"`
	OddsValue float64   `json:"odds_value"`
	Status    string    `json:"status"` // pending, won, lost, cancelled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Demo        int     `json:"demo" binding:"required"`
}

type LaunchGameResponse struct {
	StatusCode        int                    `json:"status_code"`
	StatusDescription string                 `json:"status_description"`
	Data              LaunchGameResponseData `json:"data,omitempty"`
}

type LaunchGameResponseData struct {
	URL string `json:"url"`
}
