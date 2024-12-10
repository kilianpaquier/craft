package handler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGenerated(t *testing.T) {
	t.Run("generated_doesnt_exist", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "invalid.txt")

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("not_generated_file", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("not generated"), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.False(t, generated)
	})

	t.Run("generated_folder", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "folder")
		require.NoError(t, os.Mkdir(dest, cfs.RwxRxRxRx))

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_no_lines", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, file.Close()) })

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_first_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("// Code generated by craft; DO NOT EDIT."), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_md_comment", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("<!-- Code generated by craft; DO NOT EDIT. -->"), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_second_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("\n# Code generated by craft; DO NOT EDIT."), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_json", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte(`{
			"//": "Code generated by craft; DO NOT EDIT.",
		}`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := handler.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})
}