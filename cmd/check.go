/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Looks if the current directory has a valid .env file",
	Long: `Inspects the current directory as well as its parents to see if it has a valid .env file. This command
	is run to ensure that ryuk is able to pull environment variables from this file. Note that the .env file or files
	have to be valid.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		found, err := CheckDir()
		if err != nil {
			return err
		}
		if found {
			fmt.Println("Found valid env file(s)")
		} else {
			fmt.Println("No valid env file found configured for ryuk. Run go-ryuk init to configure ryuk for this project.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
