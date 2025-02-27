CREATE TABLE IF NOT EXISTS videos (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration FLOAT,
    original_url TEXT NOT NULL,
    thumbnail_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    resolution_info VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos(user_id);
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
