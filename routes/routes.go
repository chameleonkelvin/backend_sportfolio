package routes

import (
	"scoring_app/controllers"
	"scoring_app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	accountTypeController *controllers.AccountTypeController,
	matchEventController *controllers.MatchEventController,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Match Management System API is running",
			"version": "1.0.0",
			"models":  "account_types, users, match_events, match_players, match_rounds",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public - no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Account Types routes (protected - requires authentication)
		accountTypes := v1.Group("/account-types")
		accountTypes.Use(authMiddleware.RequireAuth())
		{
			accountTypes.POST("", accountTypeController.Create)
			accountTypes.GET("", accountTypeController.GetAll)
			accountTypes.GET("/:id", accountTypeController.GetByID)
			accountTypes.PUT("/:id", accountTypeController.Update)
			accountTypes.DELETE("/:id", accountTypeController.Delete)
		}

		// Match Events routes (protected - requires authentication)
		matchEvents := v1.Group("/match-events")
		matchEvents.Use(authMiddleware.RequireAuth())
		{
			matchEvents.POST("", matchEventController.Create)
			matchEvents.GET("", matchEventController.GetAll)
			matchEvents.GET("/my-events", matchEventController.GetMyEvents)
			matchEvents.GET("/:id", matchEventController.GetByID)
			matchEvents.PUT("/:id", matchEventController.Update)
			matchEvents.DELETE("/:id", matchEventController.Delete)
		}

		// TODO: Add more endpoints for users, match_players, match_rounds
	}
}
