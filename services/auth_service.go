package services

import (
	"errors"
	"fmt"
	"scoring_app/models"
	"scoring_app/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(accountTypeID, username, fullName, email, password string, birthDate *time.Time) (*models.User, error)
	Login(username, password string) (string, *models.User, error)
	GenerateJWT(user *models.User) (string, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: 24 * time.Hour, // Token expires in 24 hours
	}
}

// Register creates a new user account
func (s *authService) Register(accountTypeID, username, fullName, email, password string, birthDate *time.Time) (*models.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &models.User{
		ID:            uuid.New().String(),
		AccountTypeID: accountTypeID,
		Username:      username,
		FullName:      fullName,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		BirthDate:     birthDate,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Load account type relationship
	user, _ = s.userRepo.FindByID(user.ID)

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(username, password string) (string, *models.User, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GenerateJWT creates a JWT token with user data in the payload
func (s *authService) GenerateJWT(user *models.User) (string, error) {
	// Create claims with user data
	claims := jwt.MapClaims{
		"user_id":         user.ID,
		"account_type_id": user.AccountTypeID,
		"username":        user.Username,
		"full_name":       user.FullName,
		"email":           user.Email,
		"exp":             time.Now().Add(s.jwtExpiry).Unix(),
		"iat":             time.Now().Unix(),
	}

	// Add account type name if available
	if user.AccountType != nil {
		claims["account_type_name"] = user.AccountType.Name
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
