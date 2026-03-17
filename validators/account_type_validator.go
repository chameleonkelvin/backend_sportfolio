package validators

// AccountTypeRequest represents the account type create/update request payload
type AccountTypeRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
}

// AccountTypeResponse represents the account type response
type AccountTypeResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
