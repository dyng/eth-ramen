package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "development"
var Commit = "development"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Long:  "Print version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s, Build Commit: %s", Version, Commit)
		},
	}
}
