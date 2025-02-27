package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cams.dev/video_upload_backend/internal/domain/entity"
)

// VideoRepository implements domain.repository.VideoRepository
type VideoRepository struct {
	db *sql.DB
}

// NewVideoRepository creates a new video repository
func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{
		db: db,
	}
}

// Create inserts a new video record
func (r *VideoRepository) Create(ctx context.Context, video *entity.Video) error {
	query := `
		INSERT INTO videos (
			id, title, description, duration, original_url, thumbnail_url,
			status, file_size, mime_type, user_id, resolution_info, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		video.ID,
		video.Title,
		video.Description,
		video.Duration,
		video.OriginalURL,
		video.ThumbnailURL,
		string(video.Status),
		video.FileSize,
		video.MimeType,
		video.UserID,
		video.ResolutionInfo,
		video.CreatedAt,
		video.UpdatedAt,
	)

	return err
}

// GetByID retrieves a video by ID
func (r *VideoRepository) GetByID(ctx context.Context, id string) (*entity.Video, error) {
	query := `
		SELECT
			id, title, description, duration, original_url, thumbnail_url,
			status, file_size, mime_type, user_id, resolution_info, created_at, updated_at
		FROM videos
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var video entity.Video
	var status string

	err := row.Scan(
		&video.ID,
		&video.Title,
		&video.Description,
		&video.Duration,
		&video.OriginalURL,
		&video.ThumbnailURL,
		&status,
		&video.FileSize,
		&video.MimeType,
		&video.UserID,
		&video.ResolutionInfo,
		&video.CreatedAt,
		&video.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video with ID %s not found", id)
		}
		return nil, err
	}

	video.Status = entity.VideoStatus(status)

	return &video, nil
}

// Update updates a video record
func (r *VideoRepository) Update(ctx context.Context, video *entity.Video) error {
	query := `
		UPDATE videos
		SET
			title = $1,
			description = $2,
			duration = $3,
			original_url = $4,
			thumbnail_url = $5,
			status = $6,
			file_size = $7,
			mime_type = $8,
			resolution_info = $9,
			updated_at = $10
		WHERE id = $11
	`

	video.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		video.Title,
		video.Description,
		video.Duration,
		video.OriginalURL,
		video.ThumbnailURL,
		string(video.Status),
		video.FileSize,
		video.MimeType,
		video.ResolutionInfo,
		video.UpdatedAt,
		video.ID,
	)

	return err
}

// List retrieves videos with pagination
func (r *VideoRepository) List(ctx context.Context, userID string, limit, offset int) ([]*entity.Video, error) {
	query := `
		SELECT
			id, title, description, duration, original_url, thumbnail_url,
			status, file_size, mime_type, user_id, resolution_info, created_at, updated_at
		FROM videos
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video

	for rows.Next() {
		var video entity.Video
		var status string

		err := rows.Scan(
			&video.ID,
			&video.Title,
			&video.Description,
			&video.Duration,
			&video.OriginalURL,
			&video.ThumbnailURL,
			&status,
			&video.FileSize,
			&video.MimeType,
			&video.UserID,
			&video.ResolutionInfo,
			&video.CreatedAt,
			&video.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		video.Status = entity.VideoStatus(status)
		videos = append(videos, &video)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
