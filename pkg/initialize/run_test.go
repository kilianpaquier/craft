package initialize_test

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	t.Run("error_already_initialized", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		file, err := os.Create(filepath.Join(destdir, craft.File))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, file.Close()) })

		// Act
		config, err := initialize.Run(ctx, destdir)

		// Assert
		assert.ErrorIs(t, err, initialize.ErrAlreadyInitialized)
		assert.Equal(t, craft.Configuration{}, config)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(destdir, craft.File), cfs.RwxRxRxRx))

		// Act
		config, err := initialize.Run(ctx, destdir)

		// Assert
		assert.ErrorContains(t, err, "exists but is not readable")
		assert.Equal(t, craft.Configuration{}, config)
	})

	t.Run("success_custom_input", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{License: lo.ToPtr("mit")}

		f := func(_ clog.Logger, config craft.Configuration, ask initialize.Ask) craft.Configuration {
			config.License = ask("Which license would you like to use ?")
			return config
		}

		inputs := []string{"mit"}
		reader := strings.NewReader(strings.Join(inputs, "\n"))

		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Act
		config, err := initialize.Run(ctx, "",
			initialize.WithLogger(clog.Std()),
			initialize.WithReader(reader),
			initialize.WithInputReaders(f))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
		assert.Contains(t, buf.String(), "Which license would you like to use ?") // just verify that ask logs as it would be expected
	})

	t.Run("success_minimal_inputs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		expected := craft.Configuration{Maintainers: []craft.Maintainer{{Name: "name"}}}

		inputs := []string{
			"name", // maintainer name
			"",     // no maintainer email
			"",     // no maintainer url
			"",     // chart generation
		}
		reader := strings.NewReader(strings.Join(inputs, "\n"))

		// Act
		config, err := initialize.Run(ctx, destdir, initialize.WithReader(reader))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_no_chart", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "name"}},
			NoChart:     true,
		}

		inputs := []string{
			"name", // maintainer name
			"",     // no maintainer email
			"",     // no maintainer url
			"f",    // chart generation
		}
		reader := strings.NewReader(strings.Join(inputs, "\n"))

		// Act
		config, err := initialize.Run(ctx, "", initialize.WithReader(reader))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_all_inputs", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{Maintainers: []craft.Maintainer{
			{Name: "name", Email: lo.ToPtr("email"), URL: lo.ToPtr("url")},
		}}

		inputs := []string{
			"name",  // maintainer name
			"email", // email
			"url",   // url
			"t",     // chart generation
		}
		reader := strings.NewReader(strings.Join(inputs, "\n"))

		// Act
		config, err := initialize.Run(ctx, "",
			initialize.WithReader(reader))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_retryable_inputs", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{Maintainers: []craft.Maintainer{{Name: "name"}}}

		inputs := []string{
			"",                    // empty maintainer first
			"name",                // maintainer name
			"",                    // no maintainer email
			"",                    // no maintainer url
			"invalid chart value", // invalid chart boolean
			"t",                   // chart generation
		}
		reader := strings.NewReader(strings.Join(inputs, "\n"))

		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Act
		config, err := initialize.Run(ctx, "",
			initialize.WithLogger(clog.Std()),
			initialize.WithReader(reader))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
		assert.Contains(t, buf.String(), "invalid chart answer 'invalid chart value', must be a boolean")
	})
}
