package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	var PORT string
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	} else {
		PORT = "15521"
	}

	db, _ := sql.Open("sqlite", "emails.db")
	db.Exec(`CREATE TABLE IF NOT EXISTS emails (
		email TEXT PRIMARY KEY,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))
		if email == "" {
			http.Error(w, "Email required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`INSERT OR IGNORE INTO emails(email) VALUES (?)`, email)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "storage failed"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "data received"})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SQLite Write Server is running"))
	})

	log.Println("Starting server on port", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
