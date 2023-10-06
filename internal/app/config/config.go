package config

import (
	"log"
	"os"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Database    DatabaseConfig
		Server      ServerConfig      `yaml:"server"`
		Map         MapConfig         `yaml:"map"`
		Playgrounds PlaygroundsConfig `yaml:"playgrounds"`
		GoogleDrive GoogleDriveConfig `yaml:"googleDrive"`
		Form        FormConfig        `yaml:"form"`
	}

	DatabaseConfig struct {
		Uri string
	}

	ServerConfig struct {
		Port int `yaml:"port"`

		Cors CorsConfig `yaml:"cors"`
	}

	CorsConfig struct {
		AllowOrigins     []string
		AllowCredentials bool          `yaml:"allowCredentials"`
		AllowMethods     []string      `yaml:"allowMethods"`
		AllowHeaders     []string      `yaml:"allowHeaders"`
		MaxAge           time.Duration `yaml:"maxAge"`
	}

	MapConfig struct {
		Bounds struct {
			SouthWest []float64 `yaml:"southWest"`
			NorthEast []float64 `yaml:"northEast"`
		} `yaml:"bounds"`

		DefaultPosition []float64 `yaml:"defaultPosition"`

		Zoom struct {
			Default int `yaml:"default"`
			Min     int `yaml:"min"`
			Max     int `yaml:"max"`
		} `yaml:"zoom"`
	}

	PlaygroundsConfig struct {
		CriticalTimeLimit time.Duration `yaml:"criticalTimeLimit"`
		PathToFile        string        `yaml:"pathToFile"`
	}

	GoogleDriveConfig struct {
		PathToGoogleServiceAccout string `yaml:"pathToGoogleServiceAccount"`
		MediaFolderId             string
	}

	FormConfig struct {
		MaxSize   datasize.ByteSize `yaml:"maxSize"`
		MaxPhotos int               `yaml:"maxPhotos"`
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func New(configPath string) *Config {
	c := new(Config)

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(file, &c); err != nil {
		panic(err)
	}

	c.Database.Uri = os.Getenv("DATABASE_URI")

	c.Server.Cors.AllowOrigins = []string{os.Getenv("CORS_ALLOW_ORIGIN")}

	c.GoogleDrive.MediaFolderId = os.Getenv("MEDIA_FOLDER_ID")

	return c
}
