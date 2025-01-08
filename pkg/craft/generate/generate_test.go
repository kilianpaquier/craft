package generate_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jarcoal/httpmock"
	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestRun_NoLang(t *testing.T) {
	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	ctx := context.Background()

	t.Run("success_chart", func(t *testing.T) {
		// Arrange
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:         &craft.CI{Name: ci},
					NoMakefile: true,
					VCS:        parser.VCS{Platform: ci},
				}

				// Act & Assert
				test(ctx, t, config)
			})
		}
	})

	t.Run("success_renovate", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					Bot:        helpers.ToPtr(craft.Renovate),
					CI:         &craft.CI{Name: ci},
					NoChart:    true,
					NoMakefile: true,
					VCS:        parser.VCS{Platform: ci},
				}

				// Act & Assert
				test(ctx, t, config)
			})
		}
	})

	t.Run("success_release", func(t *testing.T) {
		cases := []craft.CI{
			{Name: parser.GitHub, Release: &craft.Release{}},
			{Name: parser.GitHub, Release: &craft.Release{Auto: true}},
			{Name: parser.GitLab, Release: &craft.Release{}},
			{Name: parser.GitLab, Release: &craft.Release{Auto: true}},
		}
		for _, ci := range cases {
			name := fmt.Sprint(ci.Name, "_auto_", ci.Release.Auto)
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:         &ci,
					NoChart:    true,
					NoMakefile: true,
					VCS:        parser.VCS{Platform: ci.Name},
				}

				// Act & Assert
				test(ctx, t, config)
			})
		}
	})
}

func TestRun_Golang(t *testing.T) {
	ctx := context.Background()

	t.Run("success_cli", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					Bot:     helpers.ToPtr(craft.Dependabot),
					CI:      &craft.CI{Name: ci, Release: &craft.Release{}},
					NoChart: true,
					VCS:     parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetCLI("name")

					gomod := parser.Gomod{
						LangVersion: "1.23",
						ModulePath:  ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("golang", gomod)
					return nil
				}

				// Act & Assert
				test(ctx, t, config, golang)
			})
		}
	})

	t.Run("success_library", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI: &craft.CI{
						Name:    ci,
						Options: []string{craft.Sonar, craft.CodeQL, craft.Labeler},
						Release: &craft.Release{},
					},
					NoChart:    true,
					NoMakefile: true,
					VCS:        parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					gomod := parser.Gomod{
						LangVersion: "1.23",
						ModulePath:  ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("golang", gomod)
					return nil
				}

				// Act & Assert
				test(ctx, t, config, golang)
			})
		}
	})

	t.Run("success_multiple_bin_helm", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI: &craft.CI{
						Name:    ci,
						Options: []string{craft.CodeCov, craft.CodeQL, craft.Labeler},
					},
					Description: helpers.ToPtr("A useful project description"),
					Docker:      &craft.Docker{},
					VCS:         parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetJob("job-name")
					config.SetCron("cron-name")
					config.SetWorker("worker-name")

					gomod := parser.Gomod{
						LangVersion: "1.23",
						ModulePath:  ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("golang", gomod)
					return nil
				}

				// Act & Assert
				test(ctx, t, config, golang)
			})
		}
	})
}

