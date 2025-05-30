package generate_test

import (
	"os"
	"testing"

	"github.com/charmbracelet/log"

	"github.com/kilianpaquier/craft/pkg/engine"
)

func TestMain(m *testing.M) {
	engine.SetLogger(log.NewWithOptions(os.Stderr, log.Options{
		CallerFormatter: log.ShortCallerFormatter,
		Level:           log.WarnLevel,
		ReportCaller:    true,
	}))
	os.Exit(m.Run())
}
