/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		envFiles := Scan()
		Validate(envFiles)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
