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
	testfs "github.com/kilianpaquier/cli-sdk/pkg/cfs/tests"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	assertdir := filepath.Join("..", "..", "testdata", "helm")
	srcdir := "templates"

	setup := func(metadata generate.Metadata) (generate.Metadata, generate.ExecOpts) {
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"

		return metadata, generate.ExecOpts{
			FileHandlers: lo.Map(generate.MetaHandlers(), func(handler generate.MetaHandler, _ int) generate.FileHandler {
				return handler(metadata)
			}),
			EndDelim:   ">>",
			StartDelim: "<<",
			ForceAll:   true,
		}
	}

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", craft.File)
		require.NoError(t, os.MkdirAll(overrides, cfs.RwxRxRxRx))

		metadata, opts := setup(generate.Metadata{
			Configuration: craft.Configuration{
				Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			},
		})

		// Act
		err := generate.GenerateHelm(ctx, clog.Noop(), cfs.OS(), srcdir, destdir, metadata, opts)

		// Assert
		assert.ErrorContains(t, err, "read helm chart overrides")
	})

	t.Run("success_empty_values", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "empty_values")

		metadata, opts := setup(generate.Metadata{
			Configuration: craft.Configuration{
				Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			},
		})

		// Act
		err := generate.GenerateHelm(ctx, clog.Noop(), cfs.OS(), srcdir, destdir, metadata, opts)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
	})

	t.Run("success_with_dependencies", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_dependencies")

		require.NoError(t, os.Mkdir(filepath.Join(destdir, "chart"), cfs.RwxRxRxRx))
		err := cfs.CopyFile(filepath.Join(assertdir, "chart", ".craft"), filepath.Join(destdir, "chart", ".craft"))
		require.NoError(t, err)

		metadata, opts := setup(generate.Metadata{
			Configuration: craft.Configuration{
				Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			},
		})

		// Act
		err = generate.GenerateHelm(ctx, clog.Noop(), cfs.OS(), srcdir, destdir, metadata, opts)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
	})

	t.Run("success_with_resources", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_resources")

		metadata, opts := setup(generate.Metadata{
			Clis: map[string]struct{}{"cli-name": {}},
			Configuration: craft.Configuration{
				Docker:      &craft.Docker{Port: lo.ToPtr(uint16(5000))},
				Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			},
			Crons:   map[string]struct{}{"cron-name": {}},
			Jobs:    map[string]struct{}{"job-name": {}},
			Workers: map[string]struct{}{"worker-name": {}},
		})

		// Act
		err := generate.GenerateHelm(ctx, clog.Noop(), cfs.OS(), srcdir, destdir, metadata, opts)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
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
