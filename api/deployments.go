package handler

import (
	"net/http"
	"log"

	"dora-server/internal"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if internal.MongoClient() == nil {
		if err := internal.InitializeMongo(); err != nil {
			log.Printf("DB init failed: %v", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
	}

	internal.Handler(w, r)
}
