package dotenv

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Viper() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var configFile string
	switch env {
	case "docker":
		configFile = "/app/docker.env"
	case "production":
		configFile = "/app/production.env"
	default:
		configFile = ".env"
	}

	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file %s: %w", configFile, err)
	}

	return nil
}
