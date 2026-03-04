package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/betzone/backend/models"
	"github.com/betzone/backend/services"
	"github.com/betzone/backend/utils"
	"github.com/gin-gonic/gin"
)

// HealthHandler returns the API health status
// @Summary Health check
// @Description Check if the API is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API is running",
	})
}

// GetGamesHandler fetches games from Betkraft API
// @Summary Get all games
// @Description Fetch a paginated list of games from Betkraft API
// @Tags Games
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param status query int false "Game status" default(1)
// @Success 200 {object} models.BetkraftGameResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/games [get]
func GetGamesHandler(c *gin.Context, betkraftService *services.BetkraftService) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	statusStr := c.DefaultQuery("status", "1")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)
	status, _ := strconv.Atoi(statusStr)

	// Fetch games from Betkraft API
	response, err := betkraftService.GetGames(page, perPage, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch games",
			Error:   err.Error(),
		})
		return
	}

	// Return the response
	c.JSON(http.StatusOK, response)
}

// GetGameByIDHandler fetches a single game by ID
// @Summary Get a game by ID
// @Description Retrieve detailed information about a specific game
// @Tags Games
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.ApiResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/games/{id} [get]
func GetGameByIDHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement fetching a single game by ID
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Game retrieved successfully",
		Data:    gin.H{"id": id},
	})
}

// CreateBetHandler creates a new bet
// @Summary Create a new bet
// @Description Place a new bet on a casino game (requires authentication)
// @Tags Bets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.Bet true "Bet details"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 402 {object} models.ErrorResponse
// @Router /api/v1/bets [post]
func CreateBetHandler(c *gin.Context) {
	var bet models.Bet
	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	bet.UserID = userID.(string)

	// Validate bet amount
	if bet.Amount <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Bet amount must be greater than 0",
		})
		return
	}

	if bet.GameID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Game ID is required",
		})
		return
	}

	// TODO: Fetch user from database and check balance
	// For now, assume balance check passes
	// In production:
	// - Check user balance >= bet.Amount
	// - If not, return 402 (Payment Required)
	// - Deduct bet.Amount from user balance
	// - Store bet in database with status="pending"
	// - Call Betkraft API to place the bet with proper signature
	// - If Betkraft returns error, rollback balance deduction

	bet.ID = utils.GenerateUUID() // Generate unique bet ID
	bet.Status = "pending"
	bet.CreatedAt = time.Now()
	bet.UpdatedAt = time.Now()

	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "Bet placed successfully",
		Data:    bet,
	})
}

// GetBetsHandler fetches user bets
// @Summary Get user bets
// @Description Retrieve all bets for the authenticated user
// @Tags Bets
// @Produce json
// @Security Bearer
// @Success 200 {object} models.ApiResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/bets [get]
func GetBetsHandler(c *gin.Context, dbService *services.DatabaseService) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized: user not found in context",
			Error:   "user_not_found",
		})
		return
	}

	category := c.Query("category")

	var bets []models.Bet
	query := dbService.DB.Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("status = ?", category)
	}
	if err := query.Order("created_at DESC").Find(&bets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch bets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Bets retrieved successfully",
		Data:    bets,
	})
}

// GetBetByIDHandler fetches a single bet by ID
// @Summary Get a bet by ID
// @Description Retrieve detailed information about a specific bet
// @Tags Bets
// @Produce json
// @Security Bearer
// @Param id path string true "Bet ID"
// @Success 200 {object} models.ApiResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/bets/{id} [get]
func GetBetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement fetching a single bet by ID
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Bet retrieved successfully",
		Data:    gin.H{"id": id},
	})
}

// GetBetStatusHandler queries bet status from the provider
// @Summary Query bet status
// @Description Query the status of bets from the provider's API
// @Tags Bets
// @Produce json
// @Security Bearer
// @Param game_uuid path string true "Game UUID"
// @Param bet_id query string true "Comma-separated bet IDs"
// @Success 200 {object} models.BetStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/bets/status/{game_uuid} [get]
func GetBetStatusHandler(c *gin.Context, betkraftService *services.BetkraftService) {
	gameUUID := c.Param("game_uuid")
	betID := c.Query("bet_id")

	if gameUUID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Game UUID is required",
			Error:   "game_uuid_missing",
		})
		return
	}

	if betID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Bet ID is required",
			Error:   "bet_id_missing",
		})
		return
	}

	// Query bet status from provider
	response, err := betkraftService.QueryBetStatus(gameUUID, betID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to query bet status",
			Error:   err.Error(),
		})
		return
	}

	// Return the response
	c.JSON(http.StatusOK, response)
}

// GetOddsHandler fetches odds for a game
// @Summary Get game odds
// @Description Retrieve all available odds for a specific game
// @Tags Odds
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} models.ApiResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/odds/{gameId} [get]
func GetOddsHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	// TODO: Implement fetching odds for a game
	odds := []models.Odds{}
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Odds retrieved successfully",
		Data: gin.H{
			"game_id": gameID,
			"odds":    odds,
		},
	})
}

// LaunchGameHandler launches a game for a player
// @Summary Launch a game
// @Description Get a game launch URL for a specific player
// @Tags Games
// @Accept json
// @Produce json
// @Param request body models.LaunchGameRequest true "Game launch details"
// @Success 200 {object} models.LaunchGameResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/launch [post]
func LaunchGameHandler(c *gin.Context, betkraftService *services.BetkraftService) {
	var launchReq models.LaunchGameRequest

	if err := c.ShouldBindJSON(&launchReq); err != nil {
		log.Printf("[LaunchGameHandler] Validation error: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get authenticated user (if available)
	userID, exists := c.Get("user_id")
	var providerPlayerID string
	if exists {
		// Fetch user from DB using AuthService
		authService, ok := c.MustGet("authService").(*services.AuthService)
		if ok {
			user, err := authService.GetUserByID(userID.(string))
			if err == nil {
				providerPlayerID = user.ProviderPlayerID
			}
		}
	}
	// Fallback to request player_id if not found
	if providerPlayerID == "" {
		providerPlayerID = launchReq.PlayerID
	}
	launchReq.PlayerID = providerPlayerID

	log.Printf("[LaunchGameHandler] Parsed request - PlayerID: %s, PlayerName: %s, GameUUID: %s, Balance: %f, Demo: %d",
		launchReq.PlayerID, launchReq.PlayerName, launchReq.GameUUID, launchReq.Balance, launchReq.Demo)

	// Launch game with Betkraft API
	response, err := betkraftService.LaunchGame(&launchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to launch game",
			Error:   err.Error(),
		})
		return
	}

	// Return the response
	c.JSON(http.StatusOK, response)
}

// SignupHandler handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.SignupRequest true "Signup details"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/auth/signup [post]
func SignupHandler(c *gin.Context, authService *services.AuthService) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	user, token, err := authService.Signup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Signup failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User:    user,
	})
}

// SigninHandler handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.SigninRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/signin [post]
func SigninHandler(c *gin.Context, authService *services.AuthService) {
	var req models.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	user, token, err := authService.Signin(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Authentication failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	})
}

// GetProfileHandler retrieves the authenticated user's profile
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile information
// @Tags Auth
// @Produce json
// @Security Bearer
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/profile [get]
func GetProfileHandler(c *gin.Context, authService *services.AuthService) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "user_id not found in context",
		})
		return
	}

	user, err := authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		User:    user,
	})
}
