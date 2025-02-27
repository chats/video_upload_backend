package repository

import (
	"context"
	"database/sql"

	"cams.dev/video_upload_backend/internal/domain/entity"
)

// SegmentRepository implements domain.repository.SegmentRepository
type SegmentRepository struct {
	db *sql.DB
}

// NewSegmentRepository creates a new segment repository
func NewSegmentRepository(db *sql.DB) *SegmentRepository {
	return &SegmentRepository{
		db: db,
	}
}

// Create inserts a new segment record
func (r *SegmentRepository) Create(ctx context.Context, segment *entity.Segment) error {
	query := `
		INSERT INTO segments (
			id, video_id, file_name, url, resolution,
			start_time, duration, segment_index, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		segment.ID,
		segment.VideoID,
		segment.FileName,
		segment.URL,
		string(segment.Resolution),
		segment.StartTime,
		segment.Duration,
		segment.SegmentIndex,
		segment.CreatedAt,
	)

	return err
}

// GetByVideoID retrieves segments for a video
func (r *SegmentRepository) GetByVideoID(ctx context.Context, videoID string) ([]*entity.Segment, error) {
	query := `
		SELECT
			id, video_id, file_name, url, resolution,
			start_time, duration, segment_index, created_at
		FROM segments
		WHERE video_id = $1
		ORDER BY segment_index ASC
	`

	rows, err := r.db.QueryContext(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []*entity.Segment

	for rows.Next() {
		var segment entity.Segment
		var resolution string

		err := rows.Scan(
			&segment.ID,
			&segment.VideoID,
			&segment.FileName,
			&segment.URL,
			&resolution,
			&segment.StartTime,
			&segment.Duration,
			&segment.SegmentIndex,
			&segment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		segment.Resolution = entity.Resolution(resolution)
		segments = append(segments, &segment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

// GetByVideoIDAndResolution retrieves segments for a video with specific resolution
func (r *SegmentRepository) GetByVideoIDAndResolution(
	ctx context.Context,
	videoID string,
	resolution entity.Resolution,
) ([]*entity.Segment, error) {
	query := `
		SELECT
			id, video_id, file_name, url, resolution,
			start_time, duration, segment_index, created_at
		FROM segments
		WHERE video_id = $1 AND resolution = $2
		ORDER BY segment_index ASC
	`

	rows, err := r.db.QueryContext(ctx, query, videoID, string(resolution))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []*entity.Segment

	for rows.Next() {
		var segment entity.Segment
		var resolutionStr string

		err := rows.Scan(
			&segment.ID,
			&segment.VideoID,
			&segment.FileName,
			&segment.URL,
			&resolutionStr,
			&segment.StartTime,
			&segment.Duration,
			&segment.SegmentIndex,
			&segment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		segment.Resolution = entity.Resolution(resolutionStr)
		segments = append(segments, &segment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}
