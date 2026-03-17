package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"scoring_app/services"
	"time"

	"github.com/gin-gonic/gin"
)

type BannerEventController struct {
	Service services.BannerEventService
}

func NewBannerEventController(svc services.BannerEventService) *BannerEventController {
	return &BannerEventController{Service: svc}
}

// POST: /banners
func (c *BannerEventController) CreateBannerEvent(ctx *gin.Context) {
	val, exists := ctx.Get("user_id") 
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "message": "Unauthorized: user_id not found in context",
        })
        return
    }

	// Ambil data dari Multipart Form
	userID := fmt.Sprintf("%v", val)
	title := ctx.PostForm("title")
	location := ctx.PostForm("location")
	startDateStr := ctx.PostForm("start_date")
	endDateStr := ctx.PostForm("end_date")
	isActiveStr := ctx.PostForm("is_active")

	description := ctx.PostForm("description")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)
	isActive := isActiveStr == "1" || isActiveStr == "true"

	// Handle File Upload
	imagePath := c.handleFileUpload(ctx)

	input := services.CreateBannerEventInput{
		UserID:    userID,
		Title:     title,
		Location:  location,
		Image:     imagePath,
		Description: description,
		IsActive:  isActive,
		StartDate: startDate,
		EndDate:   endDate,
	}

	event, err := c.Service.CreateBannerEvent(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, event)
}

// PUT: /banners/:id
func (c *BannerEventController) UpdateBannerEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	
	// Gunakan PostForm agar user bisa update foto lewat form-data
	title := ctx.PostForm("title")
	location := ctx.PostForm("location")
	startDateStr := ctx.PostForm("start_date")
	endDateStr := ctx.PostForm("end_date")
	isActiveStr := ctx.PostForm("is_active")

	description := ctx.PostForm("description")

	var startDate, endDate time.Time
	if startDateStr != "" { startDate, _ = time.Parse("2006-01-02", startDateStr) }
	if endDateStr != "" { endDate, _ = time.Parse("2006-01-02", endDateStr) }

	// Handle File Upload (Optional saat update)
	imagePath := c.handleFileUpload(ctx)

	input := services.UpdateBannerEventInput{
		Title:     title,
		Location:  location,
		Image:     imagePath,
		Description: description,
		StartDate: startDate,
		EndDate:   endDate,
	}
	
	// Khusus IsActive (handle pointer agar bisa false)
	if isActiveStr != "" {
		active := isActiveStr == "1" || isActiveStr == "true"
		input.IsActive = &active
	}

	event, err := c.Service.UpdateBannerEvent(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	ctx.JSON(http.StatusOK, event)
}

// Helper untuk upload file agar kode tidak duplikat
func (c *BannerEventController) handleFileUpload(ctx *gin.Context) string {
	file, err := ctx.FormFile("image")
	if err != nil {
		return ""
	}

	uploadDir := "assets/uploads/banner-events"
	os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	targetPath := filepath.Join(uploadDir, filename)

	if err := ctx.SaveUploadedFile(file, targetPath); err != nil {
		return ""
	}

	return targetPath
}

// --- Method lainnya tetap sama ---

func (c *BannerEventController) GetAllBannerEvents(ctx *gin.Context) {
	events, err := c.Service.GetAllBannerEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func (c *BannerEventController) GetBannerEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	event, err := c.Service.GetBannerEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	ctx.JSON(http.StatusOK, event)
}

func (c *BannerEventController) GetBannerEventByUserId(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	events, err := c.Service.GetBannerEventByUserId(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func (c *BannerEventController) DeleteBannerEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.Service.DeleteBannerEvent(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}