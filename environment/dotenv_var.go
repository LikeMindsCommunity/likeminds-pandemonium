package environment

import (
	"path/filepath"

	"log"

	"github.com/joho/godotenv"
)

// GetDotEnvVar loads .env vars and returns value for var
func GetDotEnvVar(key string) string {
	dir, err := filepath.Abs("")
	if err != nil {
		log.Fatal(err)
	}
	dotEnvPath := filepath.Join(dir, ".env")
	dotEnv, err := godotenv.Read(dotEnvPath)

	if err != nil {
		log.Fatalf("Error reading .env file")
	}
	return dotEnv[key]
}
