package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var re = regexp.MustCompile("[^a-zA-Z0-9-]+")

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

type CreatePost struct {
	Title       string    `json:"title"`
	DisplayTime time.Time `json:"display_time,omitempty"`
}

func toSlug(title string) string {
	return strings.ToLower(re.ReplaceAllString(strings.ReplaceAll(title, " ", "-"), ""))
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
		return Post{}, fmt.Errorf("error fetching post: %w", err)
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

func ToDisplay(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		posts, err := GetPublishedPosts(db)
		if err != nil {
			log.Println("error getting posts:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(posts)
	}
}

func Create(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post CreatePost

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			log.Println("error decoding request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if post.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		if post.DisplayTime.IsZero() {
			post.DisplayTime = time.Now().UTC()
		}

		slug := toSlug(post.Title)

		_, err := db.Exec(`INSERT INTO posts (title, slug, display_time) VALUES (?, ?, ?)`, post.Title, slug, post.DisplayTime)
		if err != nil {
			log.Println("error creating post:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("created post: %+v\n slug: %s", post, slug)
		w.Write([]byte(slug))
	}
}

func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		posts, err := GetPosts(db)
		if err != nil {
			log.Println("error getting posts:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(posts)
	}
}

func SinglePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		post, err := GetPost(db, r.PathValue("slug"))
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			} else {
				log.Println("error getting post:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		json.NewEncoder(w).Encode(post)
	}
}

func DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "Missing slug parameter", http.StatusBadRequest)
			return
		}

		if err := DeletePostDb(db, slug); err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			log.Printf("Error deleting post: %v", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type UpdateDisplayTimeRequest struct {
	DisplayTime time.Time `json:"display_time"`
}

func UpdateStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "Missing slug parameter", http.StatusBadRequest)
			return
		}

		var req UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Status != "draft" && req.Status != "published" {
			http.Error(w, "Invalid status. Must be 'draft' or 'published'", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE posts SET status = ? WHERE slug = ?", req.Status, slug)
		if err != nil {
			http.Error(w, "Failed to update post status", http.StatusInternalServerError)
			log.Printf("Error updating post status: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateDisplayTime(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "Missing slug parameter", http.StatusBadRequest)
			return
		}

		var req UpdateDisplayTimeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE posts SET display_time = ? WHERE slug = ?", req.DisplayTime, slug)
		if err != nil {
			http.Error(w, "Failed to update display time", http.StatusInternalServerError)
			log.Printf("Error updating display time: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
