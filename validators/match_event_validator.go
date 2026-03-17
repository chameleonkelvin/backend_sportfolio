package validators

// MatchEventRequest represents the match event create/update request payload
type MatchEventRequest struct {
	Name         string `json:"name" binding:"required,min=3,max=255"`
	TotalCourts  int    `json:"total_courts" binding:"required,min=1"`
	GameType     string `json:"game_type" binding:"required,min=3,max=100"`
	Location     string `json:"location" binding:"max=255"`
	PlayDate     string `json:"play_date" binding:"required"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:mm:ss
	TotalPlayers int    `json:"total_players" binding:"required,min=4"`
	TeamType     string `json:"team_type" binding:"required,min=3,max=50"`
}

// MatchEventResponse represents the match event response
type MatchEventResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
