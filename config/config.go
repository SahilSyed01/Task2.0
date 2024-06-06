package config

import (
	"os"
	"log"
)

type Config struct {
	MongoURI            string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoPoolID       string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		MongoURI:            "mongodb+srv://task3-shreeraj:YIXZaFDnEmHXC3PS@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		CognitoClientID:     "3nhdeivacuskqno6992ar4bikg",
		CognitoClientSecret: "sfhflsbea9vecouaqb9eik4bu6410qijr29uklcphv9mv5k0qvs",
		CognitoPoolID:      "us-east-1_bcezkbKcV",
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: environment variable %s not set, using default value %s", key, defaultValue)
		return defaultValue
	}
	return value
}
