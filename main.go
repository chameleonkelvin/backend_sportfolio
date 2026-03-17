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
	// Parse flags
	resetDB := flag.Bool("reset", false, "Reset database")
	seedDB := flag.Bool("seed", false, "Seed initial data")
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Init database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDatabase()
	db := database.GetDB()

	// Reset / migrate DB
	if *resetDB {
		log.Println("⚠️ Resetting database...")
		if err := database.ResetDatabase(db); err != nil {
			log.Fatal(err)
		}
		if err := database.SeedData(db); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatal(err)
		}
		if *seedDB {
			if err := database.SeedData(db); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Repos & services
	userRepo := repositories.NewUserRepository(db)
	accountTypeRepo := repositories.NewAccountTypeRepository(db)
	matchEventRepo := repositories.NewMatchEventRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)
	accountTypeService := services.NewAccountTypeService(accountTypeRepo)
	matchEventService := services.NewMatchEventService(matchEventRepo)

	// Controllers
	authController := controllers.NewAuthController(authService)
	accountTypeController := controllers.NewAccountTypeController(accountTypeService)
	matchEventController := controllers.NewMatchEventController(matchEventService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Setup routes
	routes.SetupRoutes(router, authController, accountTypeController, matchEventController, authMiddleware)

	// ✅ Use Railway PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port // fallback, usually 8080
	}
	log.Printf("🚀 Starting server on :%s", port)
	log.Printf("📊 Database: %s", cfg.Database.DBName)
	log.Printf("🌍 Environment: %s", cfg.AppEnv)

	// Bind to :PORT (not localhost)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}