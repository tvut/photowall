package http

import (
	"database/sql"
	"net/http"

	"photowall/internal/auth"
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

	// admin
	adminRouter := http.NewServeMux()
	adminRouter.Handle("POST /add-post", posts.Create(db))
	adminRouter.Handle("GET /all-posts", posts.List(db))

	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", auth.RequireAuth(adminRouter)))

	return mux
}
