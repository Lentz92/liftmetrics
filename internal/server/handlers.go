package server

import (
	"context"
	"errors"
	"liftmetrics/internal/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultTimeout = 10 * time.Second

// GetLifters handles the GET request for retrieving all lifter names.
func (s *Server) GetLifters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
	defer cancel()

	lifters, err := db.GetAllLifters(ctx, s.db)
	if err != nil {
		handleDatabaseError(c, err, "Failed to fetch lifters")
		return
	}
	c.JSON(http.StatusOK, gin.H{"lifters": lifters})
}

// GetLifterDetails handles the GET request for retrieving detailed information about a specific lifter.
func (s *Server) GetLifterDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
	defer cancel()

	name := c.Param("name")

	details, err := db.GetLifterDetails(ctx, s.db, name)
	if err != nil {
		handleDatabaseError(c, err, "Failed to fetch lifter details")
		return
	}

	c.JSON(http.StatusOK, gin.H{"lifterDetails": details})
}

// GetLifterPerformance handles the GET request for a specific lifter's performance.
func (s *Server) GetLifterPerformance(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
	defer cancel()

	name := c.Param("name")

	performance, err := db.GetLifterPerformanceOverTime(ctx, s.db, name)
	if err != nil {
		handleDatabaseError(c, err, "Failed to fetch lifter performance")
		return
	}

	c.JSON(http.StatusOK, performance)
}

// GetLifterStats handles the GET request for a specific lifter's statistics.
func (s *Server) GetLifterStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
	defer cancel()

	name := c.Param("name")

	stats, err := db.GetLifterStats(ctx, s.db, name)
	if err != nil {
		handleDatabaseError(c, err, "Failed to fetch lifter stats")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// handleDatabaseError is a helper function to handle database errors consistently.
func handleDatabaseError(c *gin.Context, err error, message string) {
	if errors.Is(err, db.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lifter not found"})
	} else if errors.Is(err, context.DeadlineExceeded) {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
	}
}
