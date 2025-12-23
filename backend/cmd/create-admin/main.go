package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: create-admin <username> <password>")
	}

	username := os.Args[1]
	password := os.Args[2]

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", "photowall.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO admins (username, password_hash) VALUES (?, ?)",
		username,
		string(hash),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin user created:", username)
}
