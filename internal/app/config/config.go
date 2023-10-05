package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUri string

	PathToPlaygroundsFile string

	PathToGoogleServiceAccout string
	MediaFolderId             string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func New() *Config {
	c := new(Config)

	c.DatabaseUri = os.Getenv("DATABASE_URI")

	c.PathToPlaygroundsFile = os.Getenv("PATH_TO_PLAYGROUNDS_FILE")

	c.PathToGoogleServiceAccout = os.Getenv("PATH_TO_GOOGLE_SERVICE_ACCOUNT")
	c.MediaFolderId = os.Getenv("MEDIA_FOLDER_ID")

	return c
}
