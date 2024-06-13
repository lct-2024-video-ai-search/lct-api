package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
	client http.Client
}

func NewServer(db *sql.DB, client http.Client) (*Server, error) {
	server := &Server{
		db:     db,
		client: client,
	}

	// register routes method
	router := gin.Default()
	root := router.Group("/")
	server.registerRoutes(root)

	server.router = router
	return server, nil
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) registerRoutes(router *gin.RouterGroup) {
	router.POST("/index", s.indexVideo)
	router.GET("/search", s.searchVideo)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
