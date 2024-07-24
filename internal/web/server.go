package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"liftmetrics/internal/db"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	LifterNames []string
	DB          *sql.DB
	Router      *gin.Engine
	Templates   *template.Template
}

func NewServer(lifterNames []string, db *sql.DB, templatesDir string) *Server {
	s := &Server{
		LifterNames: lifterNames,
		DB:          db,
		Router:      gin.Default(),
	}

	// Parse templates
	s.Templates = template.Must(template.ParseGlob(templatesDir + "/*.html"))

	// Set up routes
	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// Serve static files
	s.Router.Static("/static", "./static")

	// API routes
	api := s.Router.Group("/api")
	{
		api.GET("/search", s.handleSearch)
		api.GET("/lifter-details", s.handleLifterDetails)
	}

	// Serve index page
	s.Router.GET("/", s.handleRoot)
}

func (s *Server) Start(addr string) error {
	fmt.Printf("Server is running on https://%s\n", addr)
	return s.Router.Run(addr)
}

func (s *Server) handleRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Powerlifting Analytics",
	})
}

func (s *Server) handleSearch(c *gin.Context) {
	query := strings.ToLower(strings.TrimSpace(c.Query("q")))
	log.Printf("Received search query: %q", query)

	var results []string
	queryWords := strings.Fields(query)

	numberOfMatches := 25 // Return first 25 lifters for empty query
	if len(queryWords) == 0 {
		results = s.LifterNames[:numberOfMatches]
	} else {
		for _, name := range s.LifterNames {
			if matchesAllWords(name, queryWords) {
				results = append(results, name)
				if len(results) >= numberOfMatches {
					break
				}
			}
		}
	}

	log.Printf("Number of results found: %d", len(results))
	if len(results) > 0 {
		log.Printf("First result: %q", results[0])
	}

	c.JSON(http.StatusOK, results)
}

func (s *Server) handleLifterDetails(c *gin.Context) {
	lifterName := c.Query("name")
	if lifterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lifter name is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	details, err := db.GetLifterDetails(ctx, s.DB, lifterName)
	if err != nil {
		log.Printf("Error fetching lifter details for %s: %v", lifterName, err)
		if errors.Is(err, db.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No details found for lifter"})
		} else if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, details)
}

func matchesAllWords(name string, queryWords []string) bool {
	lowerName := strings.ToLower(name)
	for _, word := range queryWords {
		if !strings.Contains(lowerName, word) {
			return false
		}
	}
	return true
}
