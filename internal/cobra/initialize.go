package cobra

import (
	"os"
	"path/filepath"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project layout",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		if cfs.Exists(dest) {
			logger.Infof("project already initialized")
			return
		}

		config, err := initialize.Run(ctx)
		if err != nil {
			logger.Fatal(err)
		}
		if err := craft.Write(dest, config); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
}
