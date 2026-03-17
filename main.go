package main

import (
	"flag"
	"log"
	"os"
	"scoring_app/config"
	"scoring_app/controllers"
	"scoring_app/database"
	"scoring_app/middleware"
	"scoring_app/repositories"
	"scoring_app/routes"
	"scoring_app/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	resetDB := flag.Bool("reset", false, "Reset database (drop all tables and recreate)")
	seedDB := flag.Bool("seed", false, "Seed initial data")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDatabase()

	db := database.GetDB()

	// Handle database reset flag
	if *resetDB {
		log.Println("⚠️  RESET DATABASE FLAG DETECTED!")
		log.Println("This will delete all data. Press Ctrl+C to cancel or wait 3 seconds...")

		if err := database.ResetDatabase(db); err != nil {
			log.Fatal("Failed to reset database:", err)
		}
		log.Println("✅ Database reset completed")
	} else {
		// Normal migration
		if err := database.AutoMigrate(db); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
	}

	// Handle seed flag
	if *seedDB || *resetDB {
		if err := database.SeedData(db); err != nil {
			log.Fatal("Failed to seed data:", err)
		}
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	accountTypeRepo := repositories.NewAccountTypeRepository(db)
	matchEventRepo := repositories.NewMatchEventRepository(db)
	matchPlayerRepo := repositories.NewMatchPlayerRepository(db)
	matchRoundRepo := repositories.NewMatchRoundRepository(db)
	leaderboardRepo := repositories.NewLeaderboardRepository(db)
	bannerEventRepo := repositories.NewBannerEventRepository(db) // Tambahkan repository untuk BannerEvent

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)
	accountTypeService := services.NewAccountTypeService(accountTypeRepo)
	matchEventService := services.NewMatchEventService(matchEventRepo)
	matchPlayerService := services.NewMatchPlayerService(
		matchPlayerRepo,
		matchEventRepo,
	)
	matchRoundService := services.NewMatchRoundService(
		matchRoundRepo,
		matchEventRepo,
		matchPlayerRepo,
	)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo)
	bannerEventService := services.NewBannerEventService(bannerEventRepo) // Tambahkan service untuk BannerEvent

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	accountTypeController := controllers.NewAccountTypeController(accountTypeService)
	matchEventController := controllers.NewMatchEventController(matchEventService)
	matchPlayerController := controllers.NewMatchPlayerController(matchPlayerService)
	matchRoundController := controllers.NewMatchRoundController(matchRoundService)
	leaderboardController := controllers.NewLeaderboardController(leaderboardService)
	bannerEventsController := controllers.NewBannerEventController(bannerEventService) // Tambahkan controller untuk BannerEvent

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS – allow requests from any origin (web dev server, mobile, etc.)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Setup routes
	routes.SetupRoutes(router, authController, accountTypeController, matchEventController, matchPlayerController, matchRoundController, leaderboardController, bannerEventsController, authMiddleware)

	// Start server (Railway-friendly)
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}
	serverAddr := ":" + port
	log.Printf("🚀 Starting server on %s", serverAddr)
	log.Printf("📊 Database: %s", cfg.Database.DBName)
	log.Printf("🌍 Environment: %s", cfg.AppEnv)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}
