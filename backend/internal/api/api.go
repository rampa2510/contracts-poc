package api

import (
	"database/sql"
	"fmt"
	"net/http"
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
		w.Write([]byte("Ok"))
	})

	info := fmt.Sprintf("Server running on addr - %s", serverConfig.addr)
	fmt.Println(info)

	server := &http.Server{Addr: serverConfig.addr, Handler: router}

	return server
}
