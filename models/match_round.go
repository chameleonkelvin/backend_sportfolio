package models

import (
	"time"
)

type MatchRound struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Relasi ke MatchEvent
	MatchID string     `gorm:"type:varchar(50);not null;index" json:"match_id"`
	Match   MatchEvent `gorm:"foreignKey:MatchID;references:ID" json:"match"`

	RoundNumber int `gorm:"not null" json:"round_number"`
	Court       int `gorm:"not null" json:"court"`

	// Relasi ke MatchPlayer
	TeamAPlayer1ID uint        `gorm:"not null;index" json:"team_a_player_1_id"`
	TeamAPlayer1   MatchPlayer `gorm:"foreignKey:TeamAPlayer1ID" json:"team_a_player_1"`

	TeamAPlayer2ID uint        `gorm:"not null;index" json:"team_a_player_2_id"`
	TeamAPlayer2   MatchPlayer `gorm:"foreignKey:TeamAPlayer2ID" json:"team_a_player_2"`

	TeamBPlayer1ID uint        `gorm:"not null;index" json:"team_b_player_1_id"`
	TeamBPlayer1   MatchPlayer `gorm:"foreignKey:TeamBPlayer1ID" json:"team_b_player_1"`

	TeamBPlayer2ID uint        `gorm:"not null;index" json:"team_b_player_2_id"`
	TeamBPlayer2   MatchPlayer `gorm:"foreignKey:TeamBPlayer2ID" json:"team_b_player_2"`

	ScoreA    int       `gorm:"default:0" json:"score_a"`
	ScoreB    int       `gorm:"default:0" json:"score_b"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateScoreItem struct {
	Court  int  `json:"court" binding:"required"`
	ScoreA *int `json:"score_a"`
	ScoreB *int `json:"score_b"`
}

type UpdateScoreResponse struct {
	Court       int    `json:"court"`
	ScoreA      int    `json:"score_a"`
	ScoreB      int    `json:"score_b"`
	TotalScore  int    `json:"total_score"`
	MatchStatus string `json:"match_status"`
	CanUpdate   bool   `json:"can_update"`
}

func (MatchRound) TableName() string {
	return "match_rounds"
}
