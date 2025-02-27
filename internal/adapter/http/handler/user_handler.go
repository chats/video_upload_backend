package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/usecase"
	"cams.dev/video_upload_backend/pkg/logger"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	logger      *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username, email, and password are required")
	}

	// Register user
	output, err := h.userUseCase.Register(c.Context(), usecase.RegisterInput{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		FullName: input.FullName,
	})

	if err != nil {
		h.logger.Error("Failed to register user", logger.Error(err))

		// Check for specific errors
		if strings.Contains(err.Error(), "already") {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  output.User,
		"token": output.Token,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if input.Email == "" || input.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	// Login user
	output, err := h.userUseCase.Login(c.Context(), usecase.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		h.logger.Error("Failed to login user", logger.Error(err))
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  output.User,
		"token": output.Token,
	})
}

// GetMe handles requests to get the current user
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}

	// Get user
	user, err := h.userUseCase.GetUserByID(c.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get user", logger.Error(err))
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateProfile handles requests to update a user's profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}

	// Parse request body
	var input struct {
		FullName  string `json:"fullName"`
		AvatarURL string `json:"avatarUrl"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Update profile
	user, err := h.userUseCase.UpdateProfile(c.Context(), userID.(string), usecase.UpdateProfileInput{
		FullName:  input.FullName,
		AvatarURL: input.AvatarURL,
	})

	if err != nil {
		h.logger.Error("Failed to update profile", logger.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// ChangePassword handles requests to change a user's password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}

	// Parse request body
	var input struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if input.OldPassword == "" || input.NewPassword == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Old and new passwords are required")
	}

	// Change password
	err := h.userUseCase.ChangePassword(c.Context(), userID.(string), usecase.ChangePasswordInput{
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	})

	if err != nil {
		h.logger.Error("Failed to change password", logger.Error(err))

		// Check for specific errors
		if strings.Contains(err.Error(), "incorrect") {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to change password")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// ListUsers handles requests to list users (admin only)
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Check if user is admin (this should be handled by an admin middleware)
	userRole := c.Locals("userRole")
	if userRole == nil || userRole != entity.RoleAdmin {
		return fiber.NewError(fiber.StatusForbidden, "Admin access required")
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Get users
	output, err := h.userUseCase.ListUsers(c.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to list users", logger.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list users")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users":      output.Users,
		"totalCount": output.TotalCount,
		"page":       output.Page,
		"limit":      output.Limit,
		"totalPages": output.TotalPages,
	})
}

// DeleteUser handles requests to delete a user (admin only)
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Check if user is admin (this should be handled by an admin middleware)
	userRole := c.Locals("userRole")
	if userRole == nil || userRole != entity.RoleAdmin {
		return fiber.NewError(fiber.StatusForbidden, "Admin access required")
	}

	// Get user ID from URL
	userID := c.Params("id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	// Delete user
	err := h.userUseCase.DeleteUser(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to delete user", logger.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// UpdateUserRole handles requests to update a user's role (admin only)
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	// Check if user is admin (this should be handled by an admin middleware)
	userRole := c.Locals("userRole")
	if userRole == nil || userRole != entity.RoleAdmin {
		return fiber.NewError(fiber.StatusForbidden, "Admin access required")
	}

	// Get user ID from URL
	userID := c.Params("id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	// Parse request body
	var input struct {
		Role string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate role
	if input.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Role is required")
	}

	// Update role
	user, err := h.userUseCase.UpdateUserRole(c.Context(), userID, entity.UserRole(input.Role))
	if err != nil {
		h.logger.Error("Failed to update user role", logger.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user role")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(router fiber.Router, authMiddleware, adminMiddleware fiber.Handler) {
	// Public routes
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected routes
	users := router.Group("/users")
	users.Use(authMiddleware)

	users.Get("/me", h.GetMe)
	users.Put("/profile", h.UpdateProfile)
	users.Put("/password", h.ChangePassword)

	// Admin routes
	users.Use(adminMiddleware)
	users.Get("/", h.ListUsers)
	users.Delete("/:id", h.DeleteUser)
	users.Put("/:id/role", h.UpdateUserRole)
}
