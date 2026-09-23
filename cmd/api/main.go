package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-ticket/config"
	_ "movie-ticket/docs"
	"movie-ticket/internal/handlers"
	"movie-ticket/internal/middleware"
	"movie-ticket/internal/models"
	"movie-ticket/internal/repositories"
	"movie-ticket/internal/services"
	"movie-ticket/internal/utils/cache"
	"movie-ticket/internal/utils/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Movie Ticket Booking REST API
// @version 1.0
// @description Production-ready RESTful API for Movie Ticket Reservation System built with Golang, Gin, GORM, and PostgreSQL.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@movieticket.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOi..."
func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. Connect Database
	log.Println("Connecting to PostgreSQL Database...")
	gormLogLevel := logger.Info
	if cfg.AppEnv == "production" {
		gormLogLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Fatal: Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Auto Migrate Entities
	log.Println("Running AutoMigration for schemas...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Cinema{},
		&models.Showtime{},
		&models.Ticket{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Init Redis Cache
	redisClient := cache.InitRedis(cfg)

	// 4. Initialize Dependency Injection
	// Repositories
	userRepo := repositories.NewUserRepository(db)
	movieRepo := repositories.NewMovieRepository(db)
	cinemaRepo := repositories.NewCinemaRepository(db)
	showtimeRepo := repositories.NewShowtimeRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)

	// Services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	movieService := services.NewMovieService(movieRepo, redisClient)
	cinemaService := services.NewCinemaService(cinemaRepo)
	showtimeService := services.NewShowtimeService(showtimeRepo, movieRepo, cinemaRepo)
	ticketService := services.NewTicketService(ticketRepo, showtimeRepo, userRepo, emailService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	movieHandler := handlers.NewMovieHandler(movieService)
	cinemaHandler := handlers.NewCinemaHandler(cinemaService)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	// 5. Setup Gin Router & Middlewares
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware(rate.Limit(30), 60)) // 30 req/s, burst 60

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		dbStatus := "connected"
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "unhealthy: " + err.Error()
		}

		redisStatus := "connected"
		if redisClient != nil && redisClient.Client() != nil {
			if err := redisClient.Client().Ping(c.Request.Context()).Err(); err != nil {
				redisStatus = "disconnected"
			}
		}

		response.Success(c, http.StatusOK, "System is healthy", gin.H{
			"status":    "ok",
			"database":  dbStatus,
			"redis":     redisStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetMe)
		}

		// Movie routes
		movies := v1.Group("/movies")
		{
			movies.GET("", movieHandler.List)
			movies.GET("/:id", movieHandler.GetByID)
			// Admin operations
			movies.POST("", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), movieHandler.Create)
			movies.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), movieHandler.Update)
			movies.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), movieHandler.Delete)
		}

		// Cinema routes
		cinemas := v1.Group("/cinemas")
		{
			cinemas.GET("", cinemaHandler.List)
			cinemas.GET("/:id", cinemaHandler.GetByID)
			// Admin operations
			cinemas.POST("", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), cinemaHandler.Create)
			cinemas.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), cinemaHandler.Update)
			cinemas.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), cinemaHandler.Delete)
		}

		// Showtime routes
		showtimes := v1.Group("/showtimes")
		{
			showtimes.GET("", showtimeHandler.List)
			showtimes.GET("/:id", showtimeHandler.GetByID)
			// Admin operations
			showtimes.POST("", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), showtimeHandler.Create)
			showtimes.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), showtimeHandler.Update)
			showtimes.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin), showtimeHandler.Delete)
		}

		// Ticket routes (Protected)
		tickets := v1.Group("/tickets", middleware.AuthMiddleware())
		{
			tickets.POST("/book", ticketHandler.BookTicket)
			tickets.GET("/my-tickets", ticketHandler.GetMyTickets)
			tickets.GET("/:id", ticketHandler.GetByID)
			tickets.PUT("/:id/cancel", ticketHandler.CancelTicket)

			// Admin operations
			tickets.GET("", middleware.RequireRole(models.RoleAdmin), ticketHandler.GetAllTickets)
			tickets.PATCH("/:id/status", middleware.RequireRole(models.RoleAdmin), ticketHandler.UpdateStatus)
		}

		// User management routes (Protected)
		users := v1.Group("/users", middleware.AuthMiddleware())
		{
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			// Admin operations
			users.GET("", middleware.RequireRole(models.RoleAdmin), userHandler.List)
			users.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), userHandler.Delete)
		}
	}

	// 6. Graceful Shutdown Server Setup
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on port %s (http://localhost:%s)", cfg.Port, cfg.Port)
		log.Printf("📖 Swagger UI available at: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if sqlDB != nil {
		_ = sqlDB.Close()
	}
	fmt.Println("Server exited cleanly.")
}
