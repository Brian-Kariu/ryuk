package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Create(key string, value string, params ...string) {
	value += "\n"
	file, err := os.OpenFile(params[0], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening env file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	env, err := godotenv.Read(params[0])
	_, ok := env[key]
	if !ok {
		text := strings.Join([]string{key, value}, "=")
		_, err = writer.WriteString(text)
		if err != nil {
			fmt.Println("Error writing to env file:", err)
			return
		}

	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing the buffer:", err)
		return
	}
	fmt.Printf("Environment %s variable set!\n", key)
}
