package validators

// MatchPlayerRequest represents the match player create/update request payload
type MatchPlayerRequest struct {
	MatchID   string `json:"match_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Gender    *int   `json:"gender" binding:"required,oneof=0 1"`
	TempScore int    `json:"temp_score"`
}

// MatchPlayerResponse represents the match player response
type MatchPlayerResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
