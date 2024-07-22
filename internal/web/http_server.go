package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"liftmetrics/internal/db"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	LifterNames       []string
	DB                *sql.DB
	InterfaceHTMLPath string
}

func NewServer(lifterNames []string, db *sql.DB, interfaceHTMLPath string) *Server {
	return &Server{
		LifterNames:       lifterNames,
		DB:                db,
		InterfaceHTMLPath: interfaceHTMLPath,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleRoot)
	http.HandleFunc("/api/search", s.handleSearch)
	http.HandleFunc("/api/lifter-details", s.handleLifterDetails)

	fmt.Println("Server is running on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	content, err := os.ReadFile(s.InterfaceHTMLPath)
	if err != nil {
		log.Printf("Error reading interface HTML file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

func (s *Server) handleLifterDetails(w http.ResponseWriter, r *http.Request) {
	lifterName := r.URL.Query().Get("name")
	if lifterName == "" {
		http.Error(w, "Lifter name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	details, err := db.GetLifterDetails(ctx, s.DB, lifterName)
	if err != nil {
		log.Printf("Error fetching lifter details for %s: %v", lifterName, err)
		if errors.Is(err, db.ErrNoRows) {
			http.Error(w, fmt.Sprintf("No details found for lifter: %s", lifterName), http.StatusNotFound)
		} else if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Error encoding JSON for lifter %s: %v", lifterName, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
