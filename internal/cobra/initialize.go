package cobra

import (
	"errors"
	"os"

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

		config, err := initialize.Run(ctx, destdir, initialize.WithLogger(_log))
		if err != nil {
			if !errors.Is(err, initialize.ErrAlreadyInitialized) {
				fatal(ctx, err)
			}
			_log.Info("project already initialized")
			return
		}

		if err := craft.Write(destdir, config); err != nil {
			fatal(ctx, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
}
