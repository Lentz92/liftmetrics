package server

import "github.com/gin-gonic/gin"

// setupRoutes defines all the routes for our server.
func (s *Server) setupRoutes() {
	// Add a root route
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to LiftMetrics API",
		})
	})

	// Create a route group for version 1 of our API
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/lifters", s.GetLifters)                             // GET /api/v1/lifters
		v1.GET("/lifters/:name", s.GetLifterDetails)                 // GET /api/v1/lifters/{name}
		v1.GET("/lifters/:name/performance", s.GetLifterPerformance) // GET /api/v1/lifters/{name}/performance
		v1.GET("/lifters/:name/stats", s.GetLifterStats)             // GET /api/v1/lifters/{name}/stats
	}
}
