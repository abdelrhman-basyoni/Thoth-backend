package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ReadEnv(env string) string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Warning Error  loading .env file")
	}

	return os.Getenv(env)
}

func ContainsValue(slice []string, valueToCheck string) bool {
	for _, v := range slice {
		if v == valueToCheck {
			return true
		}
	}
	return false
}
