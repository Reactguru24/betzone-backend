package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/player_info [post]
func PlayerInfoCallback(c *gin.Context, authService *services.AuthService, dbService *services.DatabaseService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":        400,
			"status_description": "Invalid request",
		})
		return
	}

	// Validate signature using tokenKey
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature",
		})
		return
	}

	// Get the betzone token key from environment
	tokenKey := authService.GetTokenKey()

	// Verify the signature
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract player info
	playerID, _ := request["player_id"].(string)

	// Fetch player balance from database
	user, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error fetching player info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Player not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code":        200,
		"status_description": "Success",
		"data": gin.H{
			"balance":      user.Balance,
			"currency":     "KES",
			"reference_id": playerID,
			"date":         time.Now().Format("2006-01-02 15:04:05"),
		},
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
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 402 {object} models.ErrorResponse
// @Router /api/v1/callbacks/bet [post]
func BetCallback(c *gin.Context, authService *services.AuthService, dbService *services.DatabaseService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":        400,
			"status_description": "Invalid request",
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract bet info
	playerID, _ := request["player_id"].(string)
	betAmount, _ := request["bet_amount"].(float64)
	gameUUID, _ := request["game_uuid"].(string)
	betID, _ := request["bet_id"].(string)

	// 1. Find user by playerID
	user, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error finding user %s: %v", playerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "User not found",
		})
		return
	}

	// 5. Check if balance is sufficient
	if user.Balance < betAmount {
		log.Printf("Insufficient balance for user %s. Balance: %f, Required: %f", playerID, user.Balance, betAmount)
		c.JSON(http.StatusPaymentRequired, gin.H{
			"status_code":        402,
			"status_description": "Insufficient balance",
		})
		return
	}

	// 3. Deduct betAmount from user balance
	newBalance := user.Balance - betAmount
	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process bet",
		})
		return
	}

	// 4. Create transaction record for audit trail
	txn := &models.Transaction{
		ID:            utils.GenerateUUID(),
		UserID:        playerID,
		BetID:         betID,
		Type:          "bet_placed",
		Amount:        -betAmount, // Negative because it's a debit
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		Description:   fmt.Sprintf("Bet placed on game %s", gameUUID),
		Status:        "completed",
	}

	if err := dbService.CreateTransaction(txn); err != nil {
		log.Printf("Error creating transaction: %v", err)
		// Don't fail the request, as the balance deduction was successful
	}

	// 2. Update bet status in database
	if err := dbService.UpdateBetStatus(betID, "processing"); err != nil {
		log.Printf("Error updating bet status: %v", err)
	}

	// 6. Add logging for monitoring and debugging
	log.Printf("Bet processed: playerID=%s, betID=%s, amount=%.2f, gameUUID=%s, newBalance=%.2f",
		playerID, betID, betAmount, gameUUID, newBalance)

	c.JSON(http.StatusOK, gin.H{
		"status_code":        200,
		"status_description": "Success",
		"data": gin.H{
			"balance":      newBalance,
			"currency":     "KES",
			"reference_id": betID,
			"date":         time.Now().Format("2006-01-02 15:04:05"),
		},
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
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/win [post]
func WinCallback(c *gin.Context, authService *services.AuthService, dbService *services.DatabaseService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":        400,
			"status_description": "Invalid request",
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract win info
	playerID, _ := request["player_id"].(string)
	winAmount, _ := request["win_amount"].(float64)
	betID, _ := request["bet_id"].(string)

	// 1. Find user by playerID
	user, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error finding user %s: %v", playerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "User not found",
		})
		return
	}

	// 4. Add winAmount to user balance (credit system)
	newBalance := user.Balance + winAmount
	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process win",
		})
		return
	}

	// 5. Create transaction record showing balance increase
	txn := &models.Transaction{
		ID:            utils.GenerateUUID(),
		UserID:        playerID,
		BetID:         betID,
		Type:          "bet_won",
		Amount:        winAmount, // Positive because it's a credit
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		Description:   fmt.Sprintf("Bet won with amount %.2f", winAmount),
		Status:        "completed",
	}

	if err := dbService.CreateTransaction(txn); err != nil {
		log.Printf("Error creating transaction: %v", err)
		// Don't fail the request, as the balance credit was successful
	}

	// 3. Update bet status to "won"
	if err := dbService.UpdateBetStatus(betID, "won"); err != nil {
		log.Printf("Error updating bet status: %v", err)
	}

	// 7. Add logging for audit trail
	log.Printf("Win processed: playerID=%s, betID=%s, winAmount=%.2f, newBalance=%.2f",
		playerID, betID, winAmount, newBalance)

	// 6. Send notification to player about win (TODO: implement notification system)
	// - Send email, SMS, or push notification
	// - Update real-time UI notification

	c.JSON(http.StatusOK, gin.H{
		"status_code":        200,
		"status_description": "Success",
		"data": gin.H{
			"balance":      newBalance,
			"currency":     "KES",
			"reference_id": betID,
			"date":         time.Now().Format("2006-01-02 15:04:05"),
		},
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
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/callbacks/rollback [post]
func RollbackCallback(c *gin.Context, authService *services.AuthService, dbService *services.DatabaseService) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":        400,
			"status_description": "Invalid request",
		})
		return
	}

	// Validate signature
	signature, ok := request["signature_key"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature",
		})
		return
	}

	tokenKey := authService.GetTokenKey()
	calculatedSignature := utils.HashCreate(request, tokenKey)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract rollback info
	playerID, _ := request["player_id"].(string)
	betAmount, _ := request["bet_amount"].(float64)
	betID, _ := request["bet_id"].(string)

	// 1. Find user by playerID
	user, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error finding user %s: %v", playerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "User not found",
		})
		return
	}

	// 2. Refund player balance (add bet amount back)
	newBalance := user.Balance + betAmount
	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process rollback",
		})
		return
	}

	// Create transaction record for rollback
	txn := &models.Transaction{
		ID:            utils.GenerateUUID(),
		UserID:        playerID,
		BetID:         betID,
		Type:          "rollback",
		Amount:        betAmount, // Positive because it's a refund (credit)
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		Description:   fmt.Sprintf("Bet rolled back, refunded %.2f", betAmount),
		Status:        "completed",
	}

	if err := dbService.CreateTransaction(txn); err != nil {
		log.Printf("Error creating transaction: %v", err)
		// Don't fail the request, as the balance refund was successful
	}

	// 3. Mark bet as rolled back
	if err := dbService.UpdateBetStatus(betID, "rolled_back"); err != nil {
		log.Printf("Error updating bet status: %v", err)
	}

	// Add logging for audit trail
	log.Printf("Rollback processed: playerID=%s, betID=%s, refundAmount=%.2f, newBalance=%.2f",
		playerID, betID, betAmount, newBalance)

	c.JSON(http.StatusOK, gin.H{
		"status_code":        200,
		"status_description": "Success",
		"data": gin.H{
			"balance":      newBalance,
			"currency":     "KES",
			"reference_id": betID,
			"date":         time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}
