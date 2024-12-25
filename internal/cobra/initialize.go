package cobra

import (
	"os"
	"path/filepath"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/configuration"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
	"github.com/kilianpaquier/craft/pkg/initialize/group"
)

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project layout",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		if cfs.Exists(dest) {
			logger.Info("project already initialized")
			return
		}

		config, err := initialize.Run(ctx, initialize.WithFormGroups(group.Maintainer, group.Chart))
		if err != nil {
			logger.Fatal(err)
		}

		if err := configuration.WriteYAML(dest, config, craft.EncodeOpts()...); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
}
