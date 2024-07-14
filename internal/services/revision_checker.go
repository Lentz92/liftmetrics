// Package services provides utility functions for file operations and data processing.
package services

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ErrRevisionNotFound is returned when a revision number can't be found.
var ErrRevisionNotFound = errors.New("revision not found")

// CheckRevision compares the revision number from the website with the one in the downloaded CSV file.
//
// Parameters:
//   - csvFilePath: The path to the downloaded CSV file.
//   - websiteURL: The URL of the website containing the current revision number.
//
// Returns:
//   - bool: True if an update is needed (revisions don't match), false otherwise.
//   - error: An error if any step of the process fails, nil otherwise.
func CheckRevision(csvFilePath, websiteURL string) (bool, error) {
	log.Println("Checking revision numbers...")

	// Fetch revision from website
	websiteRevision, err := getWebRevision(websiteURL)
	if err != nil {
		log.Printf("Failed to get website revision: %v", err)
		return false, err
	}

	// Extract revision from CSV file
	csvRevision, err := getCSVRevision(csvFilePath)
	if err != nil {
		log.Printf("Failed to get CSV revision: %v", err)
		return false, err
	}

	// Compare revisions
	needsUpdate := websiteRevision != csvRevision
	if needsUpdate {
		log.Printf("Update needed. Website revision: %s, CSV revision: %s", websiteRevision, csvRevision)
	} else {
		log.Printf("No update needed. Website revision: %s, CSV revision: %s", websiteRevision, csvRevision)
	}

	return needsUpdate, nil
}

// getWebRevision fetches the current revision number from the specified website.
//
// Parameters:
//   - url: The URL of the website containing the revision number.
//
// Returns:
//   - string: The revision number found on the website.
//   - error: An error if the revision number couldn't be found or if there was an HTTP error.
func getWebRevision(url string) (string, error) {
	log.Printf("Fetching revision from website: %s", url)

	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Ensure the response body is closed after we're done

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var revision string
	// Find all <li> elements within <ul> elements
	doc.Find("ul li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		// Check if the text starts with "Revision:"
		if strings.HasPrefix(text, "Revision:") {
			// Extract and trim the revision number
			revision = strings.TrimSpace(strings.TrimPrefix(text, "Revision:"))
			// Remove the trailing dot if present
			revision = strings.TrimSuffix(revision, ".")
		}
	})

	// Check if a revision number was found
	if revision == "" {
		log.Println("Revision not found on website")
		return "", ErrRevisionNotFound
	}

	log.Printf("Website revision found: %s", revision)
	return revision, nil
}

// getCSVRevision extracts the revision number from the CSV filename.
//
// Parameters:
//   - filePath: The path to the CSV file.
//
// Returns:
//   - string: The revision number extracted from the CSV filename.
//   - error: An error if the revision number couldn't be extracted.
func getCSVRevision(filePath string) (string, error) {
	log.Printf("Extracting revision from CSV filename: %s", filePath)

	// Get the base name of the file (i.e., just the filename without the path)
	filename := filepath.Base(filePath)

	// Split the filename by hyphens
	parts := strings.Split(filename, "-")

	// The revision should be the last part before the .csv extension
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		revision := strings.TrimSuffix(lastPart, ".csv")

		if revision != "" {
			log.Printf("CSV revision found: %s", revision)
			return revision, nil
		}
	}

	// If we couldn't extract a revision, return an error
	log.Println("Revision not found in CSV filename")
	return "", ErrRevisionNotFound
}
