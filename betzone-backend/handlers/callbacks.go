package handlers

import (
	"net/http"

	"github.com/betzone/backend/models"
	"github.com/betzone/backend/services"
	"github.com/betzone/backend/utils"
	"github.com/gin-gonic/gin"
)

// PlayerInfoCallback handles player info requests from Betkraft
// @Summary Player Info Callback
// @Description Respond to player info requests from Betkraft provider
// @Tags Callbacks
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Callback request from Betkraft"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/player_info [post]
func PlayerInfoCallback(c *gin.Context, authService *services.AuthService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Validate signature using tokenKey
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Missing signature",
		})
		return
	}

	// Get the betzone token key from environment
	tokenKey := authService.GetTokenKey() // Will implement this in auth service

	// Verify the signature
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
		})
		return
	}

	// Extract player info
	playerID, _ := request["player_id"].(string)

	// TODO: Fetch player balance from database
	// For now, return mock data
	c.JSON(http.StatusOK, gin.H{
		"player_id":  playerID,
		"balance":    100.00,
		"currency":   "KES",
		"status":     "active",
		"first_name": "John",
		"last_name":  "Doe",
	})
}

// BetCallback handles bet placement notifications from Betkraft
// @Summary Bet Callback
// @Description Handle bet placement notifications from Betkraft provider
// @Tags Callbacks
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Bet callback from Betkraft"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/bet [post]
func BetCallback(c *gin.Context, authService *services.AuthService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
		})
		return
	}

	// Extract bet info
	playerID, _ := request["player_id"].(string)
	betAmount, _ := request["bet_amount"].(float64)
	gameUUID, _ := request["game_uuid"].(string)

	// TODO: Store bet in database, deduct player balance
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"player_id": playerID,
		"game_uuid": gameUUID,
		"amount":    betAmount,
		"message":   "Bet recorded successfully",
	})
}

// WinCallback handles win notifications from Betkraft
// @Summary Win Callback
// @Description Handle player win notifications from Betkraft provider
// @Tags Callbacks
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Win callback from Betkraft"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/win [post]
func WinCallback(c *gin.Context, authService *services.AuthService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
		})
		return
	}

	// Extract win info
	playerID, _ := request["player_id"].(string)
	winAmount, _ := request["win_amount"].(float64)

	// TODO: Update player balance (add wins), update bet status
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"player_id":  playerID,
		"win_amount": winAmount,
		"message":    "Win recorded successfully",
	})
}

// RollbackCallback handles bet rollback notifications from Betkraft
// @Summary Rollback Callback
// @Description Handle bet rollback notifications from Betkraft provider
// @Tags Callbacks
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Rollback callback from Betkraft"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/rollback [post]
func RollbackCallback(c *gin.Context, authService *services.AuthService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
		})
		return
	}

	// Extract rollback info
	playerID, _ := request["player_id"].(string)
	betAmount, _ := request["bet_amount"].(float64)

	// TODO: Refund player balance, mark bet as rolled back
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"player_id":  playerID,
		"bet_amount": betAmount,
		"message":    "Bet rolled back successfully",
	})
}
