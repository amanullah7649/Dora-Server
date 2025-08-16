package handler

import (
	"log"
	"net/http"

	"dora-server/pkg"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize MongoDB if not already connected
	if pkg.MongoClient() == nil {
		if err := pkg.InitializeMongo(); err != nil {
			log.Printf("DB init failed: %v", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
	}

	// Call the main handler
	pkg.Handler(w, r)
}
