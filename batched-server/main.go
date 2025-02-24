package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"shared"
	"strings"
	"sync"
	"time"
)

var (
	emailQueue []string
	queueLock  sync.Mutex
	dbdir      string
	PORT       string
	fpath      string
)

func saveEmails() {
	queueLock.Lock()
	defer queueLock.Unlock()

	if len(emailQueue) == 0 {
		return
	}

	fileExists := true
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		fileExists = false
	}

	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open file:", err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write header if file is new
	if !fileExists {
		if err := writer.Write([]string{"email", "timestamp"}); err != nil {
			log.Println("Failed to write header:", err)
			return
		}
	}

	// Write each email with timestamp
	for _, email := range emailQueue {
		timestamp := time.Now().Format(time.RFC3339)
		if err := writer.Write([]string{email, timestamp}); err != nil {
			log.Println("Failed to write email:", err)
			return
		}
	}

	// Clear queue
	emailQueue = emailQueue[:0]
}

// background goroutine to batch writes
func init() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			queueLock.Lock()
			hasEmails := len(emailQueue) > 0
			queueLock.Unlock()
			if hasEmails {
				log.Println("Flushing email queue...")
				saveEmails()
			}
		}
	}()
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	if email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	queueLock.Lock()
	emailQueue = append(emailQueue, email)
	shouldFlush := len(emailQueue) >= 100
	queueLock.Unlock()

	if shouldFlush {
		saveEmails()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "data received"})
}

func main() {
	PORT, dbdir = shared.GetArgs()
	fpath = filepath.Join(dbdir, "emails.csv")

	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Batched Write Server is running"))
	})

	log.Println("Starting server on port", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
