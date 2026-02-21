package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// Validate signature from headers
	signature := c.GetHeader("x-signature-key")
	timestamp := c.GetHeader("x-timestamp")
	if signature == "" || timestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature or timestamp",
		})
		return
	}

	// Get the betzone app key from auth service
	appKey := authService.GetAppKey()
	if appKey == "" {
		log.Printf("Error: App key not configured")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Server configuration error",
		})
		return
	}

	// Generate the dynamic callback token: token = hex(MD5(hex(SHA1(appKey + timestamp))))
	callbackToken := utils.GenerateCallbackToken(appKey, timestamp)

	// Verify the signature by calculating it from request body using the callback token
	calculatedSignature := utils.HashCreate(request, callbackToken)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract player info
	playerID, ok := request["player_id"].(string)
	if !ok || playerID == "" {
		log.Printf("Missing player_id in player_info callback: %v", request["player_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: player_id",
		})
		return
	}

	log.Printf("Processing player_info callback: playerID=%s", playerID)

	// Fetch player balance from database
	user, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error fetching player info for user %s: %v", playerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Player not found",
		})
		return
	}

	log.Printf("Player info found: playerID=%s, Balance=%.2f", user.ID, user.Balance)

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

	// Validate signature from headers
	signature := c.GetHeader("x-signature-key")
	timestamp := c.GetHeader("x-timestamp")
	if signature == "" || timestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature or timestamp",
		})
		return
	}

	// Get the betzone app key from auth service
	appKey := authService.GetAppKey()
	if appKey == "" {
		log.Printf("Error: App key not configured")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Server configuration error",
		})
		return
	}

	// Generate the dynamic callback token: token = hex(MD5(hex(SHA1(appKey + timestamp))))
	callbackToken := utils.GenerateCallbackToken(appKey, timestamp)

	// Verify the signature by calculating it from request body using the callback token
	calculatedSignature := utils.HashCreate(request, callbackToken)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract bet info
	playerID, ok := request["player_id"].(string)
	if !ok || playerID == "" {
		log.Printf("Missing player_id in bet callback: %v", request["player_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: player_id",
		})
		return
	}

	// Handle both 'amount' (from Betkraft) and 'bet_amount' field names, and both float64 and string types
	var betAmount float64
	if amount, ok := request["amount"].(float64); ok {
		betAmount = amount
	} else if amountStr, ok := request["amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(amountStr, 64); err == nil {
			betAmount = parsed
		}
	} else if betAmt, ok := request["bet_amount"].(float64); ok {
		betAmount = betAmt
	} else if betAmtStr, ok := request["bet_amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(betAmtStr, 64); err == nil {
			betAmount = parsed
		}
	}

	if betAmount <= 0 {
		log.Printf("Missing or invalid bet amount in bet callback. Request: %+v", request)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: amount",
		})
		return
	}

	gameUUID, _ := request["game_uuid"].(string)
	betID, ok := request["bet_id"].(string)
	if !ok || betID == "" {
		log.Printf("Missing bet_id in bet callback: %v", request["bet_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: bet_id",
		})
		return
	}

	log.Printf("Processing bet callback: playerID=%s, betID=%s, betAmount=%.2f", playerID, betID, betAmount)

	// Check for duplicate transaction (idempotency)
	existingTxn, err := dbService.GetTransactionByBetID(betID)
	if err == nil && existingTxn != nil && existingTxn.Type == "bet_placed" {
		// Transaction already exists for this bet
		log.Printf("Duplicate bet callback detected for betID=%s, returning 202", betID)
		c.JSON(http.StatusAccepted, gin.H{
			"status_code":        202,
			"status_description": "Duplicate - already processed",
			"data": gin.H{
				"bet_id":       betID,
				"reference_id": betID,
			},
		})
		return
	}

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

	log.Printf("Found user: ID=%s, CurrentBalance=%.2f", user.ID, user.Balance)

	// 5. Check if balance is sufficient
	if user.Balance < betAmount {
		log.Printf("Insufficient balance for user %s. Balance: %.2f, Required: %.2f", playerID, user.Balance, betAmount)
		c.JSON(http.StatusPaymentRequired, gin.H{
			"status_code":        402,
			"status_description": "Insufficient balance",
		})
		return
	}

	// 3. Deduct betAmount from user balance
	newBalance := user.Balance - betAmount
	log.Printf("Deducting bet amount: %.2f from balance %.2f = %.2f", betAmount, user.Balance, newBalance)

	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process bet",
		})
		return
	}

	// Verify the balance was actually updated
	updatedUser, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error verifying balance update: %v", err)
	} else {
		log.Printf("Balance verification: Expected=%.2f, Actual=%.2f", newBalance, updatedUser.Balance)
		if updatedUser.Balance != newBalance {
			log.Printf("WARNING: Balance mismatch after update!")
		}
	}

	// 2. Create bet record in database
	bet := &models.Bet{
		ID:        betID,
		UserID:    playerID,
		GameID:    gameUUID,
		Amount:    betAmount,
		OddsValue: 0, // Odds not provided in callback, can be updated later
		Status:    "processing",
	}

	if err := dbService.CreateBet(bet); err != nil {
		log.Printf("Error creating bet record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to create bet record",
		})
		return
	}

	log.Printf("Bet record created: betID=%s, userID=%s", betID, playerID)

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

	// Validate signature from headers
	signature := c.GetHeader("x-signature-key")
	timestamp := c.GetHeader("x-timestamp")
	if signature == "" || timestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature or timestamp",
		})
		return
	}

	// Get the betzone app key from auth service
	appKey := authService.GetAppKey()
	if appKey == "" {
		log.Printf("Error: App key not configured")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Server configuration error",
		})
		return
	}

	// Generate the dynamic callback token: token = hex(MD5(hex(SHA1(appKey + timestamp))))
	callbackToken := utils.GenerateCallbackToken(appKey, timestamp)

	// Verify the signature by calculating it from request body using the callback token
	calculatedSignature := utils.HashCreate(request, callbackToken)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract and validate player_id
	playerID, ok := request["player_id"].(string)
	if !ok || playerID == "" {
		log.Printf("Missing player_id in win callback: %v", request["player_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: player_id",
		})
		return
	}

	// Extract and validate bet_id early
	betID, ok := request["bet_id"].(string)
	if !ok || betID == "" {
		log.Printf("Missing bet_id in win callback: %v", request["bet_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: bet_id",
		})
		return
	}

	// Handle 'payout_amount' field (primary), with fallback to 'win_amount' or 'amount'
	// Also handle both float64 and string types
	var payoutAmount float64
	if payout, ok := request["payout_amount"].(float64); ok {
		payoutAmount = payout
	} else if payoutStr, ok := request["payout_amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(payoutStr, 64); err == nil {
			payoutAmount = parsed
		}
	} else if winAmt, ok := request["win_amount"].(float64); ok {
		payoutAmount = winAmt
	} else if winAmtStr, ok := request["win_amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(winAmtStr, 64); err == nil {
			payoutAmount = parsed
		}
	} else if amount, ok := request["amount"].(float64); ok {
		payoutAmount = amount
	} else if amountStr, ok := request["amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(amountStr, 64); err == nil {
			payoutAmount = parsed
		}
	}

	// Extract and validate bet status (1=Pending, 2=Won, 3=Lost, 7=Voided)
	var betStatus int64
	if statusVal, ok := request["status"].(float64); ok {
		betStatus = int64(statusVal)
	} else if statusStr, ok := request["status"].(string); ok {
		if parsed, err := strconv.ParseInt(statusStr, 10, 64); err == nil {
			betStatus = parsed
		}
	}

	log.Printf("Processing win callback: playerID=%s, betID=%s, betStatus=%d, payoutAmount=%.2f", playerID, betID, betStatus, payoutAmount)

	// Handle based on bet status
	if betStatus == 3 { // Status 3 = Lost
		log.Printf("Bet lost: playerID=%s, betID=%s, no payout to process", playerID, betID)

		// Verify bet exists
		existingBet, err := dbService.GetBetByID(betID)
		if err != nil || existingBet == nil {
			log.Printf("Bet not found for lost status update: betID=%s, error=%v", betID, err)
			// Still return success to avoid retry loops
			c.JSON(http.StatusOK, gin.H{
				"status_code":        200,
				"status_description": "Success",
				"data": gin.H{
					"bet_id":       betID,
					"reference_id": betID,
				},
			})
			return
		}

		// Check if we already recorded this lost bet
		existingTxn, err := dbService.GetTransactionByBetID(betID)
		if err == nil && existingTxn != nil && existingTxn.Type == "bet_lost" {
			log.Printf("Lost bet already recorded for betID=%s", betID)
			user, _ := dbService.GetUserByID(playerID)
			c.JSON(http.StatusOK, gin.H{
				"status_code":        200,
				"status_description": "Success",
				"data": gin.H{
					"balance":      user.Balance,
					"currency":     "KES",
					"reference_id": betID,
					"date":         time.Now().Format("2006-01-02 15:04:05"),
				},
			})
			return
		}

		// Get user for transaction record
		user, err := dbService.GetUserByID(playerID)
		if err != nil {
			log.Printf("Error finding user %s for lost bet: %v", playerID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code":        500,
				"status_description": "User not found",
			})
			return
		}

		log.Printf("Recording lost bet: playerID=%s, betID=%s, betAmount=%.2f", playerID, betID, existingBet.Amount)

		// Create transaction record for lost bet (no amount change, just audit trail)
		txn := &models.Transaction{
			ID:            utils.GenerateUUID(),
			UserID:        playerID,
			BetID:         betID,
			Type:          "bet_lost",
			Amount:        0, // No balance change for lost bets
			BalanceBefore: user.Balance,
			BalanceAfter:  user.Balance,
			Description:   fmt.Sprintf("Bet lost on game %s (amount: %.2f)", existingBet.GameID, existingBet.Amount),
			Status:        "completed",
		}

		if err := dbService.CreateTransaction(txn); err != nil {
			log.Printf("Error creating lost bet transaction: %v", err)
			// Don't fail the request, but log the error
		}

		// Update bet status to "lost"
		if err := dbService.UpdateBetStatus(betID, "lost"); err != nil {
			log.Printf("Error updating bet status to lost: %v", err)
		}

		log.Printf("Lost bet recorded: playerID=%s, betID=%s, balance=%.2f", playerID, betID, user.Balance)

		// Return success with current balance
		c.JSON(http.StatusOK, gin.H{
			"status_code":        200,
			"status_description": "Success",
			"data": gin.H{
				"balance":      user.Balance,
				"currency":     "KES",
				"reference_id": betID,
				"date":         time.Now().Format("2006-01-02 15:04:05"),
			},
		})
		return
	}

	// For wins (status 2 or other), validate payout_amount is positive
	if payoutAmount <= 0 {
		log.Printf("Missing or invalid payout_amount in win callback. Request: %+v", request)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: payout_amount",
		})
		return
	}

	// Check for duplicate transaction (idempotency)
	existingTxn, err := dbService.GetTransactionByBetID(betID)
	if err == nil && existingTxn != nil && existingTxn.Type == "bet_won" {
		// Transaction already exists for this bet
		log.Printf("Duplicate win callback detected for betID=%s, returning 202", betID)
		c.JSON(http.StatusAccepted, gin.H{
			"status_code":        202,
			"status_description": "Duplicate - already processed",
			"data": gin.H{
				"bet_id":       betID,
				"reference_id": betID,
			},
		})
		return
	}

	// Verify bet exists in system
	existingBet, err := dbService.GetBetByID(betID)
	if err != nil || existingBet == nil {
		log.Printf("Bet not found for win callback: betID=%s, error=%v", betID, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "No bet found for this bet_id",
		})
		return
	}

	log.Printf("Found bet: betID=%s, userID=%s, status=%s", betID, existingBet.UserID, existingBet.Status)

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

	log.Printf("Found user: ID=%s, CurrentBalance=%.2f", user.ID, user.Balance)

	// Calculate new balance upfront for transaction record
	newBalance := user.Balance + payoutAmount

	// Create transaction record IMMEDIATELY to prevent race condition with duplicate callbacks
	txn := &models.Transaction{
		ID:            utils.GenerateUUID(),
		UserID:        playerID,
		BetID:         betID,
		Type:          "bet_won",
		Amount:        payoutAmount, // Positive because it's a credit
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		Description:   fmt.Sprintf("Bet won with payout %.2f", payoutAmount),
		Status:        "completed",
	}

	if err := dbService.CreateTransaction(txn); err != nil {
		log.Printf("Error creating transaction record (may be duplicate): %v", err)
		// If transaction creation fails, it might be a duplicate, return 202
		c.JSON(http.StatusAccepted, gin.H{
			"status_code":        202,
			"status_description": "Duplicate - already processed",
			"data": gin.H{
				"bet_id":       betID,
				"reference_id": betID,
			},
		})
		return
	}

	log.Printf("Transaction record created: txnID=%s", txn.ID)

	// 4. Add payoutAmount to user balance (credit system)
	log.Printf("Crediting payout amount: %.2f to balance %.2f = %.2f", payoutAmount, user.Balance, newBalance)

	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process win",
		})
		return
	}

	// Verify the balance was actually updated
	updatedUser, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error verifying balance update: %v", err)
	} else {
		log.Printf("Balance verification: Expected=%.2f, Actual=%.2f", newBalance, updatedUser.Balance)
		if updatedUser.Balance != newBalance {
			log.Printf("WARNING: Balance mismatch after update!")
		}
	}

	// 3. Update bet status to "won"
	if err := dbService.UpdateBetStatus(betID, "won"); err != nil {
		log.Printf("Error updating bet status: %v", err)
	}

	// 7. Add logging for audit trail
	log.Printf("Win processed: playerID=%s, betID=%s, payoutAmount=%.2f, newBalance=%.2f",
		playerID, betID, payoutAmount, newBalance)

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

	// Validate signature from headers
	signature := c.GetHeader("x-signature-key")
	timestamp := c.GetHeader("x-timestamp")
	if signature == "" || timestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Missing signature or timestamp",
		})
		return
	}

	// Get the betzone app key from auth service
	appKey := authService.GetAppKey()
	if appKey == "" {
		log.Printf("Error: App key not configured")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Server configuration error",
		})
		return
	}

	// Generate the dynamic callback token: token = hex(MD5(hex(SHA1(appKey + timestamp))))
	callbackToken := utils.GenerateCallbackToken(appKey, timestamp)

	// Verify the signature by calculating it from request body using the callback token
	calculatedSignature := utils.HashCreate(request, callbackToken)
	if signature != calculatedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code":        401,
			"status_description": "Invalid signature",
		})
		return
	}

	// Extract rollback info
	playerID, ok := request["player_id"].(string)
	if !ok || playerID == "" {
		log.Printf("Missing player_id in rollback callback: %v", request["player_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: player_id",
		})
		return
	}

	// Handle both 'bet_amount' and 'amount' field names, and both float64 and string types
	var betAmount float64
	if betAmt, ok := request["bet_amount"].(float64); ok {
		betAmount = betAmt
	} else if betAmtStr, ok := request["bet_amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(betAmtStr, 64); err == nil {
			betAmount = parsed
		}
	} else if amount, ok := request["amount"].(float64); ok {
		betAmount = amount
	} else if amountStr, ok := request["amount"].(string); ok {
		if parsed, err := strconv.ParseFloat(amountStr, 64); err == nil {
			betAmount = parsed
		}
	}

	if betAmount <= 0 {
		log.Printf("Missing or invalid bet amount in rollback callback. Request: %+v", request)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: amount",
		})
		return
	}

	betID, ok := request["bet_id"].(string)
	if !ok || betID == "" {
		log.Printf("Missing bet_id in rollback callback: %v", request["bet_id"])
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status_code":        422,
			"status_description": "Missing required parameter: bet_id",
		})
		return
	}

	log.Printf("Processing rollback callback: playerID=%s, betID=%s, betAmount=%.2f", playerID, betID, betAmount)

	// Check for duplicate transaction (idempotency)
	existingTxn, err := dbService.GetTransactionByBetID(betID)
	if err == nil && existingTxn != nil && existingTxn.Type == "rollback" {
		// Transaction already exists for this bet
		log.Printf("Duplicate rollback callback detected for betID=%s, returning 202", betID)
		c.JSON(http.StatusAccepted, gin.H{
			"status_code":        202,
			"status_description": "Duplicate - already processed",
			"data": gin.H{
				"bet_id":       betID,
				"reference_id": betID,
			},
		})
		return
	}

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

	log.Printf("Found user: ID=%s, CurrentBalance=%.2f", user.ID, user.Balance)

	// 2. Refund player balance (add bet amount back)
	newBalance := user.Balance + betAmount
	log.Printf("Refunding bet amount: %.2f to balance %.2f = %.2f", betAmount, user.Balance, newBalance)

	if err := dbService.UpdateUserBalance(playerID, newBalance); err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code":        500,
			"status_description": "Failed to process rollback",
		})
		return
	}

	// Verify the balance was actually updated
	updatedUser, err := dbService.GetUserByID(playerID)
	if err != nil {
		log.Printf("Error verifying balance update: %v", err)
	} else {
		log.Printf("Balance verification: Expected=%.2f, Actual=%.2f", newBalance, updatedUser.Balance)
		if updatedUser.Balance != newBalance {
			log.Printf("WARNING: Balance mismatch after update!")
		}
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
