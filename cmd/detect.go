/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detects if the current directory can run ryuk",
	Long: `This command runs the check command on the working directory and determines
	whether the project can run ryuk.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		found, err := CheckDir()
		if err != nil {
			return err
		}
		if found {
			fmt.Println("The current directory is ready to run ryuk commands!")
		} else {
			fmt.Println("This directory does not have ryuk configured run `go-ryuk init`")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
