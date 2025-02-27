// cmd/api/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cams.dev/video_upload_backend/internal/adapter/http"
	"cams.dev/video_upload_backend/internal/adapter/http/handler"
	"cams.dev/video_upload_backend/internal/adapter/http/middleware"
	"cams.dev/video_upload_backend/internal/adapter/repository"
	"cams.dev/video_upload_backend/internal/infrastructure/auth"
	"cams.dev/video_upload_backend/internal/infrastructure/database"
	"cams.dev/video_upload_backend/internal/infrastructure/storage"
	"cams.dev/video_upload_backend/internal/infrastructure/transcode"
	"cams.dev/video_upload_backend/internal/usecase"
	"cams.dev/video_upload_backend/pkg/config"
	applogger "cams.dev/video_upload_backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := applogger.NewLogger(applogger.Config{
		Level:  cfg.Server.LogLevel,
		AppEnv: cfg.Server.AppEnv,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Video Transcoding System API...")

	// Initialize database
	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		logger.Fatal("Failed to initialize database: " + err.Error())
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Initialize repositories
	videoRepo := repository.NewVideoRepository(db.DB())
	segmentRepo := repository.NewSegmentRepository(db.DB())
	userRepo := repository.NewUserRepository(db.DB())

	// Initialize storage
	storageRepo := storage.NewS3Storage(
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.Region,
		cfg.Storage.BucketName,
		cfg.Storage.Endpoint,
		cfg.Storage.UseSSL,
	)

	// Initialize auth service
	jwtDuration, err := time.ParseDuration(cfg.Auth.JWTExpiry)
	if err != nil {
		logger.Fatal("Invalid JWT expiry duration: " + err.Error())
	}

	jwtService := auth.NewJWTService(cfg.Auth.JWTSecret, jwtDuration)

	// Initialize transcode service
	transcodeRepo := transcode.NewFFmpegService(
		cfg.Transcode.FFmpegPath,
		cfg.Transcode.FFprobePath,
	)

	// Initialize use cases
	transcodeUseCase := usecase.NewTranscodeUseCase(
		videoRepo,
		segmentRepo,
		storageRepo,
		transcodeRepo,
	)

	videoUseCase := usecase.NewVideoUseCase(
		videoRepo,
		segmentRepo,
		storageRepo,
		transcodeUseCase,
	)

	userUseCase := usecase.NewUserUseCase(
		userRepo,
		jwtService,
	)

	// Initialize HTTP middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, logger)

	// Initialize HTTP handlers
	videoHandler := handler.NewVideoHandler(videoUseCase)
	userHandler := handler.NewUserHandler(userUseCase, logger)

	// Initialize router
	router := http.NewRouter(
		videoHandler,
		userHandler,
		authMiddleware,
		logger,
	)

	// Get fiber app
	app := router.Setup()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Shutting down server...")
		_ = app.Shutdown()
	}()

	// Start server
	logger.Info("Server starting on port " + cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}

	logger.Info("Server stopped gracefully")
}
