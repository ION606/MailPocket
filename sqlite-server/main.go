package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

const PORT = ":3000"

func main() {
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

	http.ListenAndServe(PORT, nil)
}
