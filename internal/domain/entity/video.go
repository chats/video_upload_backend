package entity

import (
	"time"
)

// VideoStatus represents the processing status of a video
type VideoStatus string

const (
	StatusPending    VideoStatus = "pending"
	StatusUploaded   VideoStatus = "uploaded"
	StatusProcessing VideoStatus = "processing"
	StatusTranscoded VideoStatus = "transcoded"
	StatusSegmented  VideoStatus = "segmented"
	StatusComplete   VideoStatus = "complete"
	StatusFailed     VideoStatus = "failed"
)

// Video represents a video entity in the system
type Video struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Duration       float64     `json:"duration"`
	OriginalURL    string      `json:"original_url"`
	ThumbnailURL   string      `json:"thumbnail_url"`
	Status         VideoStatus `json:"status"`
	FileSize       int64       `json:"file_size"`
	MimeType       string      `json:"mime_type"`
	UserID         string      `json:"user_id"`
	ResolutionInfo string      `json:"resolution_info"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
