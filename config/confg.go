package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// required envs to start up server
	if len(os.Getenv("ENV")) == 0 {
		log.Fatal("Please add ENV variable to your .env file")
	}
	if len(os.Getenv("PORT")) == 0 {
		log.Fatal("Please add PORT variable to your .env file")
	}
}

// get env variables from .env file
// throw error if not found
func StrictEnvVars(key string) string {
	value := os.Getenv(key)
	// if key not found log error
	if len(value) == 0 {
		log.Fatal("Please add " + key + " to your .env file")
	}
	return value
}

// get env variables from .env file
func EnvVars(key string) string {
	return os.Getenv(key)
}
