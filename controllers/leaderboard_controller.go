package controllers

import (
	"net/http"
	"strconv"

	// "scoring_app/models"
	"scoring_app/services"

	"github.com/gin-gonic/gin"
)

type LeaderboardController struct {
	service services.LeaderboardService
}

func NewLeaderboardController(service services.LeaderboardService) *LeaderboardController {
	return &LeaderboardController{service: service}
}

// GET LEADERBOARD ALL PLAYERS
func (c *LeaderboardController) GetLeaderboardAllPlayers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	leaderboard, total, err := c.service.GetLeaderboardAllPlayers(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Leaderboard retrieved successfully",
		"data":      leaderboard,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GET LEADERBOARD BY MATCH ID
func (c *LeaderboardController) GetLeaderboardByMatchID(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	// 1. Ambil detail event
	event, _ := c.service.GetEventDetail(matchID)

	if matchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	leaderboard, total, err := c.service.GetLeaderboardByMatchID(matchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Leaderboard retrieved successfully",
		"data": gin.H{
			"event_detail": event,
			"leaderboard":  leaderboard,
		},
		"total": total,
	})
}

func (c *LeaderboardController) GetLeaderboardByMatchIDs(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	// 1. Ambil detail event
	event, _ := c.service.GetEventDetail(matchID)

	// 2. Ambil skor
	scores, _ := c.service.GetScoresByMatchID(matchID)

	// 3. Format Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Leaderboard retrieved successfully",
		"data": gin.H{
			"event_detail": event,
			"scores":       scores,
		},
		"total": len(scores),
	})
}
