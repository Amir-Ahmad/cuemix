package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Default if ldflags not provided
var version = "0.0.0-dev"

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "cuemix",
	Short:   "cuemix - A tool to allow using external resources in Cue, such as Helm",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
