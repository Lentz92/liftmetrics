package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindCSVFile(dataDir string) (string, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "openipf") && strings.HasSuffix(file.Name(), ".csv") {
			return filepath.Join(dataDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no matching CSV file found")
}
