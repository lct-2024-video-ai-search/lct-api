package api

import (
	"github.com/gin-gonic/gin"
	"lct-backend/db"
	"net/http"
)

type Server struct {
	store              db.VideoStore
	router             *gin.Engine
	client             http.Client
	videoProcessingURL string
	videoIndexingURL   string
}

func NewServer(store db.VideoStore, videoProcessingURL string, videoIndexingURL string, client http.Client) (*Server, error) {
	server := &Server{
		store:              store,
		client:             client,
		videoIndexingURL:   videoIndexingURL,
		videoProcessingURL: videoProcessingURL,
	}

	// register routes method
	router := gin.Default()
	router.Use(allowAll())
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
	router.GET("/videos", s.videosPaged)
	router.Static("/docs", "swagger")
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorResponse(err error) ErrorResponse {
	return ErrorResponse{err.Error()}
}
