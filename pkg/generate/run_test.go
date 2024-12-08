package generate_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestRun_Error(t *testing.T) {
	ctx := context.Background()

	t.Run("error_missing_handlers_parsers", func(t *testing.T) {
		// Act
		_, err := generate.Run(ctx, craft.Configuration{})

		// Assert
		assert.ErrorIs(t, err, generate.ErrMissingHandlers)
		assert.ErrorIs(t, err, generate.ErrMissingParsers)
	})

	t.Run("error_parsing", func(t *testing.T) {
		// Arrange
		parser := func(context.Context, string, *generate.Metadata) error { return errors.New("some error") }

		// Act
		_, err := generate.Run(ctx, craft.Configuration{}, generate.WithHandlers(generate.HandlerNoop), generate.WithParsers(parser))

		// Assert
		assert.ErrorContains(t, err, "some error")
	})
}

func TestRun_NoLang(t *testing.T) {
	t.Run("success_chart", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:         &craft.CI{Name: ci},
					License:    helpers.ToPtr("mit"),
					NoMakefile: true,
					Platform:   ci,
				}

				// Act & Assert
				test(t, config)
			})
		}
	})

	t.Run("success_renovate", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					Bot:        helpers.ToPtr(craft.Renovate),
					CI:         &craft.CI{Name: ci},
					NoChart:    true,
					NoMakefile: true,
					Platform:   ci,
				}

				// Act & Assert
				test(t, config)
			})
		}
	})

	t.Run("success_release", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:         &craft.CI{Name: ci, Release: &craft.Release{}},
					NoChart:    true,
					NoMakefile: true,
					Platform:   ci,
				}

				// Act & Assert
				test(t, config)
			})
		}
	})
}

func TestRun_Golang(t *testing.T) {
	t.Run("success_cli", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:       &craft.CI{Name: ci, Release: &craft.Release{}},
					NoChart:  true,
					Platform: ci,
				}
				golang := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Binaries++
					metadata.Clis["name"] = struct{}{}
					metadata.Languages["golang"] = nil
					return nil
				}

				// Act & Assert
				test(t, config, golang)
			})
		}
	})

	t.Run("success_library", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:       &craft.CI{Name: ci, Release: &craft.Release{}},
					NoChart:  true,
					Platform: ci,
				}
				golang := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Languages["golang"] = nil
					return nil
				}

				// Act & Assert
				test(t, config, golang)
			})
		}
	})

	t.Run("success_multiple_bin_helm", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:          &craft.CI{Name: ci},
					Description: helpers.ToPtr("A useful project description"),
					Docker:      &craft.Docker{},
					NoMakefile:  true,
					Platform:    ci,
				}
				golang := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Binaries += 3
					metadata.Jobs["job-name"] = struct{}{}
					metadata.Crons["cron-name"] = struct{}{}
					metadata.Workers["worker-name"] = struct{}{}
					metadata.Languages["golang"] = parser.Gomod{LangVersion: "1.23"}
					return nil
				}

				// Act & Assert
				test(t, config, golang)
			})
		}
	})
}

func TestRun_Hugo(t *testing.T) {
	cases := []craft.CI{
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify}},
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages}},
	}
	for _, ci := range cases {
		name := fmt.Sprint("success_", ci.Name, "_", ci.Static.Name)
		t.Run(name, func(t *testing.T) {
			// Arrange
			config := craft.Configuration{
				CI:         &ci,
				NoChart:    true,
				NoMakefile: true,
				Platform:   ci.Name,
			}
			hugo := func(_ context.Context, _ string, metadata *generate.Metadata) error {
				metadata.Languages["hugo"] = nil
				return nil
			}

			// Act & Assert
			test(t, config, hugo)
		})
	}
}

func TestRun_Node(t *testing.T) {
	t.Run("success_library", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:       &craft.CI{Name: ci, Release: &craft.Release{Backmerge: true}},
					NoChart:  true,
					Platform: ci,
				}
				node := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Languages["node"] = parser.PackageJSON{Name: "craft", PackageManager: "bun@1.1.6"}
					return nil
				}

				// Act & Assert
				test(t, config, node)
			})
		}
	})

	t.Run("success_static", func(t *testing.T) {
		cases := []craft.CI{
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify}},
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages}},
		}
		for _, ci := range cases {
			name := fmt.Sprint(ci.Name, "_", ci.Static.Name)
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:         &ci,
					NoChart:    true,
					NoMakefile: true,
					Platform:   ci.Name,
				}
				node := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Binaries++
					metadata.Languages["node"] = parser.PackageJSON{Name: "craft", PackageManager: "bun@1.1.6"}
					return nil
				}

				// Act & Assert
				test(t, config, node)
			})
		}
	})
}

// test returns the verify function for every generation verification to do.
func test(t *testing.T, config craft.Configuration, parsers ...generate.Parser) {
	t.Helper()

	// add a parser to bypass git parsing in t.TempDir() not setting current craft properties
	info := func(_ context.Context, _ string, metadata *generate.Metadata) error {
		metadata.Maintainers = append(metadata.Maintainers, &craft.Maintainer{Name: "kilianpaquier"})
		metadata.ProjectHost = "github.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return nil
	}
	parsers = append(parsers, info) //nolint:revive

	// Arrange
	destdir := t.TempDir()
	assertdir := filepath.Join("..", "..", "testdata", t.Name())
	require.NoError(t, os.MkdirAll(assertdir, cfs.RwxRxRxRx))

	// Act
	_, err := generate.Run(context.Background(), config,
		generate.WithDestination(assertdir),
		generate.WithHandlers(handler.Defaults()...),
		generate.WithParsers(parser.Defaults(parsers...)...))

	// Assert
	require.NoError(t, err)
	assert.NoError(t, tests.EqualDirs(assertdir, destdir))
}
