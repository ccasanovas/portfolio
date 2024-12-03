package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
)

type FirebaseConfig struct {
	APIKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
	MeasurementID     string `json:"measurementId"`
}

func getFirebaseConfig(ctx context.Context, bucketName string, fileName string) (*FirebaseConfig, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(fileName)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open object: %v", err)
	}
	defer reader.Close()

	var config FirebaseConfig
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return &config, nil
}

func FirebaseConfigHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := os.Getenv("FIREBASE_BUCKET_NAME")
	fileName := os.Getenv("FIREBASE_FILE_NAME")

	if bucketName == "" || fileName == "" {
		http.Error(w, "FIREBASE_BUCKET_NAME or FIREBASE_FILE_NAME not set in environment variables", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	config, err := getFirebaseConfig(ctx, bucketName, fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching Firebase config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func main() {
	http.HandleFunc("/app/get-config", FirebaseConfigHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
