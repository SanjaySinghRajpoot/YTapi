-- Active: 1712148142795@@127.0.0.1@5432@ytapi
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    video_title VARCHAR(255) NOT NULL,
    description TEXT,
    publish_time TIMESTAMP NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    channel VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);