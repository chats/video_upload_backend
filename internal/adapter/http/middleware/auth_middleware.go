package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/infrastructure/auth"
	"cams.dev/video_upload_backend/pkg/logger"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	jwtService *auth.JWTService
	logger     *logger.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService *auth.JWTService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// FiberMiddleware authenticates the request for fiber
func (m *AuthMiddleware) FiberMiddleware(c *fiber.Ctx) error {
	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
	}

	// Parse token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
	}

	tokenString := parts[1]

	// Validate token
	claims, err := m.jwtService.ValidateToken(tokenString)
	if err != nil {
		m.logger.Error("Invalid token", logger.Error(err))
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Set user information in context
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("userRole", claims.Role)

	// Continue to the next handler
	return c.Next()
}

// AdminMiddleware ensures the user has admin role
func (m *AuthMiddleware) AdminMiddleware(c *fiber.Ctx) error {
	// Check if user is admin
	userRole := c.Locals("userRole")
	if userRole == nil || userRole != entity.RoleAdmin {
		return fiber.NewError(fiber.StatusForbidden, "Admin access required")
	}

	// Continue to the next handler
	return c.Next()
}
