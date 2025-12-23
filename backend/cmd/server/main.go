package main

import (
	"log"
	"net/http"

	"photowall/internal/auth"
	"photowall/internal/db"
	internalHttp "photowall/internal/http"
)

func main() {
	auth.Init()

	dbConn := db.Open("photowall.db")

	handler := internalHttp.Router(dbConn)
	handler = auth.Sessions.LoadAndSave(handler) // REQUIRED
	handler = internalHttp.CORS(handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Listening on :8080")
	log.Fatal(server.ListenAndServe())
}
