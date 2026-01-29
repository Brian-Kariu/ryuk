/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an environment variable",
	Long:  `This command takes a single key/value pair and creates an environment variable.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := viper.GetString("environment") + ".env"
		err := Create(args[0], args[1], fileName)
		if err != nil {
			return err
		}
		fmt.Printf("Env variable has been set for file, %s", fileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
