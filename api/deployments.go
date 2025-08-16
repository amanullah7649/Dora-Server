package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Deployment model aligned with Jenkins payload
type Deployment struct {
	CommitHash       string    `json:"commit_hash"`
	CommitSubject    string    `json:"commit_subject"`
	CommitBody       string    `json:"commit_body"`
	CommitTimestamp  string    `json:"commit_timestamp"`
	CommitAuthor     string    `json:"commit_author"`
	CommitAuthorEmail string   `json:"commit_author_email"`
	ReleaseVersion   string    `json:"release_version"`
	PreviousCommit   string    `json:"previous_commit"`
	FilesChanged     string    `json:"files_changed"`
	LinesChanged     string    `json:"lines_changed"`
	JenkinsBuildNum  string    `json:"jenkins_build_number"`
	JenkinsBuildURL  string    `json:"jenkins_build_url"`
	InsertedAt       time.Time `json:"inserted_at"`
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
		// Filter by query parameters
		filter := bson.M{}
		commitAuthor := r.URL.Query().Get("commit_author")
		if commitAuthor != "" {
			filter["commit_author"] = commitAuthor
		}

		limit := int64(0)
		limitParam := r.URL.Query().Get("limit")
		if limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil {
				limit = int64(l)
			}
		}

		opts := options.Find().SetSort(bson.D{{"inserted_at", -1}})
		if limit > 0 {
			opts.SetLimit(limit)
		}

		cursor, err := collection.Find(context.TODO(), filter, opts)
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
