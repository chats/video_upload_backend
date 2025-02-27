package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"cams.dev/video_upload_backend/internal/adapter/http/handler"
	"cams.dev/video_upload_backend/internal/adapter/http/middleware"
	"cams.dev/video_upload_backend/pkg/logger"
)

// Router sets up the HTTP routes
type Router struct {
	app            *fiber.App
	videoHandler   *handler.VideoHandler
	userHandler    *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
	logger         *logger.Logger
}

// NewRouter creates a new router
func NewRouter(
	videoHandler *handler.VideoHandler,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *logger.Logger,
) *Router {
	app := fiber.New(fiber.Config{
		AppName:      "Video Transcoding System",
		ErrorHandler: customErrorHandler,
	})

	return &Router{
		app:            app,
		videoHandler:   videoHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

// Setup configures the routes
func (r *Router) Setup() *fiber.App {
	// Global middleware
	r.app.Use(recover.New())
	r.app.Use(logger.FiberLogger(r.logger))
	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API v1 group
	apiV1 := r.app.Group("/api/v1")

	// Public routes
	apiV1.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Register user routes
	r.userHandler.RegisterRoutes(
		apiV1,
		r.authMiddleware.FiberMiddleware,
		r.authMiddleware.AdminMiddleware,
	)

	// Video routes (protected)
	videoRoutes := apiV1.Group("/videos")
	videoRoutes.Use(r.authMiddleware.FiberMiddleware)

	videoRoutes.Post("/", r.videoHandler.UploadVideo)
	videoRoutes.Get("/:id", r.videoHandler.GetVideo)
	apiV1.Get("/users/videos", r.videoHandler.GetVideosByUser)

	return r.app
}

// customErrorHandler handles fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}

	// Set Content-Type: application/json
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Return statuscode with error message
	return c.Status(code).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
