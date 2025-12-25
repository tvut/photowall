package images

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	DbFunctions "photowall/internal/db"
)

type UploadResponse struct {
	ImageUrls []string `json:"image_urls"`
}

func Upload(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := r.ParseMultipartForm(30 << 20) // 30MB max
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["images"]
		if len(files) == 0 {
			http.Error(w, "No files uploaded", http.StatusBadRequest)
			return
		}

		uploadTime := time.Now()

		var uploadedUrls []string

		// upload to year/month directory
		uploadDir := "uploads/" + uploadTime.Format("2006/01/")

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("Failed to create upload directory: %v", err)
			http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
			return
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				log.Printf("Failed to open uploaded file: %v", err)
				continue
			}
			defer file.Close()

			// original file name - hhmmss of upload for uniqueness
			ext := filepath.Ext(fileHeader.Filename)
			baseName := filepath.Base(fileHeader.Filename)
			fileName := strings.TrimSuffix(baseName, ext)
			uniqueName := fmt.Sprintf("%s-%s%s", fileName, uploadTime.Format("030405"), ext)
			filepath := filepath.Join(uploadDir, uniqueName)

			outFile, err := os.Create(filepath)
			if err != nil {
				log.Printf("Failed to create file: %v", err)
				continue
			}

			_, err = io.Copy(outFile, file)
			outFile.Close()
			if err != nil {
				log.Printf("Failed to save file: %v", err)
				os.Remove(filepath)
				continue
			}

			imageUrl := "/" + filepath
			imageId, err := DbFunctions.CreateImage(db, imageUrl)
			if err != nil {
				log.Printf("Failed to save image to database: %v", err)
				os.Remove(filepath)
				continue
			}

			uploadedUrls = append(uploadedUrls, imageUrl)
			log.Printf("Uploaded image: %s (ID: %d)", imageUrl, imageId)
		}

		if len(uploadedUrls) == 0 {
			http.Error(w, "Failed to upload any files", http.StatusInternalServerError)
			return
		}

		response := UploadResponse{ImageUrls: uploadedUrls}
		json.NewEncoder(w).Encode(response)
	}
}

func AttachToPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req struct {
			PostSlug  string   `json:"post_slug"`
			ImageUrls []string `json:"image_urls"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.PostSlug == "" {
			http.Error(w, "Post slug is required", http.StatusBadRequest)
			return
		}

		if len(req.ImageUrls) == 0 {
			http.Error(w, "At least one image URL is required", http.StatusBadRequest)
			return
		}

		post, err := DbFunctions.GetPost(db, req.PostSlug)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get post", http.StatusInternalServerError)
			}
			return
		}

		var errors []string = []string{}
		for i, imageUrl := range req.ImageUrls {
			imageId, err := DbFunctions.GetImageByUrl(db, imageUrl)
			if err != nil {
				log.Printf("Failed to get image ID for URL %s: %v", imageUrl, err)
				errors = append(errors, fmt.Sprintf("Failed to get image ID for URL %s: %v", imageUrl, err))
				continue
			}

			err = DbFunctions.AttachImageToPost(db, post.ID, imageId, i)
			if err != nil {
				log.Printf("Failed to attach image %d to post %d: %v", imageId, post.ID, err)
				errors = append(errors, fmt.Sprintf("Failed to attach image %d to post %d: %v", imageId, post.ID, err))
			}
		}

		if len(errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
