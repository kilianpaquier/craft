package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	cfs "github.com/kilianpaquier/craft/pkg/fs"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/logger"
)

func TestDetectGolang(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	t.Run("no_gomod", func(t *testing.T) {
		// Act
		output, exec := generate.DetectGolang(ctx, log, "", generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
		assert.Zero(t, output)
	})

	t.Run("invalid_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, craft.Gomod)
		err := os.WriteFile(gomod, []byte("an invalid go.mod file"), cfs.RwRR)
		require.NoError(t, err)

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		output, exec := generate.DetectGolang(ctx, log, destdir, generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
		assert.Zero(t, output)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, "parse go.mod:")
	})

	t.Run("missing_gomod_statements", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, craft.Gomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		output, exec := generate.DetectGolang(ctx, log, destdir, generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
		assert.Zero(t, output)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, "failed to parse go.mod statements")
		assert.Contains(t, logs, "invalid go.mod, module statement is missing")
		assert.Contains(t, logs, "invalid go.mod, go statement is missing")
	})

	t.Run("detected_no_gocmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, craft.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
			
			go 1.22`,
		), cfs.RwRR)
		require.NoError(t, err)

		input := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Configuration: craft.Configuration{Platform: craft.Github},
			Languages: map[string]any{
				"golang": generate.Gomod{
					LangVersion: "1.22",
					Platform:    craft.Github,
					ProjectHost: "github.com",
					ProjectName: "craft",
					ProjectPath: "kilianpaquier/craft",
				},
			},
			ProjectHost: "github.com",
			ProjectName: "craft",
			ProjectPath: "kilianpaquier/craft",
		}

		// Act
		output, exec := generate.DetectGolang(ctx, log, destdir, input)

		// Assert
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})

	t.Run("detected_hugo_override", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, craft.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
	
			go 1.22`,
		), cfs.RwRR)
		require.NoError(t, err)

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		input := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{Options: []string{craft.CodeCov, craft.CodeQL, craft.Dependabot}}},
			Languages:     map[string]any{},
		}
		expected := generate.Metadata{
			Configuration: craft.Configuration{
				CI:       &craft.CI{Options: []string{craft.Dependabot}},
				Platform: craft.Github,
			},
			Languages:   map[string]any{"hugo": nil},
			ProjectHost: "github.com",
			ProjectName: "craft",
			ProjectPath: "kilianpaquier/craft",
		}

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		output, exec := generate.DetectGolang(ctx, log, destdir, input)

		// Assert
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, "hugo detected, a hugo configuration file or hugo theme file is present")
	})

	t.Run("detected_all_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, craft.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft/v2
			
			go 1.22.1
			
			toolchain go1.22.2`,
		), cfs.RwRR)
		require.NoError(t, err)

		gocmd := filepath.Join(destdir, craft.Gocmd)
		for _, dir := range []string{
			gocmd,
			filepath.Join(gocmd, "cli-name"),
			filepath.Join(gocmd, "cron-name"),
			filepath.Join(gocmd, "job-name"),
			filepath.Join(gocmd, "worker-name"),
		} {
			require.NoError(t, os.Mkdir(dir, cfs.RwxRxRxRx))
		}

		input := generate.Metadata{
			Clis:      map[string]struct{}{},
			Crons:     map[string]struct{}{},
			Jobs:      map[string]struct{}{},
			Languages: map[string]any{},
			Workers:   map[string]struct{}{},
		}
		expected := generate.Metadata{
			Binaries:      4,
			Clis:          map[string]struct{}{"cli-name": {}},
			Configuration: craft.Configuration{Platform: craft.Github},
			Crons:         map[string]struct{}{"cron-name": {}},
			Jobs:          map[string]struct{}{"job-name": {}},
			Languages: map[string]any{
				"golang": generate.Gomod{
					LangVersion: "1.22.2",
					Platform:    craft.Github,
					ProjectHost: "github.com",
					ProjectName: "craft",
					ProjectPath: "kilianpaquier/craft",
				},
			},
			ProjectHost: "github.com",
			ProjectName: "craft",
			ProjectPath: "kilianpaquier/craft",
			Workers:     map[string]struct{}{"worker-name": {}},
		}

		// Act
		output, exec := generate.DetectGolang(ctx, log, destdir, input)

		// Assert
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})
}