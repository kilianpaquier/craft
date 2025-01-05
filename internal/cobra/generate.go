package cobra

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"

	schemas "github.com/kilianpaquier/craft/.schemas"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func gen(opts ...engine.GenerateOption[craft.Config]) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		// initialize configuration if it does not exist
		if !files.Exists(dest) {
			initializeCmd.Run(cmd, args) // will fatal if initialization fails
		}

		// validate configuration
		if err := files.Validate(
			func(out any) error { return files.ReadJSON(schemas.Craft, out, schemas.ReadFile) }, // read schema
			func(out any) error { return files.ReadYAML(dest, out, os.ReadFile) },               // read configuration
		); err != nil {
			logger.Fatal(err)
		}

		// read configuration
		var config craft.Config
		if err := files.ReadYAML(dest, &config, os.ReadFile); err != nil {
			logger.Fatal(err)
		}
		config.EnsureDefaults()

		// run generation
		options := slices.Concat(
			[]engine.GenerateOption[craft.Config]{
				engine.WithDestination[craft.Config](destdir),
				engine.WithLogger[craft.Config](logger),
			},
			opts,
		)
		config, err := engine.Generate(ctx, config, options...)
		if err != nil {
			logger.Fatal(err)
		}

		// save configuration again in case it was modified during generation
		if err := files.WriteYAML(dest, config, craft.EncodeOpts()...); err != nil {
			logger.Fatal(err)
		}
	}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate project layout",
	Run: gen(
		engine.WithParsers(generate.ParserGit, generate.ParserLicense, generate.ParserGolang, generate.ParserNode, generate.ParserHelm),
		engine.WithTemplates(templates.FS(), templates.All()),
	),
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
