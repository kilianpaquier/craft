package cobra

import (
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/internal/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current craft version",
	Run:   func(_ *cobra.Command, _ []string) { logger.Info(build.GetInfo()) },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
