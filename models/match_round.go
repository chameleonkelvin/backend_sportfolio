package models

import (
	"time"
)

type MatchRound struct {
	ID             uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID        string       `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;index" json:"match_id"`
	Match          *MatchEvent  `gorm:"foreignKey:MatchID;references:ID" json:"match,omitempty"`
	RoundNumber    int          `gorm:"not null" json:"round_number"`
	Court          int          `gorm:"not null" json:"court"`
	TeamAPlayer1ID uint         `gorm:"not null;index" json:"team_a_player_1_id"`
	TeamAPlayer1   *MatchPlayer `gorm:"foreignKey:TeamAPlayer1ID" json:"team_a_player_1,omitempty"`
	TeamAPlayer2ID uint         `gorm:"not null;index" json:"team_a_player_2_id"`
	TeamAPlayer2   *MatchPlayer `gorm:"foreignKey:TeamAPlayer2ID" json:"team_a_player_2,omitempty"`
	ScoreA         int          `gorm:"default:0" json:"score_a"`
	TeamBPlayer1ID uint         `gorm:"not null;index" json:"team_b_player_1_id"`
	TeamBPlayer1   *MatchPlayer `gorm:"foreignKey:TeamBPlayer1ID" json:"team_b_player_1,omitempty"`
	TeamBPlayer2ID uint         `gorm:"not null;index" json:"team_b_player_2_id"`
	TeamBPlayer2   *MatchPlayer `gorm:"foreignKey:TeamBPlayer2ID" json:"team_b_player_2,omitempty"`
	ScoreB         int          `gorm:"default:0" json:"score_b"`
	CreatedAt      time.Time    `json:"created_at"`
}

func (MatchRound) TableName() string {
	return "match_rounds"
}
