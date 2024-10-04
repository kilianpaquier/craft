package cobra

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

var (
	force    []string
	forceAll bool

	generateCmd = &cobra.Command{
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
			generate.SetLogger(_log)

			// run generation
			options := []generate.RunOption{
				generate.WithDelimiters("<<", ">>"),
				generate.WithDestination(destdir),
				generate.WithForce(force...),
				generate.WithForceAll(forceAll),
				generate.WithTemplates("templates", generate.FS()),
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
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringSliceVarP(
		&force, "force", "f", []string{},
		"force regenerating a list of templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)")

	generateCmd.Flags().BoolVar(
		&forceAll, "force-all", false,
		"force regenerating all templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)")
}
