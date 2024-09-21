package generate_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectHelm(t *testing.T) {
	ctx := context.Background()

	t.Run("no_chart_config_present", func(t *testing.T) {
		// Arrange
		metadata := generate.Metadata{Configuration: craft.Configuration{NoChart: true}}

		var buf bytes.Buffer
		log.SetOutput(&buf)
		// Act
		_, exec, err := generate.DetectHelm(ctx, clog.Std(), "", metadata)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.NotContains(t, buf.String(), fmt.Sprintf("helm chart detected, %s doesn't have no_chart key", craft.File))
	})

	t.Run("no_chart_config_absent", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Act
		_, exec, err := generate.DetectHelm(ctx, clog.Std(), "", generate.Metadata{})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Contains(t, buf.String(), fmt.Sprintf("helm chart detected, %s doesn't have no_chart key", craft.File))
	})
}

func TestGenerateHelm(t *testing.T) {
	ctx := context.Background()
	exec := generate.GenerateHelm

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, exec, "helm")

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", craft.File)
		require.NoError(t, os.MkdirAll(overrides, cfs.RwxRxRxRx))

		metadata := setup(generate.Metadata{})

		// Act
		err := generate.GenerateHelm(ctx, clog.Noop(), cfs.OS(), "", destdir, metadata, generate.ExecOpts{})

		// Assert
		assert.ErrorContains(t, err, "read helm chart overrides")
	})

	t.Run("success_empty_values", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_dependencies", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join("..", "..", "testdata", "helm", "success_dependencies")
		destdir := t.TempDir()

		require.NoError(t, os.Mkdir(filepath.Join(destdir, "chart"), cfs.RwxRxRxRx))
		err := cfs.CopyFile(filepath.Join(assertdir, "chart", ".craft"), filepath.Join(destdir, "chart", ".craft"))
		require.NoError(t, err)

		metadata := setup(generate.Metadata{})

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_resources", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Clis: map[string]struct{}{"cli-name": {}},
			Configuration: craft.Configuration{
				Docker: &craft.Docker{Port: helpers.ToPtr(uint16(5000))},
			},
			Crons:   map[string]struct{}{"cron-name": {}},
			Jobs:    map[string]struct{}{"job-name": {}},
			Workers: map[string]struct{}{"worker-name": {}},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})
}

func TestRemoveHelm(t *testing.T) {
	ctx := context.Background()

	t.Run("success_no_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")

		// Act
		err := generate.RemoveHelm(ctx, clog.Noop(), cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, chart)
	})

	t.Run("success_with_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")
		require.NoError(t, os.Mkdir(chart, cfs.RwxRxRxRx))

		// Act
		err := generate.RemoveHelm(ctx, clog.Noop(), cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, chart)
	})
}
