package shared

import (
	"log"
	"os"
)


func GetArgs() (string, string) {
	var PORT string
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	} else {
		PORT = "15521"
	}

	dbdir := "data"

	isDocker := os.Getenv("container") == "docker" || os.Getenv("DOCKER") == "true" || func() bool { _, err := os.Stat("/.dockerenv"); return err == nil }()
	if isDocker {
		dbdir = "/app/data"
	}

	if _, err := os.Stat(dbdir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbdir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	return PORT, dbdir
}