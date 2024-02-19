package cobra

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel string

	rootCmd = &cobra.Command{
		Use:   "craft",
		Short: "Craft",
		Long: `Craft stands here to generate a similar project layout for all your projects. 
Multiple coding languages are supported and even helm chart can be generated. 
For more information please consult each command specificities.`,
	}
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "set logging level")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1) // nolint:revive
	}
}

// SetLogLevel sets the logging level for logrus.
func SetLogLevel(_ *cobra.Command, _ []string) {
	level := func() logrus.Level {
		if logLevel == "" {
			return logrus.InfoLevel
		}

		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Warnf("invalid log-level '%s'. Using 'info'", logLevel)
			return logrus.InfoLevel
		}
		return level
	}()
	if level < logrus.DebugLevel {
		log.SetOutput(io.Discard) // disable goswagger logs
	}
	logrus.SetLevel(level)
}
