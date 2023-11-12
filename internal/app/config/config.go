package config

import (
	"safechildhood/pkg/converter"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Database    DatabaseConfig
		Server      ServerConfig      `yaml:"server" required:"true"`
		Map         MapConfig         `yaml:"map"`
		Playgrounds PlaygroundsConfig `yaml:"playgrounds" required:"true"`
		GoogleDrive GoogleDriveConfig `yaml:"googleDrive" required:"true"`
		Form        FormConfig        `yaml:"form"`
		Logger      LoggerConfig      `yaml:"logger"`
	}

	DatabaseConfig struct {
		Uri string `env:"DATABASE_URI" required:"true"`
	}

	ServerConfig struct {
		Port  int  `yaml:"port" default:"8080"`
		Debug bool `env:"DEBUG_MODE" default:"true"`

		Cors CorsConfig `yaml:"cors" required:"true"`
	}

	CorsConfig struct {
		AllowOrigins     []string      `env:"CORS_ALLOW_ORIGINS" required:"true"`
		AllowCredentials bool          `yaml:"allowCredentials" default:"false"`
		AllowMethods     []string      `yaml:"allowMethods" default:"[GET]"`
		AllowHeaders     []string      `yaml:"allowHeaders" default:"[Origin, Content-Length, Content-Type]"`
		ExposeHeaders    []string      `yaml:"exposeHeaders" default:"[]"`
		MaxAge           time.Duration `yaml:"maxAge" default:"10h"`
	}

	MapConfig struct {
		Bounds struct {
			SouthWest converter.Coordinates `yaml:"southWest" json:"southWest"`
			NorthEast converter.Coordinates `yaml:"northEast" json:"northEast"`
		} `yaml:"bounds" json:"bounds"`

		DefaultPosition converter.Coordinates `yaml:"defaultPosition" json:"defaultPosition"`

		Zoom struct {
			Default int `yaml:"default" json:"default"`
			Min     int `yaml:"min" json:"min"`
			Max     int `yaml:"max" json:"max"`
		} `yaml:"zoom" json:"zoom"`
	}

	PlaygroundsConfig struct {
		CriticalTimeLimit time.Duration `yaml:"criticalTimeLimit" default:"24h"`
		PathToFile        string        `yaml:"pathToFile" required:"true"`
	}

	GoogleDriveConfig struct {
		PathToGoogleServiceAccout string `yaml:"pathToGoogleServiceAccount" required:"true"`
		MediaFolderId             string `env:"MEDIA_FOLDER_ID" required:"true"`
	}

	FormConfig struct {
		MaxSize   datasize.ByteSize `yaml:"maxSize" default:"1mb"`
		MaxPhotos int               `yaml:"maxPhotos" default:"3"`
	}

	LoggerConfig struct {
		Output       string `yaml:"output" default:"/var/log/safechildhood.log"`
		OutputErrors string `yaml:"outputErrors" default:"/var/log/safechildhood.log"`
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func New(configsPath ...string) *Config {
	c := new(Config)

	if err := configor.Load(c, configsPath...); err != nil {
		panic(err)
	}

	return c
}
