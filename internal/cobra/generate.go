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
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the project layout",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		// initialize configuration if it does not exist
		if !cfs.Exists(dest) {
			initializeCmd.Run(cmd, args) // will fatal if initialization fails
		}

		// validate configuration
		if err := craft.Validate(dest); err != nil {
			logger.Fatal(err)
		}

		// read configuration
		var config craft.Configuration
		if err := craft.Read(dest, &config); err != nil {
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
		config, err := generate.Run(ctx, config, options...)
		if err != nil {
			logger.Fatal(err)
		}

		// save configuration again in case it was modified during generation
		if err := craft.Write(dest, config); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
