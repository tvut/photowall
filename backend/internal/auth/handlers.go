package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		var (
			id           int
			passwordHash string
		)

		err := db.QueryRow(
			"SELECT id, password_hash FROM admins WHERE username = ?",
			creds.Username,
		).Scan(&id, &passwordHash)

		if err != nil {
			log.Println("user not found for", creds.Username)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(passwordHash),
			[]byte(creds.Password),
		); err != nil {
			log.Println("password mismatch for", creds.Username, creds.Password)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		log.Println("logging in", creds.Username)
		Sessions.Put(r.Context(), "admin_id", id)

		w.WriteHeader(http.StatusOK)
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("logging out", Sessions.GetInt(r.Context(), "admin_id"))
		_ = Sessions.Destroy(r.Context())
		w.WriteHeader(http.StatusOK)
	}
}

func Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("me", Sessions.GetInt(r.Context(), "admin_id"))
		if Sessions.GetInt(r.Context(), "admin_id") == 0 {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
