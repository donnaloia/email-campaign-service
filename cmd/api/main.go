package main

import (
	"log"

	"github.com/donnaloia/sendpulse/internal/api"
	"github.com/donnaloia/sendpulse/internal/database"
)

func main() {
	dbConfig := database.NewDefaultConfig()
	db, err := database.Connect(dbConfig.ConnectionString())
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	server := api.NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
