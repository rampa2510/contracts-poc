package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitialiseValidation() {
	slog.Info("Validatio initiliased")
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func SendResponse(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func ValidateRequest[T any](r *http.Request, v T) (error, map[string]string) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("Error whle decoding json: %w", err), nil
	}

	err := validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err, nil
		}

		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = fmt.Sprintf("Invalid %s", err.Field())
		}
		return nil, validationErrors
	}

	return nil, nil
}
