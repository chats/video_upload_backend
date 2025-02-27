package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/domain/repository"
	"cams.dev/video_upload_backend/internal/infrastructure/auth"
)

// UserUseCase handles user-related operations
type UserUseCase struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository, jwtService *auth.JWTService) *UserUseCase {
	return &UserUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// RegisterInput represents input data for registering a user
type RegisterInput struct {
	Username string
	Email    string
	Password string
	FullName string
}

// LoginInput represents input data for logging in
type LoginInput struct {
	Email    string
	Password string
}

// UpdateProfileInput represents input data for updating a profile
type UpdateProfileInput struct {
	FullName  string
	AvatarURL string
}

// ChangePasswordInput represents input data for changing a password
type ChangePasswordInput struct {
	OldPassword string
	NewPassword string
}

// UserListOutput represents paginated user list
type UserListOutput struct {
	Users      []*entity.User
	TotalCount int
	Page       int
	Limit      int
	TotalPages int
}

// AuthOutput represents the output of authentication operations
type AuthOutput struct {
	User  *entity.User
	Token string
}

// Register handles user registration
func (uc *UserUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	// Check if email already exists
	exists, err := uc.userRepo.CheckEmailExists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already in use")
	}

	// Check if username already exists
	exists, err = uc.userRepo.CheckUsernameExists(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entity.NewUser(
		uuid.New().String(),
		input.Username,
		input.Email,
		hashedPassword,
		input.FullName,
	)

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		User:  user,
		Token: token,
	}, nil
}

// Login handles user login
func (uc *UserUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := auth.VerifyPassword(input.Password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// GetCurrentUser retrieves the current user from a JWT token
func (uc *UserUseCase) GetCurrentUser(ctx context.Context, tokenString string) (*entity.User, error) {
	// Validate token
	claims, err := uc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	return uc.userRepo.GetByID(ctx, claims.UserID)
}

// UpdateProfile updates a user's profile
func (uc *UserUseCase) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*entity.User, error) {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update profile
	user.UpdateProfile(input.FullName, input.AvatarURL)

	// Save user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (uc *UserUseCase) ChangePassword(ctx context.Context, userID string, input ChangePasswordInput) error {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := auth.VerifyPassword(input.OldPassword, user.PasswordHash); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.SetPassword(hashedPassword)

	// Save user
	return uc.userRepo.Update(ctx, user)
}

// ListUsers retrieves a paginated list of users
func (uc *UserUseCase) ListUsers(ctx context.Context, page, limit int) (*UserListOutput, error) {
	// Ensure valid pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get users
	users, totalCount, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := totalCount / limit
	if totalCount%limit != 0 {
		totalPages++
	}

	return &UserListOutput{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.Delete(ctx, id)
}

// UpdateUserRole updates a user's role
func (uc *UserUseCase) UpdateUserRole(ctx context.Context, id string, role entity.UserRole) (*entity.User, error) {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update role
	user.SetRole(role)
	user.UpdatedAt = time.Now()

	// Save user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
