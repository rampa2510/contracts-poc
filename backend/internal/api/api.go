package api

import (
	"database/sql"
	"net/http"

	"github.com/rampa2510/contracts-poc/config"
	"github.com/rampa2510/contracts-poc/internal/api/middleware"
	"github.com/rampa2510/contracts-poc/internal/api/services/contracts"
	"github.com/rampa2510/contracts-poc/internal/api/services/user"
	"github.com/rampa2510/contracts-poc/internal/utils"
	"golang.org/x/exp/slog"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (serverConfig *APIServer) InitaliseHTTPServer(env *config.EnvVars, awsClient *utils.AwsClient) *http.Server {
	utils.InitialiseValidation()

	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status": "OK",
		}
		utils.SendResponse(w, http.StatusOK, response)
	})

	// Initalise user routes
	userStorage := user.NewUserDb(serverConfig.db)
	userController := user.NewUserController(userStorage)
	user.RegisterUserRouter(router, userController)

	// Initalise contracts routes
	contractsStorage := contracts.NewContractsStorage(serverConfig.db)
	contractsController := contracts.NewContractsController(contractsStorage, awsClient, env)
	contracts.RegisterContractsRouter(router, contractsController)

	slog.Info("Server running", "addr", serverConfig.addr)

	stack := middleware.CreateStack(middleware.Logging, middleware.ErrorHandling, middleware.CorsMiddleware)

	server := &http.Server{Addr: serverConfig.addr, Handler: stack(router)}

	return server
}
