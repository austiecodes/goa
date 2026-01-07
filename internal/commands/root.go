package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gomor",
	Short: "gomor is a MCP server for memory management",
	Long:  `gomor is a MCP (Model Context Protocol) server that provides memory management capabilities.`,
}

// AddCommand adds a subcommand to the root command
func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
