package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "callbreak-go",
		Short: "A golang implementation of callbreak with Nepali Rules",
	}
)

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
