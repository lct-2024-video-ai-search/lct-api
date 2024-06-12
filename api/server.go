package api

import (
	"github.com/gin-gonic/gin"
	db "lct-backend/db/sqlc"
	"net/http"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	client http.Client
}

func NewServer(store db.Store, client http.Client) (*Server, error) {
	server := &Server{
		store:  store,
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
