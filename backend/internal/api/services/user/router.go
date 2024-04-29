package user

import (
	"log/slog"
	"net/http"
)

func RegisterUserRouter(router *http.ServeMux, userController *UserController) {
	slog.Info("Initialising User router")

	router.HandleFunc("POST /user", userController.Create)
	router.HandleFunc("GET /user", userController.GetAllUsers)
	router.HandleFunc("GET /user/{id}", userController.GetAUser)
	router.HandleFunc("PUT /user/{id}", userController.GetAUser)

	slog.Info("Initialised User router")
}
