package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Deployment struct {
	Branch          string    `json:"branch"`
	CommitHash      string    `json:"commit_hash"`
	CommitMessage   string    `json:"commit_message"`
	CommitDesc      string    `json:"commit_description"`
	CommitDate      string    `json:"commit_date"`
	AuthorName      string    `json:"author_name"`
	AuthorEmail     string    `json:"author_email"`
	BuildNumber     string    `json:"build_number"`
	BuildURL        string    `json:"build_url"`
	InsertedAt      time.Time `json:"inserted_at"`
}

var collection *mongo.Collection

func main() {
	// MongoDB connection
	serverName := "test"
	mongoURI := "mongodb+srv://doramatrix:doramatrixpassword@cluster0.5ply0tk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0" // Change this to your MongoDB URI if needed
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	collection = client.Database(fmt.Sprintf("%s_dora", serverName)).Collection("deployments")

	http.HandleFunc("/deployments", deploymentsHandler)
	http.HandleFunc("/", projectStart)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func deploymentsHandler(w http.ResponseWriter, r *http.Request) {
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

func projectStart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {


	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the Dora Matrix Deployment API!"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
