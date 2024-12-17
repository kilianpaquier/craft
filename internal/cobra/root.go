package cobra

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	logger   = log.New(os.Stderr)
	logLevel = "info"
	rootCmd  = &cobra.Command{
		Use:               "craft",
		SilenceErrors:     true, // errors are already logged by fatal function when Execute has an error
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return preRun() },
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set logging level")

	_ = preRun() // ensure logging is correctly configured with default values even when a bad input flag is given
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fatal(context.Background(), err)
	}
}

func preRun() error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		level = log.InfoLevel
	}
	logger.SetLevel(level)
	return nil
}

func fatal(_ context.Context, err error) {
	logger.Fatal(err)
}
