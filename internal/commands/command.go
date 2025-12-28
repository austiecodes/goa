package commands

import (
	mcpcmd "github.com/austiecodes/goa/internal/commands/mcp"
	memorycmd "github.com/austiecodes/goa/internal/commands/memory"
	setcmd "github.com/austiecodes/goa/internal/commands/set"
)

func init() {
	rootCmd.AddCommand(mcpcmd.McpCmd)
	rootCmd.AddCommand(memorycmd.MemoryCmd)
	rootCmd.AddCommand(setcmd.SetCmd)
}
