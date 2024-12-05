package storage

import (
	"encoding/json"
	"fmt"
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

func GetFirebaseConfig(w http.ResponseWriter, r *http.Request) {
    allowedOrigin := "https://portfoliocristianarch.web.app"
    w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	ctx := r.Context()
	bucketName := os.Getenv("FIREBASE_BUCKET_NAME")
	fileName := os.Getenv("FIREBASE_FILE_NAME")

	if bucketName == "" || fileName == "" {
		http.Error(w, "Environment variables not set", http.StatusInternalServerError)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating client: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	reader, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading object: %v", err), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	var config FirebaseConfig
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
