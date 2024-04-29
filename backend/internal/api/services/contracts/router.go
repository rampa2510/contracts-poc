package contracts

import (
	"log/slog"
	"net/http"
)

func RegisterContractsRouter(router *http.ServeMux, contractsController *ContractsController) {
	slog.Info("Initialising Contracts router")

	router.HandleFunc("POST /contract", contractsController.Create)

	slog.Info("Initialised Contracts router")
}
