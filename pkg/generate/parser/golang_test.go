package parser_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestGolang(t *testing.T) {
	ctx := context.Background()

	t.Run("no_gomod", func(t *testing.T) {
		// Arrange
		config := craft.Config{}

		// Act
		err := parser.Golang(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Zero(t, config)
	})

	t.Run("invalid_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, parser.FileGomod)
		err := os.WriteFile(gomod, []byte("an invalid go.mod file"), cfs.RwRR)
		require.NoError(t, err)

		config := craft.Config{}

		// Act
		err = parser.Golang(ctx, destdir, &config)

		// Assert
		assert.ErrorContains(t, err, "read go.mod")
		assert.Zero(t, config)
	})

	t.Run("missing_gomod_statements", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, parser.FileGomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		config := craft.Config{}

		// Act
		err = parser.Golang(ctx, destdir, &config)

		// Assert
		assert.ErrorIs(t, err, parser.ErrMissingGoStatement)
		assert.ErrorIs(t, err, parser.ErrMissingModuleStatement)
		assert.Zero(t, config)
	})

	t.Run("detected_no_gocmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, parser.FileGomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
			
			go 1.22`,
		), cfs.RwRR)
		require.NoError(t, err)

		config := craft.Config{ConfigFiles: craft.ConfigFiles{Languages: map[string]any{}}}
		expected := craft.Config{
			ConfigFiles: craft.ConfigFiles{
				Languages: map[string]any{
					"golang": parser.Gomod{
						LangVersion: "1.22",
						ModulePath:  "github.com/kilianpaquier/craft",
					},
				},
			},
			ConfigVCS: craft.ConfigVCS{
				ProjectHost: "github.com",
				ProjectName: "craft",
				ProjectPath: "kilianpaquier/craft",
				Platform:    craft.GitHub,
			},
		}

		// Act
		err = parser.Golang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("detected_hugo_override", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, parser.FileGomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
	
			go 1.22`,
		), cfs.RwRR)
		require.NoError(t, err)

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		config := craft.Config{
			ConfigFiles: craft.ConfigFiles{Languages: map[string]any{}},
		}
		expected := craft.Config{
			ConfigFiles: craft.ConfigFiles{Languages: map[string]any{"hugo": nil}},
			ConfigVCS: craft.ConfigVCS{
				ProjectHost: "github.com",
				ProjectName: "craft",
				ProjectPath: "kilianpaquier/craft",
				Platform:    craft.GitHub,
			},
		}

		// Act
		err = parser.Golang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("detected_all_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, parser.FileGomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft/v2
			
			go 1.22.1
			
			toolchain go1.22.2`,
		), cfs.RwRR)
		require.NoError(t, err)

		gocmd := filepath.Join(destdir, parser.FolderCMD)
		for _, dir := range []string{
			gocmd,
			filepath.Join(gocmd, "cli-name"),
			filepath.Join(gocmd, "cron-name"),
			filepath.Join(gocmd, "job-name"),
			filepath.Join(gocmd, "worker-name"),
		} {
			require.NoError(t, os.Mkdir(dir, cfs.RwxRxRxRx))
		}

		config := craft.Config{
			ConfigFiles: craft.ConfigFiles{
				Languages: map[string]any{},
				Clis:      map[string]struct{}{},
				Crons:     map[string]struct{}{},
				Jobs:      map[string]struct{}{},
				Workers:   map[string]struct{}{},
			},
		}
		expected := craft.Config{
			ConfigFiles: craft.ConfigFiles{
				Clis:  map[string]struct{}{"cli-name": {}},
				Crons: map[string]struct{}{"cron-name": {}},
				Jobs:  map[string]struct{}{"job-name": {}},
				Languages: map[string]any{
					"golang": parser.Gomod{
						LangVersion: "1.22.2",
						ModulePath:  "github.com/kilianpaquier/craft/v2",
					},
				},
				Workers: map[string]struct{}{"worker-name": {}},
			},
			ConfigVCS: craft.ConfigVCS{
				ProjectHost: "github.com",
				ProjectName: "craft",
				ProjectPath: "kilianpaquier/craft",
				Platform:    craft.GitHub,
			},
		}

		// Act
		err = parser.Golang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
