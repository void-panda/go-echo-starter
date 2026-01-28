package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go-echo-starter/internal/config"
	"go-echo-starter/internal/database"
	"go-echo-starter/internal/handler"
	"go-echo-starter/internal/middleware"
	"go-echo-starter/internal/repository"
	"go-echo-starter/internal/service"
	"go-echo-starter/pkg/jwt"
	"go-echo-starter/pkg/logger"
	"go-echo-starter/pkg/response"
	"go-echo-starter/pkg/validator"

	_ "go-echo-starter/docs"
)

// @title Go Echo Starter API
// @version 1.0
// @description A production-ready Go starter template using Echo framework with layered architecture.

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @schemes http https
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.Log.Level, cfg.IsDevelopment())
	log.Info().Msg("Starting Go Echo Starter application")

	// Initialize database
	db, err := database.NewPostgreSQL(&cfg.Database, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	migrator, err := database.NewMigrator(db.DB, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migrator")
	}
	if err := migrator.Up(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Initialize JWT
	jwtService := jwt.New(&cfg.JWT)

	// Initialize validator
	v := validator.New()

	// Initialize repository
	userRepo := repository.NewUserRepository(db.DB)

	// Initialize service
	userService := service.NewUserService(userRepo, log)
	authService := service.NewAuthService(userRepo, jwtService, log)

	// Initialize handler
	hdlr := handler.NewHandler(userService, authService, v, log)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup middleware
	middleware.Setup(e, log)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return response.Success(c, http.StatusOK, "OK", map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	api := e.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", hdlr.Auth.Register)
			auth.POST("/login", hdlr.Auth.Login)
			auth.GET("/me", hdlr.Auth.GetMe, middleware.JWTAuth(jwtService))
		}

		// User routes (protected)
		users := api.Group("/users", middleware.JWTAuth(jwtService))
		{
			users.POST("", hdlr.User.Create)
			users.GET("", hdlr.User.GetAll)
			users.GET("/:id", hdlr.User.GetByID)
			users.PUT("/:id", hdlr.User.Update)
			users.DELETE("/:id", hdlr.User.Delete)
		}
	}

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		log.Info().Str("address", addr).Msg("Server is starting")
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}
