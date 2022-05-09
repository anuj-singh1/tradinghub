package api

import (
	"github.com/gin-gonic/gin"
	"tradingdata/internal/config"
	db "tradingdata/internal/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	global config.GlobalInstance
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(envconfig config.Config, db db.TokenDB) (*Server, error) {
	server := &Server{
		global: config.GlobalInstance{Config: envconfig, TokenDb: db},
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	// Init Router
	router := gin.Default()
	// Route Handlers / Endpoints
	Routes(router, server.global)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start() error {
	return server.router.Run()
}
