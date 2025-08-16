package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Deployment model
type Deployment struct {
	CommitHash        string    `json:"commit_hash"`
	CommitSubject     string    `json:"commit_subject"`
	CommitBody        string    `json:"commit_body"`
	CommitTimestamp   string    `json:"commit_timestamp"`
	CommitAuthor      string    `json:"commit_author"`
	CommitAuthorEmail string    `json:"commit_author_email"`
	ReleaseVersion    string    `json:"release_version"`
	PreviousCommit    string    `json:"previous_commit"`
	FilesChanged      string    `json:"files_changed"`
	LinesChanged      string    `json:"lines_changed"`
	JenkinsBuildNum   string    `json:"jenkins_build_number"`
	JenkinsBuildURL   string    `json:"jenkins_build_url"`
	DeploymentStatus  string    `json:"deployment_status"`
	InsertedAt        time.Time `json:"inserted_at"`
}

// Global Mongo client and collection
var mongoClient *mongo.Client
var mongoCollection *mongo.Collection

// InitializeMongo initializes MongoDB connection (call once at server start)
func InitializeMongo() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("MongoDB connection error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping failed: %v", err)
	}

	mongoClient = client
	mongoCollection = client.Database("dora_db").Collection("deployments")
	return nil
}

// MongoClient returns the global MongoDB client
func MongoClient() *mongo.Client {
	return mongoClient
}

// MongoCollection returns the global collection
func MongoCollection() *mongo.Collection {
	return mongoCollection
}

// Handler is the main HTTP handler
func Handler(w http.ResponseWriter, r *http.Request) {
	if mongoClient == nil || mongoCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	switch r.URL.Path {
	case "/deployments":
		handleDeployments(w, r)
	case "/":
		handleRoot(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func handleDeployments(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":
		var dep Deployment
		if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		dep.InsertedAt = time.Now()

		commitSubject := dep.CommitSubject // ex: Merge pull request #4957 from charge-onsite/v1.0.718-main-service-alpha
		releaseVersion := dep.ReleaseVersion
		if releaseVersion == "" {
			// split commit subject with '/' and take last part
			parts := strings.Split(commitSubject, "/")
			if len(parts) > 0 {
				releaseVersion = parts[len(parts)-1]
			}
		}
		dep.ReleaseVersion = releaseVersion
		if _, err := mongoCollection.InsertOne(context.TODO(), dep); err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	

	case "GET":
		filter := bson.M{}
		if author := r.URL.Query().Get("commit_author"); author != "" {
			filter["commit_author"] = author
		}

		var limit int64
		if lStr := r.URL.Query().Get("limit"); lStr != "" {
			if l, err := strconv.Atoi(lStr); err == nil {
				limit = int64(l)
			}
		}

		opts := options.Find().SetSort(bson.D{{"inserted_at", -1}})
		if limit > 0 {
			opts.SetLimit(limit)
		}

		cursor, err := mongoCollection.Find(context.TODO(), filter, opts)
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
