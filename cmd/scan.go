package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Scan(startDir, environment string) ([]string, error) {
	var envFiles []string
	dir := startDir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return envFiles, fmt.Errorf("error getting current directory: %w", err)
		}
	}

	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			return envFiles, fmt.Errorf("error reading directory: %w", err)
		}

		found := false
		for _, file := range files {
			name := file.Name()
			if !strings.Contains(name, ".env") {
				continue
			}
			if environment != "" && !strings.Contains(name, environment) {
				continue
			}
			envFilePaths := filepath.Join(dir, name)
			envFiles = append(envFiles, envFilePaths)
			found = true
		}

		if found {
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return envFiles, nil
}
