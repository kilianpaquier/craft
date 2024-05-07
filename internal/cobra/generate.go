package cobra

import (
	"errors"
	"io/fs"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/initialize"
	"github.com/kilianpaquier/craft/internal/models"
)

var (
	generateOpts = models.GenerateOptions{
		EndDelim:   ">>",
		StartDelim: "<<",
	}

	generateCmd = &cobra.Command{
		Use:    "generate",
		Short:  "Generate the project layout",
		PreRun: SetLogLevel,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			log := logrus.WithContext(ctx)

			// init destdir for file copying and templating
			generateOpts.DestinationDir, _ = os.Getwd()

			// read config configuration
			var config models.CraftConfig
			if err := configuration.ReadCraft(generateOpts.DestinationDir, &config); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					log.WithError(err).Fatal("failed to read craft configuration, file exists but is not readable")
				}

				// init repository if craft configuration wasn't found
				config = initialize.Run(ctx)
			}
			config = configuration.EnsureDefaults(config)

			// validate craft struct
			if err := validator.New().Struct(config); err != nil {
				log.WithError(err).Fatal("failed to validate craft configuration")
			}

			// create craft runner
			runner, err := generate.NewRunner(ctx, config, generateOpts)
			if err != nil {
				log.WithError(err).Fatal("failed to create craft executor")
			}

			// generate all files
			log.Infof("start craft generation in %s", generateOpts.DestinationDir)
			config, err = runner.Run(ctx)
			if err != nil {
				log.WithError(err).Fatal("failed to execute craft generation")
			}
			if err := configuration.WriteCraft(generateOpts.DestinationDir, config); err != nil {
				log.WithError(err).Warn("failed to write config file")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringSliceVarP(
		&generateOpts.Force, "force", "f", []string{},
		"force regenerating a list of templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)")
	generateCmd.Flags().BoolVar(
		&generateOpts.ForceAll, "force-all", false,
		"force regenerating all templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)")
}
