package generate_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jarcoal/httpmock"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

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
	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	ctx := context.WithValue(context.Background(), parser.HTTPClientKey, httpClient)
	url := parser.GitLabURL + "/templates/licenses/mit"

	info := func(_ context.Context, _ string, metadata *generate.Metadata) error {
		metadata.ProjectHost = "github.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return nil
	}

	t.Run("success_chart", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

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
				test(ctx, t, config, parser.Defaults(info)...)
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
				test(ctx, t, config, parser.Defaults(info)...)
			})
		}
	})

	t.Run("success_release", func(t *testing.T) {
		cases := []craft.CI{
			{Name: craft.GitHub, Release: &craft.Release{}},
			{Name: craft.GitHub, Release: &craft.Release{Auto: true}},
			{Name: craft.GitLab, Release: &craft.Release{}},
			{Name: craft.GitLab, Release: &craft.Release{Auto: true}},
		}
		for _, ci := range cases {
			name := fmt.Sprint(ci.Name, "_auto_", ci.Release.Auto)
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:         &ci,
					NoChart:    true,
					NoMakefile: true,
					Platform:   ci.Name,
				}

				// Act & Assert
				test(ctx, t, config, parser.Defaults(info)...)
			})
		}
	})
}

func TestRun_Golang(t *testing.T) {
	ctx := context.Background()

	info := func(_ context.Context, _ string, metadata *generate.Metadata) error {
		metadata.ProjectHost = "github.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return nil
	}

	t.Run("success_cli", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					Bot:      helpers.ToPtr(craft.Dependabot),
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
				test(ctx, t, config, parser.Defaults(info, golang)...)
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
				test(ctx, t, config, parser.Defaults(info, golang)...)
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
				test(ctx, t, config, parser.Defaults(info, golang)...)
			})
		}
	})
}

func TestRun_Hugo(t *testing.T) {
	ctx := context.Background()

	info := func(_ context.Context, _ string, metadata *generate.Metadata) error {
		metadata.ProjectHost = "github.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return nil
	}

	cases := []craft.CI{
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify}},
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages, Auto: true}},
		{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages, Auto: true}},
		{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages}},
	}
	for _, ci := range cases {
		name := fmt.Sprint(ci.Name, "_", ci.Static.Name, "_auto_", ci.Static.Auto)
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
			test(ctx, t, config, parser.Defaults(info, hugo)...)
		})
	}
}

func TestRun_Node(t *testing.T) {
	ctx := context.Background()

	info := func(_ context.Context, _ string, metadata *generate.Metadata) error {
		metadata.ProjectHost = "github.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return nil
	}

	t.Run("success_package_managers", func(t *testing.T) {
		for _, tc := range []string{"bun@1.1.6", "npm@7.0.0", "pnpm@9.0.0", "yarn@1.22.10"} {
			t.Run(tc, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					CI:       &craft.CI{Name: craft.GitHub},
					NoChart:  true,
					Platform: craft.GitHub,
				}
				node := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Binaries++
					metadata.Languages["node"] = parser.PackageJSON{Name: "craft", PackageManager: tc}
					return nil
				}

				// Act & Assert
				test(ctx, t, config, parser.Defaults(info, node)...)
			})
		}
	})

	t.Run("success_library", func(t *testing.T) {
		for _, ci := range []string{craft.GitLab, craft.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Configuration{
					Bot:      helpers.ToPtr(craft.Dependabot),
					CI:       &craft.CI{Name: ci, Release: &craft.Release{Backmerge: true}},
					NoChart:  true,
					Platform: ci,
				}
				node := func(_ context.Context, _ string, metadata *generate.Metadata) error {
					metadata.Languages["node"] = parser.PackageJSON{Name: "craft", PackageManager: "bun@1.1.6"}
					return nil
				}

				// Act & Assert
				test(ctx, t, config, parser.Defaults(info, node)...)
			})
		}
	})

	t.Run("success_static", func(t *testing.T) {
		cases := []craft.CI{
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Netlify}},
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages, Auto: true}},
			{Name: craft.GitHub, Static: &craft.Static{Name: craft.Pages}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Netlify}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages, Auto: true}},
			{Name: craft.GitLab, Static: &craft.Static{Name: craft.Pages}},
		}
		for _, ci := range cases {
			name := fmt.Sprint(ci.Name, "_", ci.Static.Name, "_auto_", ci.Static.Auto)
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
				test(ctx, t, config, parser.Defaults(info, node)...)
			})
		}
	})
}

// test returns the verify function for every generation verification to do.
func test(ctx context.Context, t *testing.T, config craft.Configuration, parsers ...generate.Parser) {
	t.Helper()

	// Arrange
	config.Maintainers = append(config.Maintainers, &craft.Maintainer{Name: "kilianpaquier"})
	destdir := t.TempDir()
	assertdir := filepath.Join("..", "..", "testdata", t.Name())
	require.NoError(t, os.MkdirAll(assertdir, cfs.RwxRxRxRx))

	// Act
	_, err := generate.Run(ctx, config,
		generate.WithDestination(destdir),
		generate.WithHandlers(handler.Defaults()...),
		generate.WithParsers(parsers...))

	// Assert
	require.NoError(t, err)
	assert.NoError(t, tests.EqualDirs(assertdir, destdir))
}
