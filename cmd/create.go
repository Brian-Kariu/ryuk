package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Create(key string, value string) {
	file, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening env file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed after writing

	// Create a new writer
	writer := bufio.NewWriter(file)

	// Write text into the file
	text := strings.Join([]string{key, value}, "=")
	_, err = writer.WriteString(text)
	if err != nil {
		fmt.Println("Error writing to env file:", err)
		return
	}

	// Flush the buffer to ensure all data is written
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing the buffer:", err)
		return
	}
	fmt.Println("Environment variable set!")
}
