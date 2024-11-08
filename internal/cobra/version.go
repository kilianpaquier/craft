package cobra

import (
	"github.com/spf13/cobra"
)

var (
	// version is substituted with -ldflags
	version = "v0.0.0"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show current craft version",
		Run:   func(_ *cobra.Command, _ []string) { log.Info(version) },
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
