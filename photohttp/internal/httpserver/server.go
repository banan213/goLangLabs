package httpserver

import (
	"photohttp/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	service *service.PhotoService
}

func New(service *service.PhotoService) *Server {
	return &Server{service: service}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/app.js", "./web/app.js")
	r.GET("/health", s.health)
	r.GET("/photo/random", s.randomPhoto)
	r.GET("/photo/search", s.searchPhoto)
	r.GET("/photo/last", s.lastPhoto)
}