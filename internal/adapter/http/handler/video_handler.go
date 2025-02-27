package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"cams.dev/video_upload_backend/internal/usecase"
)

// VideoHandler handles HTTP requests related to videos
type VideoHandler struct {
	videoUseCase *usecase.VideoUseCase
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(videoUseCase *usecase.VideoUseCase) *VideoHandler {
	return &VideoHandler{
		videoUseCase: videoUseCase,
	}
}

// UploadVideo handles video upload requests
func (h *VideoHandler) UploadVideo(c *fiber.Ctx) error {
	// Get file from request
	file, err := c.FormFile("video")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to get video file: "+err.Error())
	}

	// Open and read file
	fileObj, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open video file: "+err.Error())
	}
	defer fileObj.Close()

	// Read file bytes
	buffer := make([]byte, file.Size)
	if _, err := fileObj.Read(buffer); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read video file: "+err.Error())
	}

	// Extract metadata from form
	title := c.FormValue("title", "Untitled Video")
	description := c.FormValue("description", "")

	// Extract user ID from context (would be set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		userID = "anonymous" // Fallback for testing
	}

	// Prepare upload input
	input := usecase.VideoUploadInput{
		Title:       title,
		Description: description,
		FileData:    buffer,
		FileName:    file.Filename,
		FileSize:    file.Size,
		MimeType:    file.Header.Get("Content-Type"),
		UserID:      userID.(string),
	}

	// Call use case
	video, err := h.videoUseCase.UploadVideo(c.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process video: "+err.Error())
	}

	// Return video information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Video upload successful. Processing has begun.",
		"videoId": video.ID,
		"status":  video.Status,
	})
}

// GetVideo handles requests to get video details
func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	videoID := c.Params("id")

	video, err := h.videoUseCase.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Failed to retrieve video: "+err.Error())
	}

	// Get segments
	segments, err := h.videoUseCase.GetVideoSegments(c.Context(), videoID)
	if err != nil {
		// Just log error but continue
		fmt.Printf("Failed to get segments: %v\n", err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"video":    video,
		"segments": segments,
	})
}

// GetVideosByUser handles requests to list videos for a user
func (h *VideoHandler) GetVideosByUser(c *fiber.Ctx) error {
	// Extract user ID from context
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusBadRequest, "User ID not found in request")
	}

	// Parse pagination parameters
	limit := c.QueryInt("limit", 10)  // Default: 10
	offset := c.QueryInt("offset", 0) // Default: 0

	// Call use case
	videos, err := h.videoUseCase.ListVideos(c.Context(), userID.(string), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list videos: "+err.Error())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"videos": videos,
	})
}
