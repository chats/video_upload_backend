package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cams.dev/video_upload_backend/internal/domain/entity"
)

// UserRepository implements domain.repository.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user record
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, full_name, avatar_url, role, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.AvatarURL,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT
			id, username, email, password_hash, full_name, avatar_url, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanUser(row)
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT
			id, username, email, password_hash, full_name, avatar_url, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRowContext(ctx, query, email)
	return r.scanUser(row)
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT
			id, username, email, password_hash, full_name, avatar_url, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	row := r.db.QueryRowContext(ctx, query, username)
	return r.scanUser(row)
}

// Update updates a user record
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			email = $2,
			password_hash = $3,
			full_name = $4,
			avatar_url = $5,
			role = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.AvatarURL,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, int, error) {
	query := `
		SELECT
			id, username, email, password_hash, full_name, avatar_url, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		var role string
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.AvatarURL,
			&role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		user.Role = entity.UserRole(role)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countErr := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if countErr != nil {
		return users, 0, countErr
	}

	return users, total, nil
}

// CountAll counts the total number of users
func (r *UserRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// CheckEmailExists checks if an email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

// CheckUsernameExists checks if a username already exists
func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

// Helper function to scan a user from a database row
func (r *UserRepository) scanUser(row *sql.Row) (*entity.User, error) {
	user := &entity.User{}
	var role string

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	user.Role = entity.UserRole(role)
	return user, nil
}
