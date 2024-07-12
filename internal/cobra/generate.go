package cobra

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

var (
	force    []string
	forceAll bool

	generateCmd = &cobra.Command{
		Use:    "generate",
		Short:  "Generate the project layout",
		PreRun: SetLogLevel,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			log := logrus.WithContext(ctx)
			destdir, _ := os.Getwd()

			config, err := initialize.Run(ctx, destdir, initialize.WithLogger(log))
			if err != nil && !errors.Is(err, initialize.ErrAlreadyInitialized) {
				return fmt.Errorf("initialize project: %w", err)
			}
			config = config.EnsureDefaults()

			// validate craft struct
			if err := validator.New().Struct(config); err != nil {
				return fmt.Errorf("craft config validation: %w", err)
			}

			// run generation
			config, err = generate.Run(ctx, config,
				generate.WithDelimiters("<<", ">>"),
				generate.WithDestination(destdir),
				generate.WithDetects(generate.Detects()...),
				generate.WithMetaHandlers(generate.MetaHandlers()...),
				generate.WithForce(force...),
				generate.WithForceAll(forceAll),
				generate.WithLogger(log),
				generate.WithTemplates("templates", generate.FS()),
			)
			if err != nil {
				return fmt.Errorf("run generation: %w", err)
			}

			// save craft configuration
			if err := craft.Write(destdir, config); err != nil {
				return fmt.Errorf("write craft: %w", err)
			}
			return nil
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
