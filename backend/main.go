package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func getFirebaseConfig() map[string]string {
	return map[string]string{
		"apiKey":            os.Getenv("FIREBASE_API_KEY"),
		"authDomain":        os.Getenv("FIREBASE_AUTH_DOMAIN"),
		"projectId":         os.Getenv("FIREBASE_PROJECT_ID"),
		"storageBucket":     os.Getenv("FIREBASE_STORAGE_BUCKET"),
		"messagingSenderId": os.Getenv("FIREBASE_MESSAGING_SENDER_ID"),
		"appId":             os.Getenv("FIREBASE_APP_ID"),
		"measurementId":     os.Getenv("FIREBASE_MEASUREMENT_ID"),
	}
}

func encryptFirebaseConfig(config map[string]string, secretKey string) (string, error) {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("error al serializar JSON: %w", err)
	}

	key := []byte(secretKey)
	if len(key) != 32 {
		return "", fmt.Errorf("la clave debe tener 32 bytes (256 bits)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error al crear el bloque AES: %w", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error al generar el nonce: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error al crear GCM: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, jsonData, nil)

	finalData := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(finalData), nil
}

func main() {
	firebaseConfig := getFirebaseConfig()

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY no está configurada en las variables de entorno")
	}

	encryptedConfig, err := encryptFirebaseConfig(firebaseConfig, secretKey)
	if err != nil {
		log.Fatalf("Error al encriptar configuración: %v", err)
	}

	fmt.Println("Configuración encriptada:", encryptedConfig)
}
