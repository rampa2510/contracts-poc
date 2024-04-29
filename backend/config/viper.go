package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type EnvVars struct {
	DB_FILE         string `mapstructure:"DB_FILE" validate:"required"`
	PORT            string `mapstructure:"PORT" validate:"required"`
	AWS_BUCKET_NAME string `mapstructure:"AWS_BUCKET_NAME" validate:"required"`
	AWS_REGION      string `mapstructure:"AWS_REGION" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			DB_FILE:         os.Getenv("DB_FILE"),
			PORT:            os.Getenv("PORT"),
			AWS_BUCKET_NAME: os.Getenv("AWS_BUCKET_NAME"),
			AWS_REGION:      os.Getenv("AWS_REGION"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic("Error while reading viper config")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic("Error while Unmarshaling config")
	}

	err = validate.Struct(config)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		}
		fmt.Println(validationErrors)
		panic("VALIDATION ERROR")
	}

	return
}

