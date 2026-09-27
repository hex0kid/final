package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	port := os.Getenv("TODO_PORT")
	password := os.Getenv("TODO_PASSWORD")

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("database init: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close: %v", err)
		}
	}()

	if err := server.Run(port, password); err != nil {
		log.Fatalf("server: %v", err)
	}
}
