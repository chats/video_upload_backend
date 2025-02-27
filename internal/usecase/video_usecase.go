package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/domain/repository"
)

// VideoUseCase handles video-related operations
type VideoUseCase struct {
	videoRepo        repository.VideoRepository
	segmentRepo      repository.SegmentRepository
	storageRepo      repository.StorageRepository
	transcodeUseCase *TranscodeUseCase
}

// NewVideoUseCase creates a new video use case instance
func NewVideoUseCase(
	videoRepo repository.VideoRepository,
	segmentRepo repository.SegmentRepository,
	storageRepo repository.StorageRepository,
	transcodeUseCase *TranscodeUseCase,
) *VideoUseCase {
	return &VideoUseCase{
		videoRepo:        videoRepo,
		segmentRepo:      segmentRepo,
		storageRepo:      storageRepo,
		transcodeUseCase: transcodeUseCase,
	}
}

// VideoUploadInput represents input data for uploading a video
type VideoUploadInput struct {
	Title       string
	Description string
	FileData    []byte
	FileName    string
	FileSize    int64
	MimeType    string
	UserID      string
}

// UploadVideo handles the video upload process
func (uc *VideoUseCase) UploadVideo(ctx context.Context, input VideoUploadInput) (*entity.Video, error) {
	// Generate a unique ID for the video
	videoID := uuid.New().String()

	// Create storage path for the original video
	originalVideoPath := fmt.Sprintf("uploads/%s/original/%s", videoID, filepath.Base(input.FileName))

	// Upload the original video to storage
	uploadURL, err := uc.storageRepo.UploadFile(ctx, originalVideoPath, input.FileData, input.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}

	// Create a new video entity
	video := &entity.Video{
		ID:          videoID,
		Title:       input.Title,
		Description: input.Description,
		OriginalURL: uploadURL,
		Status:      entity.StatusUploaded,
		FileSize:    input.FileSize,
		MimeType:    input.MimeType,
		UserID:      input.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Persist the video information
	if err := uc.videoRepo.Create(ctx, video); err != nil {
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	// Start the transcoding process asynchronously
	go func() {
		// Create a new context since the request context will be cancelled
		bgCtx := context.Background()

		// Update status to processing
		video.Status = entity.StatusProcessing
		if err := uc.videoRepo.Update(bgCtx, video); err != nil {
			// Log error but continue
			fmt.Printf("Failed to update video status: %v\n", err)
		}

		// Process the video (transcode and segment)
		err := uc.transcodeUseCase.ProcessVideo(bgCtx, videoID, uploadURL)

		// Update status based on the result
		if err != nil {
			video.Status = entity.StatusFailed
			fmt.Printf("Video processing failed: %v\n", err)
		} else {
			video.Status = entity.StatusComplete
		}

		// Update the video status
		if updateErr := uc.videoRepo.Update(bgCtx, video); updateErr != nil {
			fmt.Printf("Failed to update video status: %v\n", updateErr)
		}
	}()

	return video, nil
}

// GetVideoByID retrieves a video by its ID
func (uc *VideoUseCase) GetVideoByID(ctx context.Context, id string) (*entity.Video, error) {
	return uc.videoRepo.GetByID(ctx, id)
}

// GetVideoSegments retrieves all segments for a video
func (uc *VideoUseCase) GetVideoSegments(ctx context.Context, videoID string) ([]*entity.Segment, error) {
	return uc.segmentRepo.GetByVideoID(ctx, videoID)
}

// GetVideoSegmentsByResolution retrieves segments for a video filtered by resolution
func (uc *VideoUseCase) GetVideoSegmentsByResolution(
	ctx context.Context,
	videoID string,
	resolution entity.Resolution,
) ([]*entity.Segment, error) {
	return uc.segmentRepo.GetByVideoIDAndResolution(ctx, videoID, resolution)
}

// ListVideos retrieves a paginated list of videos
func (uc *VideoUseCase) ListVideos(ctx context.Context, userID string, limit, offset int) ([]*entity.Video, error) {
	return uc.videoRepo.List(ctx, userID, limit, offset)
}
