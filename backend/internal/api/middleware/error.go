package middleware

import (
	"net/http"

	"github.com/rampa2510/contracts-poc/internal/utils"
)

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}

func ErrorHandling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				apiError, ok := err.(*APIError)
				if !ok {
					apiError = &APIError{
						Message: "Internal Server Error",
						Status:  http.StatusInternalServerError,
					}
				}
				utils.SendResponse(w, apiError.Status, apiError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
