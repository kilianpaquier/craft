package cobra

import (
	"os"
	"path/filepath"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the project layout",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		config, err := func() (craft.Configuration, error) {
			// initialize configuration if it does not exist
			if !cfs.Exists(dest) {
				return initialize.Run(ctx)
			}

			// validate configuration
			if err := craft.Validate(dest); err != nil {
				return craft.Configuration{}, err
			}

			// read configuration
			var config craft.Configuration
			err := craft.Read(dest, &config)
			return config, err
		}()
		if err != nil {
			logger.Fatal(err)
		}

		// run generation
		options := []generate.RunOption{
			generate.WithDestination(destdir),
			generate.WithHandlers(handler.Defaults()...),
			generate.WithLogger(logger),
			generate.WithParsers(parser.Defaults()...),
			generate.WithTemplates(generate.TmplDir, generate.FS()),
		}
		if config, err = generate.Run(ctx, config, options...); err != nil {
			logger.Fatal(err)
		}

		// save configuration
		if err := craft.Write(dest, config); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
