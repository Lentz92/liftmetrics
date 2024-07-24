package web

import (
	"database/sql"
	"fmt"
	"liftmetrics/internal/db"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	LifterNames []string
	DB          *sql.DB
	Router      *gin.Engine
}

func NewServer(lifterNames []string, db *sql.DB, indexHTMLPath string) *Server {
	s := &Server{
		LifterNames: lifterNames,
		DB:          db,
		Router:      gin.Default(),
	}

	// Load the specific index.html file
	s.Router.LoadHTMLFiles(indexHTMLPath)

	// Set up static file serving
	staticPath := filepath.Join(filepath.Dir(indexHTMLPath), "..", "static")
	s.Router.Static("/static", staticPath)

	// Set up routes
	s.Router.GET("/", s.handleRoot)
	s.Router.GET("/api/search", s.handleSearch)
	s.Router.GET("/api/lifter-details", s.handleLifterDetails)

	return s
}

func (s *Server) Start(addr string) error {
	return s.Router.Run(addr)
}

func (s *Server) handleRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "LiftMetrics",
	})
}

func (s *Server) handleSearch(c *gin.Context) {
	query := normalizeSpace(strings.ToLower(strings.TrimSpace(c.Query("q"))))
	log.Printf("Received search query: %q", query) // Debugging line

	var results []string
	for _, name := range s.LifterNames {
		normalizedName := normalizeSpace(strings.ToLower(name))
		if strings.Contains(normalizedName, query) {
			results = append(results, name)
			if len(results) >= 25 {
				break
			}
		}
	}

	log.Printf("Number of results: %d", len(results))

	c.JSON(http.StatusOK, results)
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func (s *Server) handleLifterDetails(c *gin.Context) {
	lifterName := c.Query("name")
	if lifterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lifter name is required"})
		return
	}

	details, err := db.GetLifterDetails(c, s.DB, lifterName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching lifter details: %v", err)})
		return
	}

	c.JSON(http.StatusOK, details)
}
