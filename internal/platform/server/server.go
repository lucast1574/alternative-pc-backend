package server

import (
	"net/http"
	"time"

	listingHandler "github.com/alternative/backend/internal/modules/listing/handler"
	maintenanceHandler "github.com/alternative/backend/internal/modules/maintenance/handler"
	sellHandler "github.com/alternative/backend/internal/modules/sellrequest/handler"
	userHandler "github.com/alternative/backend/internal/modules/user/handler"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/alternative/backend/internal/platform/middleware"
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := engine.Group("/api")

	// Health & Status
	api.GET("/health", func(c *gin.Context) {
		err := database.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "database unreachable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	api.GET("/status", func(c *gin.Context) {
		dbStatus := "connected"
		if err := database.Ping(); err != nil {
			dbStatus = "disconnected"
		}
		c.JSON(http.StatusOK, gin.H{
			"service":  "alternative-pc-backend",
			"version":  "1.0.0",
			"database": dbStatus,
			"dbName":   database.DB.Name(),
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Auth (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/google", userHandler.GoogleAuth)

	// Public listings
	api.GET("/listings", listingHandler.GetApprovedListings)
	api.GET("/listings/:id", listingHandler.GetListing)

	// Public price estimates
	api.GET("/sell-requests/estimates", sellHandler.GetPriceEstimates)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// User
		protected.GET("/me", userHandler.GetMe)
		protected.PATCH("/me", userHandler.UpdateProfile)

		// Sell & maintenance (auth required to submit)
		protected.POST("/sell-requests", sellHandler.CreateSellRequest)
		protected.POST("/maintenance-requests", maintenanceHandler.CreateMaintenanceRequest)

		// Listings (seller/admin)
		sellerRoutes := protected.Group("/listings")
		sellerRoutes.Use(middleware.SellerOrAdmin())
		{
			sellerRoutes.POST("", listingHandler.CreateListing)
		}

		// My listings
		protected.GET("/my-listings", listingHandler.GetMyListings)

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.PATCH("/users/:id/role", userHandler.UpdateRole)

			admin.GET("/listings/pending", listingHandler.GetPendingListings)
			admin.PATCH("/listings/:id/approve", listingHandler.ApproveListing)
			admin.PATCH("/listings/:id/reject", listingHandler.RejectListing)

			admin.GET("/sell-requests", sellHandler.GetSellRequests)
			admin.PATCH("/sell-requests/:id", sellHandler.UpdateSellRequest)

			admin.GET("/maintenance-requests", maintenanceHandler.GetMaintenanceRequests)
			admin.PATCH("/maintenance-requests/:id", maintenanceHandler.UpdateMaintenanceRequest)
		}
	}

	return &Server{Engine: engine}
}

func (s *Server) Run(port string) error {
	return s.Engine.Run(":" + port)
}
