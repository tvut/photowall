package http

import (
	"database/sql"
	"net/http"

	"photowall/internal/auth"
	"photowall/internal/images"
	"photowall/internal/posts"
)

func Router(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// auth
	mux.Handle("POST /api/login", auth.Login(db))
	mux.Handle("POST /api/logout", auth.Logout())
	mux.Handle("GET /api/me", auth.Me())

	// public
	mux.Handle("GET /api/posts", posts.ToDisplay(db))

	// static files (uploads)
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// admin
	adminRouter := http.NewServeMux()
	adminRouter.Handle("POST /add-post", posts.Create(db))
	adminRouter.Handle("GET /posts", posts.List(db))
	adminRouter.Handle("GET /posts/{slug}", posts.SinglePost(db))
	adminRouter.Handle("PUT /posts/{slug}/status", posts.UpdateStatus(db))
	adminRouter.Handle("PUT /posts/{slug}/display-time", posts.UpdateDisplayTime(db))
	adminRouter.Handle("DELETE /posts/{slug}", posts.DeletePost(db))
	adminRouter.Handle("POST /upload-images", images.Upload(db))
	adminRouter.Handle("POST /attach-images", images.AttachToPost(db))

	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", auth.RequireAuth(adminRouter)))

	return CORS(mux)
}
