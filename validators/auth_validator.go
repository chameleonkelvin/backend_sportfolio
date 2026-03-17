package validators

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	AccountTypeID string `json:"account_type_id" binding:"required"`
	Username      string `json:"username" binding:"required,min=3,max=100"`
	FullName      string `json:"full_name" binding:"required,min=3,max=255"`
	Email         string `json:"email" binding:"required,email,max=255"`
	Password      string `json:"password" binding:"required,min=6,max=100"`
	BirthDate     string `json:"birth_date"` // Format: YYYY-MM-DD
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
