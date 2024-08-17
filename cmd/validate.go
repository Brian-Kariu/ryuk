package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Validate(envFiles []string) int {
	for _, value := range envFiles {
		env, err := godotenv.Read(value)
		fmt.Printf("Found %d env vars from %s\n", len(env), value)
		if err != nil {
			fmt.Printf("Error loading %s", err)
		}
	}
	fmt.Println("Env files are valid!")
	return 1
}
