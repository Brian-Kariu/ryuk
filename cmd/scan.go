package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Scan(params ...string) []string {
	var envFiles []string
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return envFiles
	}

	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return envFiles
		}

		found := false
		for _, file := range files {
			if strings.Contains(file.Name(), ".env") || len(params) != 0 && (strings.Contains(file.Name(), params[0]) && strings.Contains(file.Name(), ".env")) {
				envFilePaths := filepath.Join(dir, file.Name())
				envFiles = append(envFiles, envFilePaths)
				fmt.Printf("Found %s file at: %s\n", file.Name(), envFiles[len(envFiles)-1])
				found = true
			}
		}

		if found {
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			fmt.Println(".env files not found in current or parent directories.")
			break
		}

		dir = parentDir
	}

	return envFiles
}
