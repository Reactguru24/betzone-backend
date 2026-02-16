package handlers

import (
	"net/http"
	"strconv"

	"github.com/betzone/backend/models"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// HealthHandler returns the API health status
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API is running",
	})
}

// GetGamesHandler fetches games from Betkraft API
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
func CreateBetHandler(c *gin.Context) {
	var bet models.Bet
	// TODO: Implement creating a bet
	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "Bet created successfully",
		Data:    bet,
	})
}

// GetBetsHandler fetches user bets
func GetBetsHandler(c *gin.Context) {
	bets := []models.Bet{}
	// TODO: Implement fetching user bets
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Bets retrieved successfully",
		Data:    bets,
	})
}

// GetBetByIDHandler fetches a single bet by ID
func GetBetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement fetching a single bet by ID
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Bet retrieved successfully",
		Data:    gin.H{"id": id},
	})
}

// GetOddsHandler fetches odds for a game
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
func LaunchGameHandler(c *gin.Context, betkraftService *services.BetkraftService) {
	var launchReq models.LaunchGameRequest

	if err := c.ShouldBindJSON(&launchReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

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
