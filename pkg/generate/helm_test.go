package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestGenerateHelm(t *testing.T) {
	ctx := context.Background()

	execs, err := generate.DetectHelm(ctx, "", &generate.Metadata{})
	require.NoError(t, err)
	require.Len(t, execs, 1)

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, execs[0], "helm")

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", craft.File)
		require.NoError(t, os.MkdirAll(overrides, cfs.RwxRxRxRx))

		metadata := setup(generate.Metadata{})

		// Act
		err := execs[0](ctx, cfs.OS(), "", destdir, metadata, generate.ExecOpts{})

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

	execs, err := generate.DetectHelm(ctx, "", &generate.Metadata{Configuration: craft.Configuration{NoChart: true}})
	require.NoError(t, err)
	require.Len(t, execs, 1)

	t.Run("success_no_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")

		// Act
		err := execs[0](ctx, nil, "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.NoDirExists(t, chart)
	})

	t.Run("success_with_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")
		require.NoError(t, os.Mkdir(chart, cfs.RwxRxRxRx))

		// Act
		err := execs[0](ctx, nil, "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.NoDirExists(t, chart)
	})
}
