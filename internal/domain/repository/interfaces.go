package repository

import (
	"context"
	"time"

	"cams.dev/video_upload_backend/internal/domain/entity"
)

// VideoRepository defines methods for video persistence
type VideoRepository interface {
	Create(ctx context.Context, video *entity.Video) error
	GetByID(ctx context.Context, id string) (*entity.Video, error)
	Update(ctx context.Context, video *entity.Video) error
	List(ctx context.Context, userID string, limit, offset int) ([]*entity.Video, error)
}

// SegmentRepository defines methods for segment persistence
type SegmentRepository interface {
	Create(ctx context.Context, segment *entity.Segment) error
	GetByVideoID(ctx context.Context, videoID string) ([]*entity.Segment, error)
	GetByVideoIDAndResolution(ctx context.Context, videoID string, resolution entity.Resolution) ([]*entity.Segment, error)
}

// StorageRepository defines methods for object storage operations
type StorageRepository interface {
	UploadFile(ctx context.Context, fileName string, data []byte, contentType string) (string, error)
	GetFile(ctx context.Context, fileName string) ([]byte, error)
	GeneratePresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, error)
}

// TranscodeRepository defines methods for video transcoding operations
type TranscodeRepository interface {
	Transcode(ctx context.Context, inputURL string, outputPath string, resolution entity.Resolution, fps int) error
	Segment(ctx context.Context, videoPath string, segmentDuration int, outputPath string) ([]string, error)
	GetVideoInfo(ctx context.Context, videoPath string) (duration float64, width int, height int, err error)
}

// UserRepository defines methods for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, int, error)
	CountAll(ctx context.Context) (int, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}
