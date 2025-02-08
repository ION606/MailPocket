package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const PORT = ":3000"

var (
	emailQueue []string
	queueLock  sync.Mutex
)


func saveEmails() {
	queueLock.Lock()
	defer queueLock.Unlock()

	if len(emailQueue) == 0 {
		return
	}

	fileExists := true
	if _, err := os.Stat("emails.csv"); os.IsNotExist(err) {
		fileExists = false
	}

	f, err := os.OpenFile("emails.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// write header if file is new
	if !fileExists {
		if err := writer.Write([]string{"email", "timestamp"}); err != nil {
			return
		}
	}

	// write each email with timestamp
	for _, email := range emailQueue {
		timestamp := time.Now().Format(time.RFC3339)
		if err := writer.Write([]string{email, timestamp}); err != nil {
			return
		}
	}

	// clear queue
	emailQueue = emailQueue[:0]
}

// background goroutine to batch writes
func init() {
	go func() {
		for {
			time.Sleep(5 * time.Second);
			queueLock.Lock();
			hasEmails := len(emailQueue) > 0;
			queueLock.Unlock();
			if hasEmails {
				saveEmails();
			}
		}
	}();	
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
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Batched Write Server is running"))
	})
	http.ListenAndServe(PORT, nil)
}
