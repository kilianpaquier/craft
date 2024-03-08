package templating_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/templating"
)

func TestExecute(t *testing.T) {
	t.Run("error_mkdir", func(t *testing.T) {
		// Arrange
		dir := filepath.Join(t.TempDir(), "dir")
		require.NoError(t, os.Mkdir(dir, filesystem.RwxRxRxRx))

		// create empty file (at midlevel) to ensure os.MkdirAll fails
		dest := filepath.Join(dir, "file.txt", "file.txt")
		file, err := os.Create(filepath.Dir(dest))
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = templating.Execute(nil, dest, nil)

		// Assert
		assert.ErrorContains(t, err, "failed to create destination directory")
	})

	t.Run("error_execute", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		dest := filepath.Join(tmp, "template-result.txt")

		// not parsing any file with template to ensure tmpl.Execute fails
		tmpl := template.New("template.txt").
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap())

		// Act
		err := templating.Execute(tmpl, dest, nil)

		// Assert
		assert.ErrorContains(t, err, "failed to template")
		assert.ErrorContains(t, err, `"template.txt" is an incomplete or empty template`)
	})

	t.Run("error_remove", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)

		dest := filepath.Join(tmp, "dir", "path", "file.txt")
		require.NoError(t, os.MkdirAll(filepath.Dir(dest), filesystem.RwxRxRxRx))

		// create a file in dest to ensure os.Remove in Execute function fails
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = templating.Execute(tmpl, filepath.Dir(dest), data)

		// Assert
		assert.ErrorContains(t, err, fmt.Sprintf("failed to remove %s before rewritting it", filepath.Dir(dest)))
	})

	t.Run("success_dest_exists", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)

		// create dest to ensure os.Remove works
		dest := filepath.Join(tmp, "template-result.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = templating.Execute(tmpl, dest, data)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "hey ! A name", string(bytes))
	})

	t.Run("success_shell", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)

		dest := filepath.Join(tmp, "template-result.sh")

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = templating.Execute(tmpl, dest, data)

		// Assert
		assert.NoError(t, err)
		info, err := os.Stat(dest)
		assert.NoError(t, err)
		if runtime.GOOS == "linux" {
			assert.Equal(t, info.Mode(), filesystem.RwxRxRxRx)
		}
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "hey ! A name", string(bytes))
	})
}
