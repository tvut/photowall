package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	DbFunctions "photowall/internal/db"
	"time"
)

type CreatePost struct {
	Title       string    `json:"title"`
	DisplayTime time.Time `json:"display_time,omitempty"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type UpdateDisplayTimeRequest struct {
	DisplayTime time.Time `json:"display_time"`
}

func ToDisplay(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		posts, err := DbFunctions.GetPublishedPosts(db)
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

		slug, err := DbFunctions.CreatePost(db, post.Title, post.DisplayTime)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(slug))
	}
}

func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		posts, err := DbFunctions.GetPosts(db)
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

		post, err := DbFunctions.GetPost(db, r.PathValue("slug"))
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

		if err := DbFunctions.DeletePostDb(db, slug); err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			log.Printf("Error deleting post: %v", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		var req UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Status != "draft" && req.Status != "published" {
			http.Error(w, "Invalid status. Must be 'draft' or 'published'", http.StatusBadRequest)
			return
		}

		err := DbFunctions.UpdatePostStatus(db, slug, req.Status)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
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

		var req UpdateDisplayTimeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := DbFunctions.UpdatePostDisplayTime(db, slug, req.DisplayTime)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update display time", http.StatusInternalServerError)
			log.Printf("Error updating display time: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
