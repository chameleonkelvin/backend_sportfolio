package controllers

import (
	"net/http"
	"scoring_app/services"
	"scoring_app/validators"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
// @route POST /api/v1/auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req validators.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	// Parse birth date if provided
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid birth_date format. Use YYYY-MM-DD (e.g., 2023-10-02)",
			})
			return
		}
		birthDate = &parsedDate
	}

	user, err := ctrl.authService.Register(
		req.AccountTypeID,
		req.Username,
		req.FullName,
		req.Email,
		req.Password,
		birthDate,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generate token for the new user
	token, err := ctrl.authService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":              user.ID,
				"account_type_id": user.AccountTypeID,
				"username":        user.Username,
				"full_name":       user.FullName,
				"email":           user.Email,
				"birth_date":      user.BirthDate,
				"created_at":      user.CreatedAt,
			},
		},
	})
}

// Login handles user authentication
// @route POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req validators.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	token, user, err := ctrl.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":              user.ID,
				"account_type_id": user.AccountTypeID,
				"account_type":    user.AccountType,
				"username":        user.Username,
				"full_name":       user.FullName,
				"email":           user.Email,
				"birth_date":      user.BirthDate,
				"created_at":      user.CreatedAt,
			},
		},
	})
}
