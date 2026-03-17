package validators

type CreateScoreRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Category    string  `json:"category" binding:"required,min=1,max=100"`
	Value       float64 `json:"value" binding:"required,min=0"`
	MaxValue    float64 `json:"max_value" binding:"required,min=0"`
	Description string  `json:"description" binding:"max=500"`
}

type UpdateScoreRequest struct {
	Category    string  `json:"category" binding:"required,min=1,max=100"`
	Value       float64 `json:"value" binding:"required,min=0"`
	MaxValue    float64 `json:"max_value" binding:"required,min=0"`
	Description string  `json:"description" binding:"max=500"`
}
