package cobra

import (
	"net/http"
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

func gen(generators ...engine.Generator[craft.Config]) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		dest := filepath.Join(wd, craft.File)
		logger.Infof("generating layout in %s", wd)

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
		engine.SetLogger(logger)
		parsers := []engine.Parser[craft.Config]{
			generate.ParserGit,
			generate.ParserGolang,
			generate.ParserNode,
			// misc parsers
			generate.ParserShell,
			generate.ParserTmpl,
			// must be kept last since it marshals config and merges it with chart overrides
			generate.ParserHelm,
		}
		config, err := engine.Generate(ctx, wd, config, parsers, generators)
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
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate project layout",
	Run: gen(
		generate.GeneratorGitignore(http.DefaultClient),                                                                               // gitignore
		generate.GeneratorLicense(http.DefaultClient),                                                                                 // license
		engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.Dependabot(), templates.Renovate())),                        // bot
		engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.CodeCov(), templates.Sonar())),                              // coverage
		engine.GeneratorTemplates(templates.FS(), templates.Docker()),                                                                 // docker
		engine.GeneratorTemplates(templates.FS(), templates.Golang()),                                                                 // golang
		engine.GeneratorTemplates(templates.FS(), templates.Misc()),                                                                   // misc
		engine.GeneratorTemplates(templates.FS(), templates.Makefile()),                                                               // makefile
		engine.GeneratorTemplates(templates.FS(), templates.Chart()),                                                                  // chart
		engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.GitHub(), templates.GitLab(), templates.SemanticRelease())), // ci
	),
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
