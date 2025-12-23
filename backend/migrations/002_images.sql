-- Create images table
CREATE TABLE images (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  url TEXT NOT NULL
);

-- Create junction table for posts to images
CREATE TABLE post_images (
  post_id INTEGER NOT NULL,
  image_id INTEGER NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (post_id, image_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

-- Remove content column from posts
ALTER TABLE posts DROP COLUMN content;
