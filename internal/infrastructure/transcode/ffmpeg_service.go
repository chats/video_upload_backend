package transcode

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"cams.dev/video_upload_backend/internal/domain/entity"
	"cams.dev/video_upload_backend/internal/domain/repository"
)

// FFmpegService implements repository.TranscodeRepository using FFmpeg
type FFmpegService struct {
	ffmpegPath  string
	ffprobePath string
}

// NewFFmpegService creates a new FFmpeg service
func NewFFmpegService(ffmpegPath, ffprobePath string) repository.TranscodeRepository {
	return &FFmpegService{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
	}
}

// getResolutionParams returns width and height for a given resolution
func (s *FFmpegService) getResolutionParams(resolution entity.Resolution) (int, int) {
	switch resolution {
	case entity.Resolution1080p:
		return 1920, 1080
	case entity.Resolution720p:
		return 1280, 720
	default:
		return 1280, 720 // Default to 720p
	}
}

// Transcode transcodes a video to a specific resolution and fps
func (s *FFmpegService) Transcode(
	ctx context.Context,
	inputURL string,
	outputPath string,
	resolution entity.Resolution,
	fps int,
) error {
	width, height := s.getResolutionParams(resolution)

	// Prepare the FFmpeg command
	args := []string{
		"-i", inputURL,
		"-c:v", "libx264",
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
		"-r", strconv.Itoa(fps),
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-y", // Overwrite output file if it exists
		outputPath,
	}

	// Create the command
	cmd := exec.CommandContext(ctx, s.ffmpegPath, args...)

	// Capture stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg transcode failed: %w", err)
	}

	return nil
}

// Segment segments a video into parts of specified duration
func (s *FFmpegService) Segment(
	ctx context.Context,
	videoPath string,
	segmentDuration int,
	outputPattern string,
) ([]string, error) {
	// Create directory for output pattern if it doesn't exist
	dir := filepath.Dir(outputPattern)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare the FFmpeg command for segmenting
	args := []string{
		"-i", videoPath,
		"-c", "copy", // Copy without re-encoding
		"-map", "0",
		"-f", "segment",
		"-segment_time", strconv.Itoa(segmentDuration),
		"-segment_format", "mpegts",
		"-segment_list", filepath.Join(dir, "playlist.m3u8"),
		"-segment_list_type", "m3u8",
		outputPattern,
	}

	// Create the command
	cmd := exec.CommandContext(ctx, s.ffmpegPath, args...)

	// Capture stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg segment failed: %w", err)
	}

	// Get list of segment files
	pattern := strings.Replace(outputPattern, "%03d", "*", 1)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list segment files: %w", err)
	}

	return matches, nil
}

// GetVideoInfo gets information about a video file
func (s *FFmpegService) GetVideoInfo(
	ctx context.Context,
	videoPath string,
) (duration float64, width int, height int, err error) {
	// Prepare the FFprobe command for getting duration
	durationArgs := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	}

	// Execute the command for duration
	durationCmd := exec.CommandContext(ctx, s.ffprobePath, durationArgs...)
	durationOutput, err := durationCmd.Output()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ffprobe duration failed: %w", err)
	}

	// Parse duration
	durationStr := strings.TrimSpace(string(durationOutput))
	duration, err = strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	// Prepare the FFprobe command for getting video dimensions
	dimensionArgs := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		videoPath,
	}

	// Execute the command for dimensions
	dimensionCmd := exec.CommandContext(ctx, s.ffprobePath, dimensionArgs...)
	dimensionOutput, err := dimensionCmd.Output()
	if err != nil {
		return duration, 0, 0, fmt.Errorf("ffprobe dimensions failed: %w", err)
	}

	// Parse dimensions
	dimensions := strings.TrimSpace(string(dimensionOutput))
	parts := strings.Split(dimensions, "x")
	if len(parts) != 2 {
		return duration, 0, 0, fmt.Errorf("unexpected dimensions format: %s", dimensions)
	}

	width, err = strconv.Atoi(parts[0])
	if err != nil {
		return duration, 0, 0, fmt.Errorf("failed to parse width: %w", err)
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		return duration, width, 0, fmt.Errorf("failed to parse height: %w", err)
	}

	return duration, width, height, nil
}
