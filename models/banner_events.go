package models

import (
	"time"

	"gorm.io/gorm"
)

type BannerEvent struct {
	ID          string         `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID      string         `gorm:"type:varchar(50);index" json:"user_id" validate:"required"`
	Title       string         `gorm:"size:255;not null" json:"title" validate:"required,min=5"`
	Location    string         `gorm:"size:255" json:"location"`
	Image       string         `gorm:"type:longtext" json:"image"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	StartDate   time.Time      `json:"start_date" validate:"required"`
	EndDate     time.Time      `json:"end_date" validate:"required,gtfield=StartDate"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BannerEvent) TableName() string {
	return "banner_events"
}
