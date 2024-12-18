package cobra

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	logger    = log.New(os.Stderr)
	logLevel  = "info"
	logFormat = "text"
	rootCmd   = &cobra.Command{
		Use:               "craft",
		SilenceErrors:     true, // don't print errors with cobra, let logger.Fatal handle them
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return preRun() },
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set logging level")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", `set logging format (either "text" or "json")`)

	_ = preRun() // ensure logging is correctly configured with default values even when a bad input flag is given
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func preRun() error {
	switch logFormat {
	case "text":
		logger.SetFormatter(log.TextFormatter)
	case "json":
		logger.SetFormatter(log.JSONFormatter)
	default:
		return errors.New(`invalid --log-format argument, must be either "json" or "text"`)
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		level = log.InfoLevel
	}
	logger.SetLevel(level)
	return nil
}
