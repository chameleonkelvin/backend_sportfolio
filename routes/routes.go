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
	matchPlayerController *controllers.MatchPlayerController,
	matchRoundController *controllers.MatchRoundController,
	leaderboardController *controllers.LeaderboardController,
	bannerEventController *controllers.BannerEventController,
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

		// Match Players routes (protected - requires authentication)
		matchPlayers := v1.Group("/match-players")
		matchPlayers.Use(authMiddleware.RequireAuth())
		{
			matchPlayers.POST("", matchPlayerController.Create)
			matchPlayers.POST("/batch", matchPlayerController.CreateBatch)

			matchPlayers.GET("", matchPlayerController.GetAll)
			matchPlayers.GET("/:id", matchPlayerController.GetByID)
			matchPlayers.GET("/match/:match_id", matchPlayerController.GetByMatchID)

			matchPlayers.PUT("/:id", matchPlayerController.Update)

			matchPlayers.DELETE("/:id", matchPlayerController.Delete)
			matchPlayers.DELETE("/match/:match_id", matchPlayerController.DeleteByMatchID)
		}

		// Match Rounds routes (protected - requires authentication)
		matchRounds := v1.Group("/match-rounds")
		matchRounds.Use(authMiddleware.RequireAuth())
		{
			matchRounds.POST("", matchRoundController.Create)
			matchRounds.POST("/:match_id/pairing", matchRoundController.CreatePairing)
			matchRounds.GET("", matchRoundController.GetAll)
			matchRounds.GET("/:id", matchRoundController.GetByID)
			matchRounds.GET("/match/:match_id", matchRoundController.GetByMatchID)
			matchRounds.PATCH("/:id/score", matchRoundController.UpdateScores)
			matchRounds.PUT("/:id", matchRoundController.Update)
			matchRounds.DELETE("/:id", matchRoundController.Delete)
			matchRounds.DELETE("/match/:match_id", matchRoundController.DeleteByMatchID)
		}

		// Leaderboard routes (protected - requires authentication)
		leaderboard := v1.Group("/leaderboard")
		leaderboard.Use(authMiddleware.RequireAuth())
		{
			leaderboard.GET("/all-players", leaderboardController.GetLeaderboardAllPlayers)
			leaderboard.GET("/match/:match_id", leaderboardController.GetLeaderboardByMatchID)
			leaderboard.GET("/score-match/:match_id", leaderboardController.GetLeaderboardByMatchIDs)
			// leaderboard.GET("/match/:match_id/detail", leaderboardController.GetLeaderboardByMatchIDs)
			// leaderboard.GET("/session/:session", leaderboardController.GetLeaderboardBySession)
			// leaderboard.GET("/overall", leaderboardController.GetOverallLeaderboard)
		}

		// Banner Events routes (protected - requires authentication)
		bannerEvents := v1.Group("/banner-events")
		bannerEvents.Use(authMiddleware.RequireAuth())
		{
			bannerEvents.POST("", bannerEventController.CreateBannerEvent)
			bannerEvents.GET("", bannerEventController.GetAllBannerEvents)
			bannerEvents.GET("/my-events/:userID", bannerEventController.GetBannerEventByUserId)
			bannerEvents.GET("/:id", bannerEventController.GetBannerEventByID)
			bannerEvents.PUT("/:id", bannerEventController.UpdateBannerEvent)
			bannerEvents.DELETE("/:id", bannerEventController.DeleteBannerEvent)
		}

		// TODO: Add more endpoints for users, match_rounds
	}
}
