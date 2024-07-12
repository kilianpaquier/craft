package cobra

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize the project layout",
	PreRun: SetLogLevel,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		log := logrus.WithContext(ctx)
		destdir, _ := os.Getwd()

		config, err := initialize.Run(ctx, destdir, initialize.WithLogger(log))
		if err != nil {
			if !errors.Is(err, initialize.ErrAlreadyInitialized) {
				return fmt.Errorf("initialize project: %w", err)
			}
			return nil
		}

		if err := craft.Write(destdir, config); err != nil {
			return fmt.Errorf("write craft: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
