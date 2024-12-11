package cobra

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

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

		// validate configuration
		if err := craft.Validate(dest); err != nil {
			fatal(ctx, err)
		}

		// read or initialize configuration
		var config craft.Configuration
		if err := craft.Read(dest, &config); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				fatal(ctx, err)
			}
			if config, err = initialize.Run(ctx); err != nil {
				fatal(ctx, err)
			}
		}

		// run generation
		options := []generate.RunOption{
			generate.WithDestination(destdir),
			generate.WithHandlers(handler.Defaults()...),
			generate.WithLogger(log),
			generate.WithParsers(parser.Defaults()...),
			generate.WithTemplates("_templates", generate.FS()),
		}
		config, err := generate.Run(ctx, config, options...)
		if err != nil {
			fatal(ctx, err)
		}

		// save configuration
		if err := craft.Write(dest, config); err != nil {
			fatal(ctx, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
