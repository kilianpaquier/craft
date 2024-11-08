package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectGolang(t *testing.T) {
	ctx := context.Background()

	t.Run("no_gomod", func(t *testing.T) {
		// Arrange
		config := generate.Metadata{}

		// Act
		exec, err := generate.DetectGolang(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, exec)
		assert.Zero(t, config)
	})

	t.Run("invalid_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, craft.Gomod)
		err := os.WriteFile(gomod, []byte("an invalid go.mod file"), cfs.RwRR)
		require.NoError(t, err)

		config := generate.Metadata{}

		// Act
		exec, err := generate.DetectGolang(ctx, destdir, &config)

		// Assert
		assert.ErrorContains(t, err, "read go.mod")
		assert.Empty(t, exec)
		assert.Zero(t, config)
	})

	t.Run("missing_gomod_statements", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, craft.Gomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		config := generate.Metadata{}

		// Act
		exec, err := generate.DetectGolang(ctx, destdir, &config)

		// Assert
		assert.ErrorContains(t, err, "read go.mod")
		assert.ErrorContains(t, err, "invalid go.mod, module statement is missing")
		assert.ErrorContains(t, err, "invalid go.mod, go statement is missing")
		assert.Empty(t, exec)
		assert.Zero(t, config)
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

		config := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Configuration: craft.Configuration{Platform: craft.GitHub},
			Languages: map[string]any{
				"golang": generate.Gomod{
					LangVersion: "1.22",
					Platform:    craft.GitHub,
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
		exec, err := generate.DetectGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, config)
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

		config := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{Options: []string{craft.CodeCov, craft.CodeQL, craft.Sonar}}},
			Languages:     map[string]any{},
		}
		expected := generate.Metadata{
			Configuration: craft.Configuration{
				CI:       &craft.CI{Options: []string{}},
				Platform: craft.GitHub,
			},
			Languages:   map[string]any{"hugo": nil},
			ProjectHost: "github.com",
			ProjectName: "craft",
			ProjectPath: "kilianpaquier/craft",
		}

		// Act
		exec, err := generate.DetectGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, config)
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

		config := generate.Metadata{
			Clis:      map[string]struct{}{},
			Crons:     map[string]struct{}{},
			Jobs:      map[string]struct{}{},
			Languages: map[string]any{},
			Workers:   map[string]struct{}{},
		}
		expected := generate.Metadata{
			Binaries:      4,
			Clis:          map[string]struct{}{"cli-name": {}},
			Configuration: craft.Configuration{Platform: craft.GitHub},
			Crons:         map[string]struct{}{"cron-name": {}},
			Jobs:          map[string]struct{}{"job-name": {}},
			Languages: map[string]any{
				"golang": generate.Gomod{
					LangVersion: "1.22.2",
					Platform:    craft.GitHub,
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
		exec, err := generate.DetectGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, config)
	})
}
