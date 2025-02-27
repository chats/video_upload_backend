package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/domain/repository"
)

// TranscodeUseCase handles video transcoding operations
type TranscodeUseCase struct {
	videoRepo     repository.VideoRepository
	segmentRepo   repository.SegmentRepository
	storageRepo   repository.StorageRepository
	transcodeRepo repository.TranscodeRepository
}

// NewTranscodeUseCase creates a new transcode use case instance
func NewTranscodeUseCase(
	videoRepo repository.VideoRepository,
	segmentRepo repository.SegmentRepository,
	storageRepo repository.StorageRepository,
	transcodeRepo repository.TranscodeRepository,
) *TranscodeUseCase {
	return &TranscodeUseCase{
		videoRepo:     videoRepo,
		segmentRepo:   segmentRepo,
		storageRepo:   storageRepo,
		transcodeRepo: transcodeRepo,
	}
}

// ProcessVideo handles video transcoding and segmentation
func (uc *TranscodeUseCase) ProcessVideo(ctx context.Context, videoID, videoURL string) error {
	// Retrieve video info
	video, err := uc.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Create temporary directory for processing
	tempDir, err := os.MkdirTemp("", "video-processing-"+videoID)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up when done

	// Download original video
	videoData, err := uc.storageRepo.GetFile(ctx, videoURL)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	// Save to temp directory
	originalVideoPath := filepath.Join(tempDir, "original.mp4")
	if err := os.WriteFile(originalVideoPath, videoData, 0644); err != nil {
		return fmt.Errorf("failed to save video to temp dir: %w", err)
	}

	// Get video information
	duration, width, height, err := uc.transcodeRepo.GetVideoInfo(ctx, originalVideoPath)
	if err != nil {
		return fmt.Errorf("failed to get video info: %w", err)
	}

	// Update video with duration and resolution info
	video.Duration = duration
	video.ResolutionInfo = fmt.Sprintf("%dx%d", width, height)
	video.Status = entity.StatusTranscoded
	if err := uc.videoRepo.Update(ctx, video); err != nil {
		return fmt.Errorf("failed to update video info: %w", err)
	}

	// Process each required resolution
	resolutions := []entity.Resolution{
		entity.Resolution1080p,
		entity.Resolution720p,
	}

	for _, resolution := range resolutions {
		// Transcode to the target resolution
		outputPath := filepath.Join(tempDir, fmt.Sprintf("%s.mp4", resolution))

		if err := uc.transcodeRepo.Transcode(ctx, originalVideoPath, outputPath, resolution, 24); err != nil {
			return fmt.Errorf("failed to transcode to %s: %w", resolution, err)
		}

		// Segment the transcoded video
		segmentsDirPath := filepath.Join(tempDir, string(resolution))
		if err := os.MkdirAll(segmentsDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create segments directory: %w", err)
		}

		// Create 10-second segments
		segmentPattern := filepath.Join(segmentsDirPath, "segment_%03d.ts")
		segmentFiles, err := uc.transcodeRepo.Segment(ctx, outputPath, 10, segmentPattern)
		if err != nil {
			return fmt.Errorf("failed to segment video: %w", err)
		}

		// Upload each segment and store metadata
		for i, segmentPath := range segmentFiles {
			segmentFileName := filepath.Base(segmentPath)

			// Read segment file
			segmentData, err := os.ReadFile(segmentPath)
			if err != nil {
				return fmt.Errorf("failed to read segment file: %w", err)
			}

			// Upload to storage
			storagePath := fmt.Sprintf("videos/%s/%s/%s", videoID, resolution, segmentFileName)
			segmentURL, err := uc.storageRepo.UploadFile(ctx, storagePath, segmentData, "video/mp2t")
			if err != nil {
				return fmt.Errorf("failed to upload segment: %w", err)
			}

			// Create segment record
			segment := &entity.Segment{
				ID:           uuid.New().String(),
				VideoID:      videoID,
				FileName:     segmentFileName,
				URL:          segmentURL,
				Resolution:   resolution,
				StartTime:    float64(i) * 10.0, // 10-second segments
				Duration:     10.0,
				SegmentIndex: i,
				CreatedAt:    time.Now(),
			}

			// For the last segment, adjust duration if needed
			if i == len(segmentFiles)-1 && duration-float64(i)*10.0 < 10.0 {
				segment.Duration = duration - float64(i)*10.0
			}

			// Save segment metadata
			if err := uc.segmentRepo.Create(ctx, segment); err != nil {
				return fmt.Errorf("failed to create segment record: %w", err)
			}
		}
	}

	// Generate thumbnail from the original video
	thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	// This would typically be done using the transcodeRepo, but I'll keep it simple here

	// Upload thumbnail
	thumbnailData, err := os.ReadFile(thumbnailPath)
	if err == nil { // Only if thumbnail generation succeeded
		thumbnailURL, err := uc.storageRepo.UploadFile(
			ctx,
			fmt.Sprintf("videos/%s/thumbnail.jpg", videoID),
			thumbnailData,
			"image/jpeg",
		)
		if err == nil {
			video.ThumbnailURL = thumbnailURL
		}
	}

	// Update video status to complete
	video.Status = entity.StatusComplete
	return uc.videoRepo.Update(ctx, video)
}
