package db

import (
	"database/sql"
	"fmt"
	"photowall/internal/utilities"
	"strings"
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayTime time.Time `json:"display_time"`
}

type DisplayPost struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	PhotoUrls   []string  `json:"photos"`
	DisplayTime time.Time `json:"display_time"`
}

func GetPublishedPosts(db *sql.DB) ([]DisplayPost, error) {
	rows, err := db.Query(`
		SELECT p.title, p.slug, p.display_time,
			COALESCE(GROUP_CONCAT(i.url, '||' ORDER BY pi.position), '') as photo_urls
		FROM posts p
		LEFT JOIN post_images pi ON pi.post_id = p.id
		LEFT JOIN images i ON i.id = pi.image_id
		WHERE p.status = 'published'
		GROUP BY p.id
		ORDER BY p.display_time DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []DisplayPost
	for rows.Next() {
		var post DisplayPost
		var photoUrls string
		if err := rows.Scan(&post.Title, &post.Slug, &post.DisplayTime, &photoUrls); err != nil {
			return nil, err
		}

		if photoUrls != "" {
			post.PhotoUrls = strings.Split(photoUrls, "||")
		} else {
			post.PhotoUrls = []string{}
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []DisplayPost{}
	}
	return posts, nil
}

func GetPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT *
		FROM posts p
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Status, &post.CreatedAt, &post.DisplayTime); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []Post{}
	}
	return posts, nil
}

func GetPost(db *sql.DB, slug string) (Post, error) {
	var post Post
	err := db.QueryRow(`
        SELECT id, title, slug, status, created_at, display_time
        FROM posts
        WHERE slug = ?
    `, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Status,
		&post.CreatedAt,
		&post.DisplayTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, err
		}
		return Post{}, err
	}

	return post, nil
}

func DeletePostDb(db *sql.DB, slug string) error {
	result, err := db.Exec(`
        DELETE
        FROM posts
        WHERE slug = ?
    `, slug)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdatePostStatus(db *sql.DB, slug string, status string) error {
	result, err := db.Exec(`
        UPDATE posts 
        SET status = ? 
        WHERE slug = ?
    `, status, slug)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdatePostDisplayTime(db *sql.DB, slug string, displayTime time.Time) error {
	result, err := db.Exec(`
        UPDATE posts 
        SET display_time = ? 
        WHERE slug = ?
    `, displayTime, slug)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func CreatePost(db *sql.DB, title string, displayTime time.Time) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	if displayTime.IsZero() {
		return "", fmt.Errorf("display time cannot be zero")
	}

	slug := utilities.ToSlug(title)

	_, err := db.Exec(`
        INSERT INTO posts (title, slug, display_time, status) 
        VALUES (?, ?, ?, 'draft')
    `, title, slug, displayTime)

	if err != nil {
		return "", fmt.Errorf("error creating post: %w", err)
	}

	return slug, nil
}

func CreateImage(db *sql.DB, url string) (int64, error) {
	result, err := db.Exec(`
        INSERT INTO images (url) 
        VALUES (?)
    `, url)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetImageByUrl(db *sql.DB, url string) (int64, error) {
	var id int64
	err := db.QueryRow(`
        SELECT id FROM images WHERE url = ?
    `, url).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func AttachImageToPost(db *sql.DB, postId int, imageId int64, position int) error {
	_, err := db.Exec(`
        INSERT OR REPLACE INTO post_images (post_id, image_id, position) 
        VALUES (?, ?, ?)
    `, postId, imageId, position)
	return err
}
