package validators

// MatchRoundRequest represents the match round create/update request payload
type MatchRoundRequest struct {
	MatchID        string `json:"match_id" binding:"required,uuid4"`
	RoundNumber    int    `json:"round_number" binding:"required,min=1"`
	Court          int    `json:"court" binding:"required,min=1"`
	TeamAPlayer1ID uint   `json:"team_a_player_1_id" binding:"required"`
	TeamAPlayer2ID uint   `json:"team_a_player_2_id,omitempty"`
	TeamBPlayer1ID uint   `json:"team_b_player_1_id" binding:"required"`
	TeamBPlayer2ID uint   `json:"team_b_player_2_id,omitempty"`
	ScoreA         int    `json:"score_a,omitempty"`
	ScoreB         int    `json:"score_b,omitempty"`
}

// MatchRoundResponse represents the match round response
type MatchRoundResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
