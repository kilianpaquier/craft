package generate_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestParserGolang(t *testing.T) {
	ctx := t.Context()

	t.Run("error_read_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(destdir, parser.FileGomod), files.RwxRxRxRx))

		// Act
		err := generate.ParserGolang(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, fmt.Sprintf("read '%s'", parser.FileGomod))
	})

	t.Run("success_no_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		config := craft.Config{}

		// Act
		err := generate.ParserGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Zero(t, config)
	})

	t.Run("success_hugo", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		err := os.WriteFile(filepath.Join(destdir, parser.FileGomod), []byte(
			`module github.com/kilianpaquier/craft

			go 1.22`,
		), files.RwRR)
		require.NoError(t, err)

		hugoconfig, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		require.NoError(t, hugoconfig.Close())

		expected := craft.Config{
			Languages: map[string]any{
				"hugo": parser.HugoConfig{},
			},
		}
		config := craft.Config{}

		// Act
		err = generate.ParserGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_go_no_cmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		err := os.WriteFile(filepath.Join(destdir, parser.FileGomod), []byte(
			`module github.com/kilianpaquier/craft

			go 1.22`,
		), files.RwRR)
		require.NoError(t, err)

		expected := craft.Config{
			Languages: map[string]any{
				"go": parser.Gomod{
					Module: "github.com/kilianpaquier/craft",
					Go:     "1.22",
				},
			},
		}
		config := craft.Config{}

		// Act
		err = generate.ParserGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_go_cmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		err := os.WriteFile(filepath.Join(destdir, parser.FileGomod), []byte(
			`module github.com/kilianpaquier/craft

			go 1.22`,
		), files.RwRR)
		require.NoError(t, err)

		cmd := filepath.Join(destdir, parser.FolderCMD)
		require.NoError(t, os.Mkdir(cmd, files.RwxRxRxRx))
		cli := filepath.Join(cmd, "name")
		require.NoError(t, os.Mkdir(cli, files.RwxRxRxRx))
		main, err := os.Create(filepath.Join(cli, "main.go"))
		require.NoError(t, err)
		require.NoError(t, main.Close())

		expected := craft.Config{
			Executables: parser.Executables{
				Clis: map[string]struct{}{"name": {}},
			},
			Languages: map[string]any{
				"go": parser.Gomod{
					Module: "github.com/kilianpaquier/craft",
					Go:     "1.22",
				},
			},
		}
		config := craft.Config{}

		// Act
		err = generate.ParserGolang(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
