package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
)

const SentinelKey = "GORYUK"

type validEnv struct {
	file    string
	envVars int
}

func Validate(envFiles []string) ([]validEnv, error) {
	validEnvs := []validEnv{}
	for _, value := range envFiles {
		env, err := godotenv.Read(value)
		if err != nil {
			return validEnvs, fmt.Errorf("error loading %s: %w", value, err)
		}
		envInfo := validEnv{file: value, envVars: len(env)}
		_, ok := env[SentinelKey]
		if ok {
			validEnvs = append(validEnvs, envInfo)
		}
	}
	return validEnvs, nil
}
