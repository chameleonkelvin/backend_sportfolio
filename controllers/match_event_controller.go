package controllers

import (
	"net/http"
	"scoring_app/services"
	"scoring_app/validators"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MatchEventController struct {
	service services.MatchEventService
}

func NewMatchEventController(service services.MatchEventService) *MatchEventController {
	return &MatchEventController{
		service: service,
	}
}

// Create handles creating a new match event
// @route POST /api/v1/match-events
func (ctrl *MatchEventController) Create(c *gin.Context) {
	var req validators.MatchEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	// Get user info from JWT context
	userID, _ := c.Get("user_id")
	accountTypeID, _ := c.Get("account_type_id")

	// Parse play date
	playDate, err := time.Parse("2006-01-02", req.PlayDate)
	if err != nil {
		// Try parsing with time
		playDate, err = time.Parse("2006-01-02 15:04:05", req.PlayDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid play_date format. Use YYYY-MM-DD or YYYY-MM-DD HH:mm:ss",
			})
			return
		}
	}

	matchEvent, err := ctrl.service.Create(
		userID.(string),
		accountTypeID.(string),
		req.Name,
		req.TotalCourts,
		req.GameType,
		req.Location,
		playDate,
		req.TotalPlayers,
		req.TeamType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Match event created successfully",
		"data":    matchEvent,
	})
}

// GetAll handles retrieving all match events with pagination
// @route GET /api/v1/match-events
func (ctrl *MatchEventController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	matchEvents, total, err := ctrl.service.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match events retrieved successfully",
		"data": gin.H{
			"match_events": matchEvents,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetByID handles retrieving a single match event by ID
// @route GET /api/v1/match-events/:id
func (ctrl *MatchEventController) GetByID(c *gin.Context) {
	id := c.Param("id")

	matchEvent, err := ctrl.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match event retrieved successfully",
		"data":    matchEvent,
	})
}

// GetMyEvents handles retrieving match events created by the authenticated user
// @route GET /api/v1/match-events/my-events
func (ctrl *MatchEventController) GetMyEvents(c *gin.Context) {
	userID, _ := c.Get("user_id")

	matchEvents, err := ctrl.service.GetByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Your match events retrieved successfully",
		"data":    matchEvents,
	})
}

// Update handles updating an existing match event
// @route PUT /api/v1/match-events/:id
func (ctrl *MatchEventController) Update(c *gin.Context) {
	id := c.Param("id")
	var req validators.MatchEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	// Get user info from JWT context
	userID, _ := c.Get("user_id")

	// Parse play date
	playDate, err := time.Parse("2006-01-02", req.PlayDate)
	if err != nil {
		// Try parsing with time
		playDate, err = time.Parse("2006-01-02 15:04:05", req.PlayDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid play_date format. Use YYYY-MM-DD or YYYY-MM-DD HH:mm:ss",
			})
			return
		}
	}

	matchEvent, err := ctrl.service.Update(
		id,
		userID.(string),
		req.Name,
		req.TotalCourts,
		req.GameType,
		req.Location,
		playDate,
		req.TotalPlayers,
		req.TeamType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match event updated successfully",
		"data":    matchEvent,
	})
}

// Delete handles deleting a match event
// @route DELETE /api/v1/match-events/:id
func (ctrl *MatchEventController) Delete(c *gin.Context) {
	id := c.Param("id")

	// Get user info from JWT context
	userID, _ := c.Get("user_id")

	err := ctrl.service.Delete(id, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match event deleted successfully",
	})
}
