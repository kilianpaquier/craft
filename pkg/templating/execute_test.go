package templating_test

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/templating"
)

func TestExecute(t *testing.T) {
	t.Run("error_mkdir", func(t *testing.T) {
		// Arrange
		dir := filepath.Join(t.TempDir(), "dir")
		require.NoError(t, os.Mkdir(dir, cfs.RwxRxRxRx))

		// create empty file (at midlevel) to ensure os.MkdirAll fails
		dest := filepath.Join(dir, "file.txt", "file.txt")
		file, err := os.Create(filepath.Dir(dest))
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = templating.Execute(nil, nil, dest)

		// Assert
		assert.ErrorContains(t, err, "create directory")
	})

	t.Run("error_execute", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		dest := filepath.Join(tmp, "template-result.txt")

		// not parsing any file with template to ensure tmpl.Execute fails
		tmpl := template.New("template.txt").Funcs(templating.FuncMap())

		// Act
		err := templating.Execute(tmpl, nil, dest)

		// Assert
		assert.ErrorContains(t, err, "template execution")
		assert.ErrorContains(t, err, `"template.txt" is an incomplete or empty template`)
	})

	t.Run("error_write_dir", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), cfs.RwRR)
		require.NoError(t, err)

		// create a file in dest to ensure WriteFile fails since it's a directory
		dest := filepath.Join(tmp, "dir")
		require.NoError(t, os.MkdirAll(filepath.Dir(dest), cfs.RwxRxRxRx))

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(templating.FuncMap()).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = templating.Execute(tmpl, data, filepath.Dir(dest))

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success_dest_exists", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), cfs.RwRR)
		require.NoError(t, err)

		// create dest to ensure os.Remove works
		dest := filepath.Join(tmp, "template-result.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(templating.FuncMap()).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = templating.Execute(tmpl, data, dest)

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "hey ! A name", string(bytes))
	})
}
