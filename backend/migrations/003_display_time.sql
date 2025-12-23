-- Add display_time column to posts
ALTER TABLE posts ADD COLUMN display_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
