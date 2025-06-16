package generate_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

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
	"github.com/kilianpaquier/craft/testutils"
)

func TestGenerate_NoLang(t *testing.T) {
	ctx := t.Context()

	t.Run("success_chart", func(t *testing.T) {
		// Arrange
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:      &craft.CI{Name: ci},
					Exclude: []string{craft.Makefile},
					VCS:     parser.VCS{Platform: ci},
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
					Bot:     helpers.ToPtr(craft.Renovate),
					CI:      &craft.CI{Name: ci},
					Exclude: []string{craft.Chart, craft.Makefile},
					VCS:     parser.VCS{Platform: ci},
				}

				// Act & Assert
				test(ctx, t, config)
			})
		}
	})

	t.Run("success_precommit", func(t *testing.T) {
		for _, precommit := range []bool{true, false} {
			t.Run(strconv.FormatBool(precommit), func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:      &craft.CI{Name: parser.GitHub},
					Exclude: []string{craft.Chart, craft.Makefile},
				}
				if !precommit {
					config.Exclude = append(config.Exclude, craft.PreCommit)
				} else {
					config.Include = append(config.Include, craft.PreCommit+":auto-commit")
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
					CI:      &ci,
					Exclude: []string{craft.Chart, craft.Makefile},
					VCS:     parser.VCS{Platform: ci.Name},
				}

				// Act & Assert
				test(ctx, t, config)
			})
		}
	})
}

func TestGenerate_Shell(t *testing.T) {
	ctx := t.Context()

	shell := func(_ context.Context, _ string, config *craft.Config) error {
		config.SetLanguage("shell", nil)
		return nil
	}

	t.Run("success_ci", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:      &craft.CI{Name: ci},
					Exclude: []string{craft.Chart, craft.Makefile},
				}

				// Act & Assert
				test(ctx, t, config, shell)
			})
		}
	})

	t.Run("success_precommit", func(t *testing.T) {
		for _, precommit := range []bool{true, false} {
			t.Run(strconv.FormatBool(precommit), func(t *testing.T) {
				// Arrange
				config := craft.Config{Exclude: []string{craft.Chart, craft.Makefile}}
				if !precommit {
					config.Exclude = append(config.Exclude, craft.PreCommit)
				}

				// Act & Assert
				test(ctx, t, config, shell)
			})
		}
	})
}

func TestGenerate_Golang(t *testing.T) {
	ctx := t.Context()

	t.Run("success_cli", func(t *testing.T) {
		for _, ci := range []string{parser.GitLab, parser.GitHub} {
			t.Run(ci, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					Bot:     helpers.ToPtr(craft.Dependabot),
					CI:      &craft.CI{Name: ci, Release: &craft.Release{}},
					Exclude: []string{craft.Chart},
					VCS:     parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetCLI("name")

					gomod := parser.Gomod{
						Go:     "1.23",
						Module: ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("go", gomod)
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
					Exclude: []string{craft.Chart, craft.Makefile},
					VCS:     parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					gomod := parser.Gomod{
						Go:     "1.23",
						Module: ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("go", gomod)
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
					Exclude:     []string{craft.Shell},
					VCS:         parser.VCS{Platform: ci},
				}
				golang := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetJob("job-name")
					config.SetCron("cron-name")
					config.SetWorker("worker-name")

					gomod := parser.Gomod{
						Go:     "1.23",
						Module: ci + ".com/kilianpaquier/craft",
					}
					config.VCS = gomod.AsVCS()
					config.SetLanguage("go", gomod)
					return nil
				}

				// Act & Assert
				test(ctx, t, config, golang)
			})
		}
	})
}

func TestGenerate_Hugo(t *testing.T) {
	ctx := t.Context()

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
				CI:      &ci,
				Exclude: []string{craft.Chart},
				VCS:     parser.VCS{Platform: ci.Name},
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

func TestGenerate_Node(t *testing.T) {
	ctx := t.Context()

	t.Run("success_package_managers", func(t *testing.T) {
		for _, tc := range []string{"bun@1.1.6", "npm@7.0.0", "pnpm@9.0.0", "yarn@1.22.10"} {
			t.Run(tc, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					CI:      &craft.CI{Name: parser.GitHub},
					Exclude: []string{craft.Chart},
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
		type testcase struct {
			Bot            string
			CI             string
			PackageManager string
		}
		cases := []testcase{
			{Bot: craft.Renovate, CI: parser.GitLab, PackageManager: "bun@1.1.6"},
			{Bot: craft.Dependabot, CI: parser.GitLab, PackageManager: "bun@1.1.6"},
			{Bot: craft.Renovate, CI: parser.GitLab, PackageManager: "npm@7.0.0"},
			{Bot: craft.Dependabot, CI: parser.GitLab, PackageManager: "npm@7.0.0"},

			{Bot: craft.Renovate, CI: parser.GitHub, PackageManager: "bun@1.1.6"},
			{Bot: craft.Dependabot, CI: parser.GitHub, PackageManager: "bun@1.1.6"},
			{Bot: craft.Renovate, CI: parser.GitHub, PackageManager: "npm@7.0.0"},
			{Bot: craft.Dependabot, CI: parser.GitHub, PackageManager: "npm@7.0.0"},
		}

		for _, tc := range cases {
			name := fmt.Sprint(tc.CI, "_", tc.Bot, "_", tc.PackageManager)
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Config{
					Bot: helpers.ToPtr(tc.Bot),
					CI: &craft.CI{
						Name:    tc.CI,
						Auth:    craft.Auth{Maintenance: helpers.ToPtr(craft.PersonalToken)},
						Release: &craft.Release{Backmerge: true},
					},
					Exclude: []string{craft.Chart},
					VCS:     parser.VCS{Platform: tc.CI},
				}
				node := func(_ context.Context, _ string, config *craft.Config) error {
					config.SetLanguage("node", parser.PackageJSON{Name: "craft", PackageManager: tc.PackageManager})
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
					CI:      &ci,
					Exclude: []string{craft.Chart, craft.Makefile},
					VCS:     parser.VCS{Platform: ci.Name},
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
	assertdir := filepath.Join(testutils.Testdata(t), t.Name())
	require.NoError(t, os.MkdirAll(assertdir, files.RwxRxRxRx))

	// Act
	_, err := engine.Generate(ctx, destdir, config,
		slices.Concat(parsers, []engine.Parser[craft.Config]{ParserInfo, generate.ParserGolang, generate.ParserNode, generate.ParserShell, generate.ParserChart}),
		[]engine.Generator[craft.Config]{
			engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.Dependabot(), templates.Renovate())),
			engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.CodeCov(), templates.Sonar())),
			engine.GeneratorTemplates(templates.FS(), templates.Docker()),
			engine.GeneratorTemplates(templates.FS(), templates.Golang()),
			engine.GeneratorTemplates(templates.FS(), templates.Misc()),
			engine.GeneratorTemplates(templates.FS(), templates.Makefile()),
			engine.GeneratorTemplates(templates.FS(), templates.Chart()),
			engine.GeneratorTemplates(templates.FS(), slices.Concat(templates.GitHub(), templates.GitLab(), templates.SemanticRelease())),
		})

	// Assert
	require.NoError(t, err)
	assert.NoError(t, compare.Dirs(assertdir, destdir))
}
