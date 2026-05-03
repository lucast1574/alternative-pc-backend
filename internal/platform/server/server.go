package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.Default()
	return &Server{
		Engine: engine,
	}
}

func (s *Server) Run(port string) error {
	return s.Engine.Run(":" + port)
}
