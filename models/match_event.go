package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchEvent struct {
	ID           string         `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;primaryKey" json:"id"`
	UserID       string         `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;index" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	TotalCourts  int            `gorm:"not null" json:"total_courts"`
	GameType     string         `gorm:"type:varchar(100);not null" json:"game_type"`
	Location     string         `gorm:"type:varchar(255);not null" json:"location"`
	PlayDate     time.Time      `gorm:"not null" json:"play_date"`
	TotalPlayers int            `gorm:"not null" json:"total_players"`
	TeamType     string         `gorm:"type:varchar(50);not null" json:"team_type"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Players      []MatchPlayer  `gorm:"foreignKey:MatchID" json:"players,omitempty"`
	Rounds       []MatchRound   `gorm:"foreignKey:MatchID" json:"rounds,omitempty"`
}

func (MatchEvent) TableName() string {
	return "match_events"
}
