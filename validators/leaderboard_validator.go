package validators

import (
	"scoring_app/models"
)

type LeaderboardRequest struct {
	MatchID string `json:"match_id" binding:"required"`
	Page    int    `json:"page" binding:"required,min=1"`
	Size    int    `json:"size" binding:"required,min=1,max=100"`
}

type SessionLeaderboardRequest struct {
	Session string `json:"session" binding:"required"`
	Page    int    `json:"page" binding:"required,min=1"`
	Size    int    `json:"size" binding:"required,min=1,max=100"`
}

type OverallLeaderboardRequest struct {
	Page int `json:"page" binding:"required,min=1"`
	Size int `json:"size" binding:"required,min=1,max=100"`
}

func (r *LeaderboardRequest) ToServiceRequest() *models.LeaderboardEntry {
	return &models.LeaderboardEntry{
		// UserID:   0, // UserID akan diisi di service layer
		Name:  "", // Username akan diisi di service layer
		Score: 0,  // Score akan dihitung di service layer
		Rank:  0,  // Rank akan dihitung di service layer
	}
}
