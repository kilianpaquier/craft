package cobra

import (
	"path/filepath"

	"github.com/spf13/cobra"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/initialize"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new craft project",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		dest := filepath.Join(wd, craft.File)

		if files.Exists(dest) {
			logger.Info("project already initialized")
			return
		}

		config, err := engine.Initialize(ctx, engine.WithFormGroups(initialize.Maintainer, initialize.License))
		if err != nil {
			logger.Fatal(err)
		}

		if err := files.WriteYAML(dest, config, craft.EncodeOpts()...); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
}
