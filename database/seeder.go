package database

import (
	"log"
	"scoring_app/models"
	"time"

	"gorm.io/gorm"
)

// SeedAccountTypes - Seed data untuk account types
func SeedAccountTypes(db *gorm.DB) error {
	accountTypes := []models.AccountType{
		{
			ID:          "admin",
			Name:        "Administrator",
			Description: "Full access to the system",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "organizer",
			Name:        "Event Organizer",
			Description: "Can create and manage match events",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "player",
			Name:        "Player",
			Description: "Can participate in matches",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, accountType := range accountTypes {
		var existing models.AccountType
		result := db.Where("id = ?", accountType.ID).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&accountType).Error; err != nil {
				return err
			}
			log.Printf("Created account type: %s", accountType.Name)
		} else {
			log.Printf("Account type already exists: %s", accountType.Name)
		}
	}

	return nil
}

// SeedData - Seed all initial data
func SeedData(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := SeedAccountTypes(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
