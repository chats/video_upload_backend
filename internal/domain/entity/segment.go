package entity

import (
	"time"
)

// Resolution defines video resolution quality
type Resolution string

const (
	Resolution1080p Resolution = "1080p"
	Resolution720p  Resolution = "720p"
)

// Segment represents a transcoded video segment
type Segment struct {
	ID           string     `json:"id"`
	VideoID      string     `json:"video_id"`
	FileName     string     `json:"file_name"`
	URL          string     `json:"url"`
	Resolution   Resolution `json:"resolution"`
	StartTime    float64    `json:"start_time"`
	Duration     float64    `json:"duration"` // Typically 10 seconds as specified
	SegmentIndex int        `json:"segment_index"`
	CreatedAt    time.Time  `json:"created_at"`
}
