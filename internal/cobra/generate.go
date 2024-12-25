package cobra

import (
	"os"
	"path/filepath"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/spf13/cobra"

	schemas "github.com/kilianpaquier/craft/.schemas"
	"github.com/kilianpaquier/craft/pkg/configuration"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
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
		read := func() ([]byte, error) { return schemas.ReadFile(schemas.Craft) }
		if err := configuration.Validate(dest, read); err != nil {
			logger.Fatal(err)
		}

		// read configuration
		var config craft.Config
		if err := configuration.ReadYAML(dest, &config); err != nil {
			logger.Fatal(err)
		}

		// run generation
		options := []generate.RunOption[craft.Config]{
			generate.WithDestination[craft.Config](destdir),
			generate.WithHandlers(handler.Defaults()...),
			generate.WithLogger[craft.Config](logger),
			generate.WithParsers(parser.Defaults()...),
			generate.WithTemplates[craft.Config](generate.TmplDir, generate.FS()),
		}
		config, err := generate.Run(ctx, config, options...)
		if err != nil {
			logger.Fatal(err)
		}

		// save configuration again in case it was modified during generation
		if err := configuration.WriteYAML(dest, config, craft.EncodeOpts()...); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
