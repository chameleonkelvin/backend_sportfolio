package database

import (
	"log"
	"scoring_app/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate in order: parent tables first, then child tables
	err := db.AutoMigrate(
		&models.AccountType{},
		&models.User{},
		&models.MatchEvent{},
		&models.MatchPlayer{},
		&models.MatchRound{},
		&models.BannerEvent{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// DropOldTables - Function to drop old tables (use with caution!)
func DropOldTables(db *gorm.DB) error {
	log.Println("Dropping old tables...")

	// Drop old tables that are no longer used
	if db.Migrator().HasTable("scores") {
		if err := db.Migrator().DropTable("scores"); err != nil {
			return err
		}
		log.Println("Dropped table: scores")
	}

	log.Println("Old tables cleanup completed")
	return nil
}

// ResetDatabase - Drop all tables and recreate (DEVELOPMENT ONLY!)
func ResetDatabase(db *gorm.DB) error {
	log.Println("WARNING: Resetting database - all data will be lost!")

	// Drop all tables
	err := db.Migrator().DropTable(
		&models.MatchRound{},
		&models.MatchPlayer{},
		&models.MatchEvent{},
		&models.User{},
		&models.AccountType{},
	)

	if err != nil {
		return err
	}

	log.Println("All tables dropped")

	// Recreate tables
	return AutoMigrate(db)
}
