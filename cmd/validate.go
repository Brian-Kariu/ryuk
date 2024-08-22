package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
)

type validEnv struct {
	file    string
	envVars int
}

func Validate(envFiles []string) []validEnv {
	validEnvs := []validEnv{}
	for _, value := range envFiles {
		env, err := godotenv.Read(value)
		envInfo := validEnv{file: value, envVars: len(env)}
		_, ok := env["GORYUK"]
		if ok {
			validEnvs = append(validEnvs, envInfo)
		}
		if err != nil {
			fmt.Printf("Error loading %s", err)
		}
	}
	return validEnvs
}
