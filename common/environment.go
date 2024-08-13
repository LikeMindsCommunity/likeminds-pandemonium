package common

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

// GoDotEnvVariable to load/read the .env file and return the value of the key
func GoDotEnvVariable(key string) string {
	// load .env file
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// Uncomment this to make it work with IDE debug mode (tested on GoLand)
	//dir = "/Users/<user_name>/path_to_pandemonium_root_directory"
	environmentPath := filepath.Join(dir, ".env")
	envs, err := godotenv.Read(environmentPath)

	if err != nil {
		log.Fatalf("Error reading .env file")
	}
	return envs[key]
}
