package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;primaryKey" json:"id"`
	AccountTypeID string         `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;index" json:"account_type_id"`
	AccountType   *AccountType   `gorm:"foreignKey:AccountTypeID;references:ID" json:"account_type,omitempty"`
	Username      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	FullName      string         `gorm:"type:varchar(255);not null" json:"full_name"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	OTPCode       string         `gorm:"type:varchar(10)" json:"otp_code,omitempty"`
	BirthDate     *time.Time     `gorm:"type:date" json:"birth_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	MatchEvents   []MatchEvent   `gorm:"foreignKey:UserID" json:"match_events,omitempty"`
}

func (User) TableName() string {
	return "users"
}
