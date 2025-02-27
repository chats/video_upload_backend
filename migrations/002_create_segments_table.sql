CREATE TABLE IF NOT EXISTS segments (
    id VARCHAR(36) PRIMARY KEY,
    video_id VARCHAR(36) NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    resolution VARCHAR(10) NOT NULL,
    start_time FLOAT NOT NULL,
    duration FLOAT NOT NULL,
    segment_index INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (video_id, resolution, segment_index)
);

CREATE INDEX IF NOT EXISTS idx_segments_video_id ON segments(video_id);
