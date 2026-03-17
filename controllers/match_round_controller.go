package controllers

import (
	"net/http"
	"strconv"

	"scoring_app/models"
	"scoring_app/services"

	"github.com/gin-gonic/gin"
)

type MatchRoundController struct {
	service services.MatchRoundService
}

func NewMatchRoundController(service services.MatchRoundService) *MatchRoundController {
	return &MatchRoundController{service: service}
}

// CREATE
func (c *MatchRoundController) Create(ctx *gin.Context) {
	var req models.MatchRound

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	if err := c.service.Create(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match round", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Match round created successfully", "data": req})
}

// CREATE PAIRING
func (c *MatchRoundController) CreatePairing(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	if matchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	rounds, err := c.service.CreatePairing(matchID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Pairing created successfully",
		"data":    rounds,
	})
}

// UPDATE SCORES
func (c *MatchRoundController) UpdateScores(ctx *gin.Context) {
	roundIDParam := ctx.Param("id")

	roundID, err := strconv.Atoi(roundIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round ID"})
		return
	}

	// Coba bind sebagai array dulu
	var bulkReq []models.UpdateScoreItem
	if err := ctx.ShouldBindJSON(&bulkReq); err == nil {
		if len(bulkReq) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Empty bulk payload"})
			return
		}

		responses, err := c.service.UpdateScoresBulk(uint(roundID), bulkReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Bulk scores updated successfully",
			"data":    responses,
		})
		return
	}

	// Kalau bukan array, coba single object
	var singleReq models.UpdateScoreItem
	if err := ctx.ShouldBindJSON(&singleReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	responses, err := c.service.UpdateScoresBulk(uint(roundID), []models.UpdateScoreItem{singleReq})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(responses) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Score updated successfully",
			"data":    responses[0],
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No updates performed"})
	}
}

// UPDATE
func (c *MatchRoundController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round ID"})
		return
	}

	var req models.MatchRound
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// paksa ID dari URL
	req.ID = uint(id)

	if err := c.service.Update(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update match round", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Match round updated successfully",
	})
}

// GET ALL
func (c *MatchRoundController) GetAll(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rounds, total, err := c.service.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve match rounds",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match rounds retrieved successfully",
		"data": gin.H{
			"match_rounds": rounds,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GET ROUNDS BY MATCH ID
func (c *MatchRoundController) GetByMatchID(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	rounds, err := c.service.GetByMatchID(matchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match rounds", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Match rounds retrieved successfully", "data": rounds})
}

// GET ROUND BY ID
func (c *MatchRoundController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round ID"})
		return
	}

	round, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Match round not found", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Match round retrieved successfully", "data": round})
}

// DELETE
func (c *MatchRoundController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete match round", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Match round deleted successfully"})
}

// DELETE BY MATCH ID
func (c *MatchRoundController) DeleteByMatchID(ctx *gin.Context) {
	matchID := ctx.Param("match_id")

	if err := c.service.DeleteByMatchID(matchID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete match rounds", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Match rounds deleted successfully"})
}
