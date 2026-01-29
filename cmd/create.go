package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Create(key string, value string, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening env file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	env, err := godotenv.Read(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading env file: %w", err)
	}
	_, ok := env[key]
	if !ok {
		text := fmt.Sprintf("%s=%s\n", key, value)
		_, err = writer.WriteString(text)
		if err != nil {
			return fmt.Errorf("error writing to env file: %w", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing the buffer: %w", err)
	}
	return nil
}
