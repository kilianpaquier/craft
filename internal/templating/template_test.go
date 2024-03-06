package templating_test

import (
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
	t.Run("success_dest_exists", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		src := filepath.Join(tmp, "template.txt")
		dest := filepath.Join(tmp, "template-result.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		err = os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)
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
		src := filepath.Join(tmp, "template.txt")
		dest := filepath.Join(tmp, "template-result.sh")

		err := os.WriteFile(src, []byte("{{ .name }}"), filesystem.RwRR)
		require.NoError(t, err)
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
