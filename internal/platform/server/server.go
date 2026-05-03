package server

import (
	"net/http"
	"time"

	"github.com/alternative/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.Default()

	// CORS middleware
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check - pings MongoDB
	engine.GET("/api/health", func(c *gin.Context) {
		err := database.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "database unreachable",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Status endpoint
	engine.GET("/api/status", func(c *gin.Context) {
		dbStatus := "connected"
		err := database.Ping()
		if err != nil {
			dbStatus = "disconnected: " + err.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"service":  "alternative-pc-backend",
			"version":  "1.0.0",
			"database": dbStatus,
			"dbName":   database.DB.Name(),
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	return &Server{
		Engine: engine,
	}
}

func (s *Server) Run(port string) error {
	return s.Engine.Run(":" + port)
}
