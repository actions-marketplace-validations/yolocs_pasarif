package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pasarif",
	Short: "Parse a SARIF",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// add commands.
	rootCmd.AddCommand(checkCmd, queryCmd)
}
