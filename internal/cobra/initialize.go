package cobra

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

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

		var config craft.Configuration
		err := craft.Read(dest, &config)
		if err == nil {
			logger.Infof("project already initialized")
			return
		}
		if !errors.Is(err, fs.ErrNotExist) {
			fatal(ctx, err)
		}
		if config, err = initialize.Run(ctx); err != nil {
			fatal(ctx, err)
		}
		if err := craft.Write(dest, config); err != nil {
			fatal(ctx, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
}
