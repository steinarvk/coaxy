package cmd

import "github.com/spf13/cobra"

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "commands for debugging coaxy",
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
