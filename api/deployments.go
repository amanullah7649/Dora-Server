// package main
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Deployment struct {
	Branch        string    `json:"branch"`
	CommitHash    string    `json:"commit_hash"`
	CommitMessage string    `json:"commit_message"`
	CommitDesc    string    `json:"commit_description"`
	CommitDate    string    `json:"commit_date"`
	AuthorName    string    `json:"author_name"`
	AuthorEmail   string    `json:"author_email"`
	BuildNumber   string    `json:"build_number"`
	BuildURL      string    `json:"build_url"`
	InsertedAt    time.Time `json:"inserted_at"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		http.Error(w, "MONGODB_URI not set", http.StatusInternalServerError)
		return
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		http.Error(w, "MongoDB connection error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("dora_db")
	collection := db.Collection("deployments")

	switch r.URL.Path {
	case "/deployments":
		handleDeployments(w, r, collection)
	case "/":
		handleRoot(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func handleDeployments(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	switch r.Method {
	case "POST":
		var dep Deployment
		if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		dep.InsertedAt = time.Now()

		_, err := collection.InsertOne(context.TODO(), dep)
		if err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	case "GET":
		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.TODO())

		var deployments []Deployment
		if err := cursor.All(context.TODO(), &deployments); err != nil {
			http.Error(w, "Failed to parse data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deployments)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the Dora Matrix Deployment API!"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
