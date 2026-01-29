package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
)

func RemoveEnv(envs []string, filename string) error {
	env, err := godotenv.Read(filename)
	if err != nil {
		return fmt.Errorf("error reading env file: %w", err)
	}

	for _, v := range envs {
		delete(env, v)
	}
	err = godotenv.Write(env, filename)
	if err != nil {
		return fmt.Errorf("error writing env file: %w", err)
	}
	return nil
}
