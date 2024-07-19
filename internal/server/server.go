package server

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"liftmetrics/internal/db"
	"log"
	"sync"
	"time"
)

// Server represents the HTTP server, its router, and database connection.
type Server struct {
	router      *gin.Engine
	db          *sql.DB
	lifterNames []db.LifterName
	mutex       sync.RWMutex
}

// NewServer creates and returns a new Server instance.
func NewServer(database *sql.DB) (*Server, error) {
	server := &Server{
		router: gin.Default(),
		db:     database,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := server.updateLifterNames(ctx); err != nil {
		return nil, err
	}

	go server.periodicallyUpdateLifterNames()

	server.setupRoutes()
	return server, nil
}

// updateLifterNames fetches all lifter names from the database and updates the server's cache.
func (s *Server) updateLifterNames(ctx context.Context) error {
	names, err := db.GetAllLifters(ctx, s.db)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.lifterNames = names
	s.mutex.Unlock()

	return nil
}

// Run starts the HTTP server on the specified address.
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetLifterNames returns the cached list of lifter names.
func (s *Server) GetLifterNames() []db.LifterName {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lifterNames
}

// periodicallyUpdateLifterNames updates the lifter names cache every hour.
func (s *Server) periodicallyUpdateLifterNames() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		if err := s.updateLifterNames(ctx); err != nil {
			log.Printf("Error updating lifter names: %v", err)
		}
		cancel()
	}
}
