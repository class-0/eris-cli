package commands

import (
	"fmt"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Display Eris's Platform Version.",
	Long:  `Display the versions of what your platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Println("Eris CLI Version: " + VERSION)
	},
}
