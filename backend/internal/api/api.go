package api

import (
	"database/sql"
	"net/http"

	// "github.com/rampa2510/contracts-poc/internal/api/middleware"
	// "github.com/rampa2510/contracts-poc/internal/api/services/user"
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

func (serverConfig *APIServer) InitaliseHTTPServer() *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status": "OK",
		}
		utils.SendResponse(w, http.StatusOK, response)
	})

	router.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"status": "Ok",
		}
		utils.SendResponse(w, http.StatusOK, resp)
	})

	// userStorage := user.NewUserDb(serverConfig.db)
	// userController := user.NewUserController(userStorage)
	// user.RegisterUserRouter(router, userController)

	slog.Info("Server running", "addr", serverConfig.addr)

	// stack := middleware.CreateStack(middleware.Logging, middleware.ErrorHandling)

	server := &http.Server{Addr: serverConfig.addr, Handler: router}

	return server
}
