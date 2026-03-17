package controllers

import (
	"net/http"
	"scoring_app/services"
	"scoring_app/validators"

	"github.com/gin-gonic/gin"
)

type AccountTypeController struct {
	service services.AccountTypeService
}

func NewAccountTypeController(service services.AccountTypeService) *AccountTypeController {
	return &AccountTypeController{
		service: service,
	}
}

// Create handles creating a new account type
// @route POST /api/v1/account-types
func (ctrl *AccountTypeController) Create(c *gin.Context) {
	var req validators.AccountTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	accountType, err := ctrl.service.Create(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account type created successfully",
		"data":    accountType,
	})
}

// GetAll handles retrieving all account types
// @route GET /api/v1/account-types
func (ctrl *AccountTypeController) GetAll(c *gin.Context) {
	accountTypes, err := ctrl.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account types retrieved successfully",
		"data":    accountTypes,
	})
}

// GetByID handles retrieving a single account type by ID
// @route GET /api/v1/account-types/:id
func (ctrl *AccountTypeController) GetByID(c *gin.Context) {
	id := c.Param("id")

	accountType, err := ctrl.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account type retrieved successfully",
		"data":    accountType,
	})
}

// Update handles updating an existing account type
// @route PUT /api/v1/account-types/:id
func (ctrl *AccountTypeController) Update(c *gin.Context) {
	id := c.Param("id")
	var req validators.AccountTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	accountType, err := ctrl.service.Update(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account type updated successfully",
		"data":    accountType,
	})
}

// Delete handles deleting an account type
// @route DELETE /api/v1/account-types/:id
func (ctrl *AccountTypeController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account type deleted successfully",
	})
}
