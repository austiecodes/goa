package commands

import (
	setcmd "github.com/austiecodes/goa/internal/commands/set"
)

func init() {
	rootCmd.AddCommand(setcmd.SetCmd)
}

