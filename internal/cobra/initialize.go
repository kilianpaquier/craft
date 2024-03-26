package cobra

import (
	"errors"
	"io/fs"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/initialize"
	"github.com/kilianpaquier/craft/internal/models"
)

var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize the project layout",
	PreRun: SetLogLevel,
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		log := logrus.WithContext(ctx)
		destdir, _ := os.Getwd()

		// read craft configuration
		var craft models.CraftConfig
		if err := configuration.ReadCraft(destdir, &craft); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				log.WithError(err).Fatal("failed to read craft configuration, file exists but is not readable")
			}

			// init repository if craft configuration wasn't found
			craft = initialize.Run(ctx)
			if err := configuration.WriteCraft(destdir, craft); err != nil {
				log.WithError(err).Warn("failed to write config file")
			}
		} else {
			log.Info("project already initialized, .craft file exists")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
