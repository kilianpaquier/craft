package cobra

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
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

		config, err := initialize.Run(ctx, destdir)
		if err != nil && !errors.Is(err, initialize.ErrAlreadyInitialized) {
			fatal(ctx, err)
		}
		config.EnsureDefaults()

		// validate craft struct
		if err := validator.New().Struct(config); err != nil {
			fatal(ctx, err)
		}

		// run generation
		options := []generate.RunOption{
			generate.WithDestination(destdir),
			generate.WithHandlers(handler.Defaults()...),
			generate.WithLogger(log),
			generate.WithParsers(parser.Defaults()...),
			generate.WithTemplates("_templates", generate.FS()),
		}
		config, err = generate.Run(ctx, config, options...)
		if err != nil {
			fatal(ctx, err)
		}

		// save craft configuration
		if err := craft.Write(destdir, config); err != nil {
			fatal(ctx, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
