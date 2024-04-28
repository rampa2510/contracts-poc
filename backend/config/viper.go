package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	DB_FILE string `mapstructure:"DB_FILE"`
	PORT    string `mapstructure:"PORT"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			DB_FILE: os.Getenv("DB_FILE"),
			PORT:    os.Getenv("PORT"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.DB_FILE == "" {
		err = errors.New("DB_FILE is required")
		return
	}

	if config.PORT == "" {
		err = errors.New("PORT is required")
		return
	}

	return
}
