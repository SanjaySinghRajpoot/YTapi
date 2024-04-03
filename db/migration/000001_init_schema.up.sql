CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    video_title VARCHAR(255) NOT NULL,
    description TEXT,
    publish_time TIMESTAMP NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    channel VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);