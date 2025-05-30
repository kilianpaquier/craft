package generate_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	compare "github.com/kilianpaquier/compare/pkg"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func TestHelmTemplate(t *testing.T) {
	ctx := t.Context()
	testdir := filepath.Join("..", "..", "testdata", t.Name())
	generate := engine.GeneratorTemplates(templates.FS(), templates.Chart()) // chart generation

	// read all tests (simpler in case new test cases would be added)
	tests, err := os.ReadDir(testdir)
	require.NoError(t, err)

	// run all tests
	for _, test := range tests {
		if !test.IsDir() {
			continue
		}

		t.Run(test.Name(), func(t *testing.T) {
			// Arrange
			assertdir := filepath.Join(testdir, test.Name())
			expected := filepath.Join(assertdir, "manifest.yaml")
			actual := filepath.Join(t.TempDir(), "manifest.yaml")

			// generate chart files
			destdir := t.TempDir()
			require.NoError(t, generate(ctx, destdir, craft.Config{
				Languages: map[string]any{"helm": map[string]any{"projectName": "craft"}},
			}))
			chartdir := filepath.Join(destdir, "chart")

			// remove default values since we use custom ones
			require.NoError(t, os.Remove(filepath.Join(chartdir, "values.yaml")))

			// copy chart additional inputs for given test
			chartinput := filepath.Join(assertdir, "chart")
			if files.Exists(chartinput) {
				require.NoError(t, os.CopyFS(chartdir, os.DirFS(chartinput)))
			}

			// Act
			manifest, err := template(ctx, chartdir, filepath.Join(assertdir, "values.yaml"))
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(actual, []byte(manifest), files.RwRR))

			// Assert
			assert.NoError(t, compare.Files(expected, actual))
		})
	}
}

// template runs a dry run of an helm install and returns the computed manifest.
func template(ctx context.Context, chartdir, valuesFile string) (string, error) {
	// load chart
	chart, err := loader.LoadDir(chartdir)
	if err != nil {
		return "", fmt.Errorf("load chart dir: %w", err)
	}

	// load values
	rawValues, err := os.ReadFile(valuesFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("read file: %w", err)
	}
	var values map[string]any
	if err := yaml.Unmarshal(rawValues, &values); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	client := action.NewInstall(&action.Configuration{})
	client.ClientOnly = true
	client.DisableHooks = true
	client.DryRun = true
	client.IncludeCRDs = false
	client.Namespace = chart.Name()
	client.ReleaseName = chart.Name()

	// run install and retrieve resulting manifest
	release, err := client.RunWithContext(ctx, chart, values)
	if err != nil {
		return "", fmt.Errorf("template chart: %w", err)
	}
	return release.Manifest, nil
}
