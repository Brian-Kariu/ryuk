package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	Workspaces       []WorkspaceConfig
	CurrentWorkspace string
	CurrentEnv       string
)

func updateWorkspaces(w WorkspaceConfig) error {
	err := checkWorkspaceExists(w.Name)
	if err != nil {
		fmt.Errorf("Workspace config: %s\n", err)
	}

	Workspaces = append(Workspaces, w)
	viper.Set("workspaces", Workspaces)
	if err := viper.WriteConfig(); err != nil {
		fmt.Errorf("Error saving workspaces : %v\n", err)
	}
	return nil
}

// TODO: Link environment struct to this
type WorkspaceConfig struct {
	ID          string              `mapstructure:"id"`
	Name        string              `mapstructure:"name"`
	Path        string              `mapstructure:"path"`
	Environment map[string]struct{} `mapstructure:"environment"`
}

func DeleteWorkspace(id string) {
	for i, ws := range Workspaces {
		if ws.ID == id {
			Workspaces = append(Workspaces[:i], Workspaces[i+1:]...)
			break
		}
	}
	viper.Set("workspaces", Workspaces)
	if err := viper.WriteConfig(); err != nil {
		fmt.Errorf("Error saving workspaces : %v\n", err)
	}
}

func UpdateWorkspace(name, env string) {
	currentWorkspaceIndex := 0
	ws, err := GetWorkspace(name)
	envSet := map[string]struct{}{}
	if err != nil {
		fmt.Println("Error fetching workspace, %s", err)
	}
	if len(ws.Environment) == 0 {
		envSet[env] = struct{}{}
		ws.Environment = envSet
	}
	if len(ws.Environment) != 0 {
		ws.Environment[env] = struct{}{}
	}
	for i, cw := range Workspaces {
		if ws.ID == cw.ID {
			currentWorkspaceIndex = i
		}
	}
	Workspaces[currentWorkspaceIndex].Environment = ws.Environment
	fmt.Println("Current env var: %s", ws.Environment)
	fmt.Println("Current workspaces: %s", Workspaces)
	viper.Set("workspaces", Workspaces)
	if err := viper.WriteConfig(); err != nil {
		fmt.Errorf("Error saving workspaces : %v\n", err)
	}
}

func GetWorkspace(name string) (WorkspaceConfig, error) {
	for _, ws := range Workspaces {
		if ws.Name == name {
			return ws, nil
		}
	}
	return WorkspaceConfig{}, fmt.Errorf("Couldn't find workspace ", name)
}

func checkWorkspaceExists(name string) error {
	err := viper.UnmarshalKey("workspaces", &Workspaces)
	if err != nil {
		fmt.Println("Error fetching workspaces: %v", err)
	}

	for _, ws := range Workspaces {
		if ws.Name == name {
			return fmt.Errorf("Workspace '%s' already exists.\n", name)
		}
	}
	return nil
}

func NewWorkspaceConfig(name string, environment []string) {
	filePath := filepath.Join(BasePath, name)
	id := uuid.New().String()
	envSet := map[string]struct{}{}
	for _, env := range environment {
		envSet[env] = struct{}{}
	}
	newWorkspace := WorkspaceConfig{
		ID:          id,
		Name:        name,
		Path:        filePath,
		Environment: envSet,
	}
	err := updateWorkspaces(newWorkspace)
	if err != nil {
		log.Println("Error creating workspace %s", name)
	}
	log.Printf("%s workspace has been created.\n", name)
}
