package engine_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestApplyTemplate(t *testing.T) {
	mocklog := func(t *testing.T, l engine.Logger) {
		t.Helper()

		initial := engine.GetLogger()
		engine.SetLogger(l)
		t.Cleanup(func() { engine.SetLogger(initial) })
	}

	t.Run("error_missing_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		// Act
		err := engine.ApplyTemplate(os.DirFS(destdir), destdir, engine.Template[testconfig]{}, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "localize path")
	})

	t.Run("error_read_template_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "dir"}
		require.NoError(t, os.Mkdir(filepath.Join(destdir, template.Out), files.RwxRxRxRx))

		// Act
		err := engine.ApplyTemplate(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "should generate")
	})

	t.Run("error_template_invalid_globs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "file.txt"}

		buf := strings.Builder{}
		logger := engine.NewTestLogger(&buf)
		mocklog(t, logger)

		// Act
		err := engine.ApplyTemplate(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, buf.String(), fmt.Sprintf("empty template 'globs', skipping '%s' generation", template.Out))
	})

	t.Run("error_parse_template_globs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Globs: []string{"invalid.txt"},
			Out:   "file.txt",
		}

		// Act
		err := engine.ApplyTemplate(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "parse template file(s)")
	})

	t.Run("success_template_already_exists", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "file.txt"}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Out), []byte("some not empty file"), files.RwRR))

		buf := strings.Builder{}
		logger := engine.NewTestLogger(&buf)
		mocklog(t, logger)

		// Act
		err := engine.ApplyTemplate(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, buf.String(), fmt.Sprintf("not generating '%s' since it already exists", template.Out))
	})
}

func TestApplyPatches(t *testing.T) {
	t.Run("error_missing_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, engine.Template[testconfig]{}, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "localize path")
	})

	t.Run("error_missing_template_patch", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "parse template patch")
	})

	t.Run("error_template_patch", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Patches[0]), []byte("{{ .invalid }}"), files.RwRR))

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "template patch execution")
	})

	t.Run("error_invalid_patch_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Patches[0]), []byte(`
diff --git a/file.txt b/file.txt
index 332d5ce..39af8aa 100644
--- a/file.txt
+++ b/file.txt
@@ -1,0 +1,2 @@
+value`), files.RwRR))

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "parse git patch")
	})

	t.Run("error_apply_patch", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Patches[0]), []byte(`
diff --git a/file.txt b/file.txt
index 332d5ce..39af8aa 100644
--- a/file.txt
+++ b/file.txt
@@ -2,0 +2,1 @@
+value`), files.RwRR))

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "apply diff number")
	})

	t.Run("success_create", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Patches[0]), []byte(`
diff --git a/file.txt b/file.txt
index 332d5ce..39af8aa 100644
--- a/file.txt
+++ b/file.txt
@@ -1,0 +1,1 @@
+value`), files.RwRR))

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(filepath.Join(destdir, template.Out))
		require.NoError(t, err)
		assert.Equal(t, "value", string(content))
	})

	t.Run("success_update", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Out:     "file.txt",
			Patches: []string{"file.patch"},
		}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Out), []byte("some not empty file"), files.RwRR))
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Patches[0]), []byte(`
diff --git a/file.txt b/file.txt
index 332d5ce..39af8aa 100644
--- a/file.txt
+++ b/file.txt
@@ -1 +1 @@
-some not empty file
\ No newline at end of file
+some replaced value in non empty file
\ No newline at end of file`), files.RwRR))

		// Act
		err := engine.ApplyPatches(os.DirFS(destdir), destdir, template, testconfig{})

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(filepath.Join(destdir, template.Out))
		require.NoError(t, err)
		assert.Equal(t, "some replaced value in non empty file", string(content))
	})
}

func TestExecuteTemplate(t *testing.T) {
	t.Run("error_mkdir", func(t *testing.T) {
		// Arrange
		dir := filepath.Join(t.TempDir(), "dir")
		require.NoError(t, os.Mkdir(dir, files.RwxRxRxRx))

		// create empty file (at midlevel) to ensure os.MkdirAll fails
		dest := filepath.Join(dir, "file.txt", "file.txt")
		file, err := os.Create(filepath.Dir(dest))
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = engine.ExecuteTemplate(nil, nil, dest)

		// Assert
		assert.ErrorContains(t, err, "mkdir")
	})

	t.Run("error_execute", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		dest := filepath.Join(tmp, "template-result.txt")

		// not parsing any file with template to ensure tmpl.Execute fails
		tmpl := template.New("template.txt").Funcs(engine.FuncMap(tmp))

		// Act
		err := engine.ExecuteTemplate(tmpl, nil, dest)

		// Assert
		assert.ErrorContains(t, err, "template execution")
		assert.ErrorContains(t, err, `"template.txt" is an incomplete or empty template`)
	})

	t.Run("error_write_dir", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), files.RwRR)
		require.NoError(t, err)

		// create a file in dest to ensure WriteFile fails since it's a directory
		dest := filepath.Join(tmp, "dir")
		require.NoError(t, os.MkdirAll(filepath.Dir(dest), files.RwxRxRxRx))

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(engine.FuncMap(tmp)).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = engine.ExecuteTemplate(tmpl, data, filepath.Dir(dest))

		// Assert
		assert.ErrorContains(t, err, "open file")
	})

	t.Run("success_dest_exists", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()

		// create template file
		src := filepath.Join(tmp, "template.txt")
		err := os.WriteFile(src, []byte("{{ .name }}"), files.RwRR)
		require.NoError(t, err)

		// create dest to ensure os.Remove works
		dest := filepath.Join(tmp, "template-result.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		data := map[string]string{"name": "hey ! A name"}

		tmpl, err := template.New("template.txt").
			Funcs(engine.FuncMap(tmp)).
			ParseFiles(src)
		require.NoError(t, err)

		// Act
		err = engine.ExecuteTemplate(tmpl, data, dest)

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "hey ! A name", string(content))
	})
}
