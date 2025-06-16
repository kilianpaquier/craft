package parser_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestReadGomod(t *testing.T) {
	t.Run("no_gomod", func(t *testing.T) {
		// Act
		_, err := parser.ReadGomod("")

		// Assert
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("invalid_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, parser.FileGomod)
		err := os.WriteFile(gomod, []byte("an invalid go.mod file"), files.RwRR)
		require.NoError(t, err)

		// Act
		_, err = parser.ReadGomod(destdir)

		// Assert
		assert.ErrorContains(t, err, "parse modfile")
	})

	t.Run("missing_module_statement", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, parser.FileGomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		// Act
		_, err = parser.ReadGomod(destdir)

		// Assert
		assert.ErrorIs(t, err, parser.ErrMissingModuleStatement)
	})

	t.Run("missing_go_statement", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		err := os.WriteFile(filepath.Join(destdir, parser.FileGomod), []byte("module github.com/kilianpaquier/craft"), files.RwRR)
		require.NoError(t, err)

		// Act
		_, err = parser.ReadGomod(destdir)

		// Assert
		assert.ErrorIs(t, err, parser.ErrMissingGoStatement)
	})

	t.Run("golang_detected", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		err := os.WriteFile(filepath.Join(destdir, parser.FileGomod), []byte(
			`module github.com/kilianpaquier/craft

			go 1.22`,
		), files.RwRR)
		require.NoError(t, err)

		expectedMod := parser.Gomod{
			Go:     "1.22",
			Module: "github.com/kilianpaquier/craft",
			Tools:  []string{},
		}
		expectedVCS := parser.VCS{
			ProjectHost: "github.com",
			ProjectName: "craft",
			ProjectPath: "kilianpaquier/craft",
			Platform:    parser.GitHub,
		}

		// Act
		mod, err := parser.ReadGomod(destdir)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedMod, mod)
		assert.Equal(t, expectedVCS, mod.AsVCS())
	})
}

func TestReadGocmd(t *testing.T) {
	t.Run("not_detected_no_main", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		executable := filepath.Join(destdir, parser.FolderCMD, "cli-name")
		require.NoError(t, os.MkdirAll(executable, files.RwxRxRxRx))

		// Act
		executables, err := parser.ReadGoCmd(destdir)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, parser.Executables{}, executables)
	})

	t.Run("detected_executables", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		cmd := filepath.Join(destdir, parser.FolderCMD)
		require.NoError(t, os.Mkdir(cmd, files.RwxRxRxRx))
		for _, executable := range []string{"name", "cron-name", "job-name", "worker-name"} {
			dir := filepath.Join(cmd, executable)
			require.NoError(t, os.Mkdir(dir, files.RwxRxRxRx))
			main, err := os.Create(filepath.Join(dir, "main.go"))
			require.NoError(t, err)
			require.NoError(t, main.Close())
		}

		expected := parser.Executables{
			Clis:    map[string]struct{}{"name": {}},
			Crons:   map[string]struct{}{"cron-name": {}},
			Jobs:    map[string]struct{}{"job-name": {}},
			Workers: map[string]struct{}{"worker-name": {}},
		}

		// Act
		executables, err := parser.ReadGoCmd(destdir)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, executables)
	})
}

func TestHugo(t *testing.T) {
	t.Run("detected_hugo", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		expected := parser.HugoConfig{IsTheme: false}

		// Act
		config, ok := parser.Hugo(destdir)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, expected, config)
	})

	t.Run("detected_hugo_theme", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		hugo, err := os.Create(filepath.Join(destdir, "theme.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		expected := parser.HugoConfig{IsTheme: true}

		// Act
		config, ok := parser.Hugo(destdir)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, expected, config)
	})
}
