/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset an environment variable",
	Long:  `This command removes one or more environment variable in the current working environment`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := viper.GetString("environment") + ".env"
		err := RemoveEnv(args, filename)
		if err != nil {
			return err
		}
		envVars := strings.Join(args, ", ")
		fmt.Printf("Env variable %s have been unset", envVars)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}