func TestRun_Hugo(t *testing.T) {
	ctx := context.Background()

	cases := []craft.CI{
		{Name: parser.GitHub, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
		{Name: parser.GitHub, Static: &craft.Static{Name: craft.Netlify}},
		{Name: parser.GitHub, Static: &craft.Static{Name: craft.Pages, Auto: true}},
		{Name: parser.GitHub, Static: &craft.Static{Name: craft.Pages}},
		{Name: parser.GitLab, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
		{Name: parser.GitLab, Static: &craft.Static{Name: craft.Netlify}},
		{Name: parser.GitLab, Static: &craft.Static{Name: craft.Pages, Auto: true}},
		{Name: parser.GitLab, Static: &craft.Static{Name: craft.Pages}},
	}
	for _, ci := range cases {
		name := fmt.Sprint(ci.Name, "_", ci.Static.Name, "_auto_", ci.Static.Auto)
		t.Run(name, func(t *testing.T) {
			// Arrange
			config := craft.Config{
				CI:         &ci,
				NoChart:    true,
				NoMakefile: true,
				VCS:        parser.VCS{Platform: ci.Name},
			}
			hugo := func(_ context.Context, _ string, config *craft.Config) error {
				config.SetLanguage("hugo", nil)
				return nil
			}

			// Act & Assert
			test(ctx, t, config, hugo)
		})
	}
}

func TestRun_Node(t *testing.T) {
	ctx := context.Background()

	t.Run("success_package_managers", func(t *testing.T) {
		for _, tc := range []string{"bun@1.1.6", "npm@7.0.0", "pnpm@9.0.0", "yarn@1.22.10"} {
			t.Run(tc, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:      &craft.CI{Name: parser.GitHub},
					NoChart: true,
					VCS:     parser.VCS{Platform: parser.GitHub},
				}
				node := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetWorker("index.js")
					config.SetLanguage("node", parser.PackageJSON{Name: "craft", PackageManager: tc})
					return nil
				}

				// Act & Assert
				test(ctx, t, config, node)
			})
		}
	})

	t.Run("success_library", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					Bot: helpers.ToPtr(craft.Renovate),
					CI: &craft.CI{
						Name:    ci,
						Auth:    craft.Auth{Maintenance: helpers.ToPtr(craft.PersonalToken)},
						Release: &craft.Release{Backmerge: true},
					},
					NoChart: true,
					VCS:     parser.VCS{Platform: ci},
				}
				node := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetLanguage("node", parser.PackageJSON{Name: "craft", PackageManager: "bun@1.1.6"})
					return nil
				}

				// Act & Assert
				test(ctx, t, config, node)
			})
		}
	})

	t.Run("success_static", func(t *testing.T) {
		cases := []craft.CI{
			{Name: parser.GitHub, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
			{Name: parser.GitHub, Static: &craft.Static{Name: craft.Netlify}},
			{Name: parser.GitHub, Static: &craft.Static{Name: craft.Pages, Auto: true}},
			{Name: parser.GitHub, Static: &craft.Static{Name: craft.Pages}},
			{Name: parser.GitLab, Static: &craft.Static{Name: craft.Netlify, Auto: true}},
			{Name: parser.GitLab, Static: &craft.Static{Name: craft.Netlify}},
			{Name: parser.GitLab, Static: &craft.Static{Name: craft.Pages, Auto: true}},
			{Name: parser.GitLab, Static: &craft.Static{Name: craft.Pages}},
		}
		for _, ci := range cases {
			name := fmt.Sprint(ci.Name, "_", ci.Static.Name, "_auto_", ci.Static.Auto)
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:         &ci,
					NoChart:    true,
					NoMakefile: true,
					VCS:        parser.VCS{Platform: ci.Name},
				}
				node := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetWorker("index.js")
					config.SetLanguage("node", parser.PackageJSON{Name: "craft", PackageManager: "bun@1.1.6"})
					return nil
				}

				// Act & Assert
				test(ctx, t, config, node)
			})
		}
	})
}

func ParserInfo(_ context.Context, _ string, config *craft.Config) error {
	config.VCS = parser.VCS{
		Platform:    config.Platform,
		ProjectHost: config.Platform + ".com",
		ProjectName: "craft",
		ProjectPath: "kilianpaquier/craft",
	}
	return nil
}

// test verifies every generation with provided config, parser and t.Name folder expected results.
func test(ctx context.Context, t *testing.T, config craft.Config, parsers ...engine.Parser[craft.Config]) {
	t.Helper()

	// Arrange
	config.Maintainers = append(config.Maintainers, &craft.Maintainer{Name: "kilianpaquier"})
	destdir := t.TempDir()
	assertdir := filepath.Join("..", "..", "..", "testdata", t.Name())
	require.NoError(t, os.MkdirAll(assertdir, files.RwxRxRxRx))

	// Act
	_, err := engine.Generate(ctx, config,
		engine.WithDestination[craft.Config](destdir),
		engine.WithTemplates(templates.FS(), templates.All()),
		engine.WithParsers(slices.Concat(
			parsers,
			[]engine.Parser[craft.Config]{ParserInfo, generate.ParserLicense, generate.ParserGolang, generate.ParserNode, generate.ParserHelm},
		)...))

	// Assert
	require.NoError(t, err)
	assert.NoError(t, compare.Dirs(assertdir, destdir))
}
