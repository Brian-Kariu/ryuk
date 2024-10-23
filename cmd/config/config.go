package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var Workspaces []WorkspaceConfig

type WorkspaceConfig struct {
	Name        string   `mapstructure:"name"`
	Path        string   `mapstructure:"path"`
	Environment []string `mapstructure:"environment"`
}

func checkWorkspaceExists(name string) error {
	err := viper.UnmarshalKey("workspaces", &Workspaces)
	if err != nil {
		fmt.Println("Error fetching workspaces:", err)
	}

	for _, ws := range Workspaces {
		if ws.Name == name {
			return fmt.Errorf("Workspace '%s' already exists.\n", name)
		}
	}
	return nil
}

func NewWorkspaceConfig(name string, environment []string) {
	err := checkWorkspaceExists(name)
	if err != nil {
		fmt.Printf("Workspace config: %s\n", err)
	}
	filePath := filepath.Join(BasePath, name+".db")
	newWorkspace := WorkspaceConfig{
		Name:        name,
		Path:        filePath,
		Environment: environment,
	}
	Workspaces = append(Workspaces, newWorkspace)
}
