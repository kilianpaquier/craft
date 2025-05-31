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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yannh/kubeconform/pkg/validator"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/kilianpaquier/craft/internal/helpers"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestHelmTemplate(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	testdir := filepath.Join("..", "..", "testdata", t.Name())
	generate := engine.GeneratorTemplates(templates.FS(), templates.Chart()) // chart generation

	// read all tests (simpler in case new test cases would be added)
	tests, err := os.ReadDir(testdir)
	require.NoError(t, err)

	// parse kube version since lint doesn't take latest one
	kubeVersion, err := chartutil.ParseKubeVersion("1.33")
	require.NoError(t, err)

	// run all tests
	for _, test := range tests {
		if !test.IsDir() {
			continue
		}

		t.Run(test.Name(), func(t *testing.T) {
			t.Parallel()

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
			r, err := template(ctx, kubeVersion, chartdir, filepath.Join(assertdir, "values.yaml"))
			if err != nil {
				t.Fatal(err, helpers.FromPtr(r).Manifest)
			}
			require.NoError(t, os.WriteFile(actual, []byte(r.Manifest), files.RwRR))

			t.Run("lint", func(t *testing.T) {
				t.Parallel()

				// Act
				result, err := lint(kubeVersion, chartdir, filepath.Join(assertdir, "values.yaml"))
				require.NoError(t, err)

				// Assert
				assert.False(t, action.HasWarningsOrErrors(result), result.Messages)
			})

			t.Run("template", func(t *testing.T) {
				t.Parallel()

				// Assert
				assert.NoError(t, compare.Files(expected, actual))
			})

			t.Run("kubeconform", func(t *testing.T) {
				t.Parallel()

				// Arrange
				v, err := validator.New(nil, validator.Opts{Strict: true}) // keep kube version as master (under the hood)
				require.NoError(t, err)

				file, err := os.Open(actual) // close is handled by ValidateWithContext ...
				require.NoError(t, err)

				// Act
				results := v.ValidateWithContext(ctx, actual, file)

				// Assert
				for _, result := range results {
					assert.NotContains(t, []validator.Status{validator.Invalid, validator.Error}, result.Status, result.Err)
				}
			})
		})
	}
}

// lint runs a lint on given chart directory with input values file as templating.
func lint(kubeVersion *chartutil.KubeVersion, chartdir, valuesFile string) (*action.LintResult, error) {
	// load values
	rawValues, err := os.ReadFile(valuesFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var values map[string]any
	if err := yaml.Unmarshal(rawValues, &values); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	client := action.NewLint()
	client.Strict = true
	client.KubeVersion = kubeVersion

	// run lint
	result := client.Run([]string{chartdir}, values)
	return result, nil
}

// template runs a dry run of an helm install and returns the computed manifest.
func template(ctx context.Context, kubeVersion *chartutil.KubeVersion, chartdir, valuesFile string) (*release.Release, error) {
	// load chart
	chart, err := loader.LoadDir(chartdir)
	if err != nil {
		return nil, fmt.Errorf("load chart dir: %w", err)
	}

	// load values
	rawValues, err := os.ReadFile(valuesFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var values map[string]any
	if err := yaml.Unmarshal(rawValues, &values); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	client := action.NewInstall(&action.Configuration{})
	client.ClientOnly = true
	client.DryRun = true
	client.KubeVersion = kubeVersion
	client.ReleaseName = chart.Name()

	// run install and retrieve resulting manifest
	r, err := client.RunWithContext(ctx, chart, values)
	if err != nil {
		return r, fmt.Errorf("template chart: %w", err)
	}
	return r, nil
}
