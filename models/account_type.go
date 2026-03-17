package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType struct {
	ID          string         `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Users       []User         `gorm:"foreignKey:AccountTypeID" json:"users,omitempty"`
}

func (AccountType) TableName() string {
	return "account_types"
}
