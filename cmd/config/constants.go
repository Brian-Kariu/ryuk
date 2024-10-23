package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var BasePath string = ""

func InitConstants() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	BasePath = filepath.Join(home, ".ryuk/")
	if BasePath == "" {
		log.Fatal("Error assigning base path config.")
	}
}
