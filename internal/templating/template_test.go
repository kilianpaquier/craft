package templating_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/templating"
)

func TestExecute(t *testing.T) {
	t.Run("success_dest_exists", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		src := filepath.Join(tmp, "template.txt")
		dest := filepath.Join(tmp, "template-result.txt")
		_, err := os.Create(dest)
		require.NoError(t, err)

		err = os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)
		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New(path.Base(src)).
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
		assert.Equal(t, string(bytes), "hey ! A name")
	})

	t.Run("success_shell", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		src := filepath.Join(tmp, "template.txt")
		dest := filepath.Join(tmp, "template-result.sh")

		err := os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)
		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New(path.Base(src)).
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
		assert.Equal(t, info.Mode(), filesystem.RwxRxRxRx)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, string(bytes), "hey ! A name")
	})
}
