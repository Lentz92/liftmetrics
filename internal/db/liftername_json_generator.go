package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateLifterJSON(ctx context.Context, db *sql.DB, dataDir string) error {
	lifters, err := GetAllLifters(ctx, db)
	if err != nil {
		return fmt.Errorf("getting all lifters: %w", err)
	}

	lifterNames := make([]string, len(lifters))
	for i, lifter := range lifters {
		lifterNames[i] = lifter.Name
	}

	jsonData, err := json.Marshal(lifterNames)
	if err != nil {
		return fmt.Errorf("marshaling lifter names: %w", err)
	}

	jsonFilePath := filepath.Join(dataDir, "lifters.json")
	err = os.WriteFile(jsonFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("writing JSON file: %w", err)
	}

	return nil
}
