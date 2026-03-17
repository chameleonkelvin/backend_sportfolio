package models

import (
	"time"
)

type MatchPlayer struct {
	ID        uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID   string      `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;index" json:"match_id"`
	Match     *MatchEvent `gorm:"foreignKey:MatchID;references:ID" json:"match,omitempty"`
	Name      string      `gorm:"type:varchar(255);not null" json:"name"`
	Gender    int         `gorm:"not null;comment:'0 = perempuan, 1 = pria'" json:"gender"`
	TempScore int         `gorm:"default:0" json:"temp_score"`
	CreatedAt time.Time   `json:"created_at"`
}

func (MatchPlayer) TableName() string {
	return "match_players"
}
