package controllers

import (
	"net/http"
	"scoring_app/models"
	"scoring_app/services"
	"scoring_app/validators"

	"github.com/gin-gonic/gin"
	"strconv"
)

type MatchPlayerController struct {
	service services.MatchPlayerService
}

func NewMatchPlayerController(service services.MatchPlayerService) *MatchPlayerController {
	return &MatchPlayerController{
		service: service,
	}
}

// CREATE
func (c *MatchPlayerController) Create(ctx *gin.Context) {
	var req validators.MatchPlayerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validators.MatchPlayerResponse{
			Success: false,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
		return
	}

	player := models.MatchPlayer{
		MatchID:   req.MatchID,
		Name:      req.Name,
		Gender:    *req.Gender,
		TempScore: req.TempScore,
	}

	if err := c.service.Create(&player); err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to create match player",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match player created successfully",
		Data:    player,
	})
}

// BULK CREATE
func (c *MatchPlayerController) CreateBatch(ctx *gin.Context) {
	var req []validators.MatchPlayerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validators.MatchPlayerResponse{
			Success: false,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
		return
	}

	var players []models.MatchPlayer
	for _, r := range req {
		player := models.MatchPlayer{
			MatchID:   r.MatchID,
			Name:      r.Name,
			Gender:    *r.Gender,
			TempScore: r.TempScore,
		}
		players = append(players, player)
	}

	if err := c.service.CreateBatch(players); err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to create match players",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match players created successfully",
		Data:    players,
	})
}

// GET ALL DATA
func (c *MatchPlayerController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	players, total, err := c.service.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve match players",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match players retrieved successfully",
		"data": gin.H{
			"match_players": players,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GET ALL BY MATCH ID (Query Param)
func (c *MatchPlayerController) GetByMatchID(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	players, err := c.service.FindByMatchID(matchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to retrieve match players",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match players retrieved successfully",
		Data:    players,
	})
}

// GET BY ID
func (c *MatchPlayerController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	player, err := c.service.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, validators.MatchPlayerResponse{
			Success: false,
			Message: "Match player not found",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match player retrieved successfully",
		Data:    player,
	})
}

// UPDATE
func (c *MatchPlayerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req validators.MatchPlayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validators.MatchPlayerResponse{
			Success: false,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
		return
	}

	player, err := c.service.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, validators.MatchPlayerResponse{
			Success: false,
			Message: "Match player not found",
			Data:    err.Error(),
		})
		return
	}

	player.MatchID = req.MatchID
	player.Name = req.Name
	player.Gender = *req.Gender
	player.TempScore = req.TempScore

	if err := c.service.Update(player); err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to update match player",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match player updated successfully",
		Data:    player,
	})
}

// DELETE BY ID
func (c *MatchPlayerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to delete match player",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match player deleted successfully",
	})
}

// ========================
// DELETE BY MATCH ID
// ========================
func (c *MatchPlayerController) DeleteByMatchID(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	if err := c.service.DeleteByMatchID(matchID); err != nil {
		ctx.JSON(http.StatusInternalServerError, validators.MatchPlayerResponse{
			Success: false,
			Message: "Failed to delete match players by match ID",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, validators.MatchPlayerResponse{
		Success: true,
		Message: "Match players deleted successfully by match ID",
	})
}
