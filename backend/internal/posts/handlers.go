package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var re = regexp.MustCompile("[^a-zA-Z0-9-]+")

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

func List(db *sql.DB) http.HandlerFunc {
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

func Create(w http.ResponseWriter, r *http.Request) {
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
		post.DisplayTime = time.Now()
	}

	slug := toSlug(post.Title)

	log.Printf("creating post: %+v\n slug: %s", post, slug)
	w.Write([]byte(slug))
}
