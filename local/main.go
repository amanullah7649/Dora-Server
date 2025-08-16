package main

import (
	"log"
	"net/http"
	"os"

	"dora-server/pkg"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// Initialize MongoDB
	if err := pkg.InitializeMongo(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", pkg.Handler)
	http.HandleFunc("/deployments", pkg.Handler)

	log.Printf("Local server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
