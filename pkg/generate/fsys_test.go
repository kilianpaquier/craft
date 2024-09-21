package generate_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	testfs "github.com/kilianpaquier/cli-sdk/pkg/cfs/tests"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestIsGenerated(t *testing.T) {
	t.Run("generated_doesnt_exist", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "invalid.txt")

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("not_generated_file", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("not generated"), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.False(t, generated)
	})

	t.Run("generated_folder", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "folder")
		require.NoError(t, os.Mkdir(dest, cfs.RwxRxRxRx))

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_no_lines", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, file.Close()) })

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_first_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("// Code generated by craft; DO NOT EDIT."), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_md_comment", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("<!-- Code generated by craft; DO NOT EDIT. -->"), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_second_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("\n# Code generated by craft; DO NOT EDIT."), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_json", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte(`{
			"//": "Code generated by craft; DO NOT EDIT.",
		}`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		generated := generate.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})
}

func TestExec_Generic(t *testing.T) {
	ctx := context.Background()
	exec := generate.DefaultExec("lang_generic")

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, exec, "generic")

	t.Run("success_releases_github", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Github,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{Auto: true, Backmerge: true},
				},
				NoMakefile: true,
				Platform:   craft.Github,
			},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})

	t.Run("success_options", func(t *testing.T) {
		for _, tc := range []string{craft.Github, craft.Gitlab} {
			t.Run("success_options_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:    tc,
							Options: craft.CIOptions(),
						},
						NoMakefile: true,
						Platform:   tc,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_release_gitlab", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Gitlab,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{Auto: true, Backmerge: true},
				},
				NoMakefile: true,
				Platform:   craft.Gitlab,
			},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})

	t.Run("success_bot_github", func(t *testing.T) {
		for _, tc := range []string{craft.Dependabot, craft.Renovate} {
			t.Run("success_bot_github_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						Bot:        helpers.ToPtr(tc),
						NoMakefile: true,
						Platform:   craft.Github,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_makefile", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{NoMakefile: false},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})
}

func TestExec_Golang(t *testing.T) {
	ctx := context.Background()
	golang := generate.DefaultExec("lang_golang")

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Languages = map[string]any{"golang": generate.Gomod{LangVersion: "1.22"}}
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, golang, "golang")

	t.Run("success_releases_github", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Github,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.PersonalToken)},
					Release: &craft.Release{Backmerge: true},
				},
				NoMakefile: true,
				Platform:   craft.Github,
			},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})

	t.Run("success_options", func(t *testing.T) {
		for _, tc := range []string{craft.Github, craft.Gitlab} {
			t.Run("success_options_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:    tc,
							Options: craft.CIOptions(),
						},
						NoMakefile: true,
						Platform:   tc,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_release_gitlab", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Gitlab,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{Backmerge: true},
				},
				Platform: craft.Gitlab,
			},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_bot_github", func(t *testing.T) {
		for _, tc := range []string{craft.Dependabot, craft.Renovate} {
			t.Run("success_bot_github_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						Bot:        helpers.ToPtr(tc),
						NoMakefile: true,
						Platform:   craft.Github,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_binary_none", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_binary_cli", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Binaries: 1,
			Clis:     map[string]struct{}{"cli-name": {}},
			Configuration: craft.Configuration{
				Docker: &craft.Docker{},
			},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_binary_cron", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Binaries: 1,
			Configuration: craft.Configuration{
				Docker: &craft.Docker{},
			},
			Crons: map[string]struct{}{"cron-name": {}},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_binary_job", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Binaries: 1,
			Configuration: craft.Configuration{
				Docker: &craft.Docker{},
			},
			Jobs: map[string]struct{}{"job-name": {}},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_binary_worker", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Binaries: 1,
			Configuration: craft.Configuration{
				Docker: &craft.Docker{},
			},
			Workers: map[string]struct{}{"worker-name": {}},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_binaries", func(t *testing.T) {
		for _, tc := range []string{craft.Github, craft.Gitlab} {
			t.Run("success_binary_all_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Binaries: 4,
					Clis:     map[string]struct{}{"cli-name": {}},
					Configuration: craft.Configuration{
						CI:           &craft.CI{Name: tc},
						Docker:       &craft.Docker{Port: helpers.ToPtr(uint16(5000)), Registry: helpers.ToPtr("example.com")},
						License:      helpers.ToPtr("mit"),
						NoGoreleaser: true,
						Platform:     tc,
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
	})
}

func TestExec_Hugo(t *testing.T) {
	ctx := context.Background()
	hugo := generate.DefaultExec("lang_hugo")

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Languages = map[string]any{"hugo": generate.Gomod{LangVersion: "1.22"}}
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, hugo, "hugo")

	t.Run("success_releases_github", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Github,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubApp)},
					Release: &craft.Release{},
				},
				NoMakefile: true,
				Platform:   craft.Github,
			},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})

	t.Run("success_options", func(t *testing.T) {
		for _, tc := range []string{craft.Github, craft.Gitlab} {
			t.Run("success_options_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:    tc,
							Options: craft.CIOptions(),
						},
						NoMakefile: true,
						Platform:   tc,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_release_gitlab", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Gitlab,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{},
				},
				Platform: craft.Gitlab,
			},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_bots_github", func(t *testing.T) {
		for _, tc := range []string{craft.Dependabot, craft.Renovate} {
			t.Run("success_bot_github_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						Bot:        helpers.ToPtr(tc),
						NoMakefile: true,
						Platform:   craft.Github,
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_statics", func(t *testing.T) {
		cases := []struct {
			CI     string
			Static string
		}{
			{CI: craft.Github, Static: craft.Netlify},
			{CI: craft.Github, Static: craft.Pages},
			{CI: craft.Gitlab, Static: craft.Netlify},
			{CI: craft.Gitlab, Static: craft.Pages},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("success_static_%s_%s", tc.Static, tc.CI), func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:   tc.CI,
							Static: &craft.Static{Name: tc.Static},
						},
						Platform: tc.CI,
					},
				})
				destdir := t.TempDir()

				// Act & Assert
				verify(t, destdir, metadata)
			})
		}
	})
}

func TestExec_Nodejs(t *testing.T) {
	ctx := context.Background()
	nodejs := generate.DefaultExec("lang_nodejs")

	setup := func(metadata generate.Metadata) generate.Metadata {
		metadata.Maintainers = []*craft.Maintainer{{Name: "maintainer name"}}
		metadata.NoMakefile = true
		metadata.ProjectHost = "example.com"
		metadata.ProjectName = "craft"
		metadata.ProjectPath = "kilianpaquier/craft"
		return metadata
	}

	verify := test(ctx, nodejs, "nodejs")

	t.Run("success_release_github", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Github,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{},
				},
				NoMakefile: true,
				Platform:   craft.Github,
			},
			Languages: map[string]any{
				"nodejs": generate.PackageJSON{PackageManager: "bun@1.0.0"},
			},
		})
		destdir := t.TempDir()

		// Act & Asset
		verify(t, destdir, metadata)
	})

	t.Run("success_options", func(t *testing.T) {
		for _, tc := range []string{craft.Github, craft.Gitlab} {
			t.Run("success_options_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:    tc,
							Options: craft.CIOptions(),
						},
						NoMakefile: true,
						Platform:   tc,
					},
					Languages: map[string]any{
						"nodejs": generate.PackageJSON{PackageManager: "pnpm@9.0.0"},
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_release_gitlab", func(t *testing.T) {
		// Arrange
		metadata := setup(generate.Metadata{
			Configuration: craft.Configuration{
				CI: &craft.CI{
					Name:    craft.Gitlab,
					Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
					Release: &craft.Release{},
				},
				Platform: craft.Gitlab,
			},
			Languages: map[string]any{
				"nodejs": generate.PackageJSON{PackageManager: "pnpm@9.0.0"},
			},
		})
		destdir := t.TempDir()

		// Act & Assert
		verify(t, destdir, metadata)
	})

	t.Run("success_bots_github", func(t *testing.T) {
		for _, tc := range []string{craft.Dependabot, craft.Renovate} {
			t.Run("success_both_github_"+tc, func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						Bot:        helpers.ToPtr(tc),
						NoMakefile: true,
						Platform:   craft.Github,
					},
					Languages: map[string]any{
						"nodejs": generate.PackageJSON{PackageManager: "yarn@2.4.3"},
					},
				})
				destdir := t.TempDir()

				// Act & Asset
				verify(t, destdir, metadata)
			})
		}
	})

	t.Run("success_statics", func(t *testing.T) {
		cases := []struct {
			CI     string
			Static string
		}{
			{CI: craft.Github, Static: craft.Netlify},
			{CI: craft.Github, Static: craft.Pages},
			{CI: craft.Gitlab, Static: craft.Netlify},
			{CI: craft.Gitlab, Static: craft.Pages},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("success_static_%s_%s", tc.Static, tc.CI), func(t *testing.T) {
				// Arrange
				metadata := setup(generate.Metadata{
					Configuration: craft.Configuration{
						CI: &craft.CI{
							Name:   tc.CI,
							Static: &craft.Static{Name: tc.Static},
						},
						Platform: tc.CI,
					},
					Languages: map[string]any{
						"nodejs": generate.PackageJSON{PackageManager: "bun@1.0.0"},
					},
				})
				destdir := t.TempDir()

				// Act & Assert
				verify(t, destdir, metadata)
			})
		}
	})
}

// test returns the verify function for every generation verification to do.
func test(ctx context.Context, exec generate.Exec, name string) func(t *testing.T, destdir string, metadata generate.Metadata) {
	srcdir := "templates"
	assertdir := filepath.Join("..", "..", "testdata", name)

	return func(t *testing.T, destdir string, metadata generate.Metadata) {
		t.Helper()

		// Arrange
		assertdir := filepath.Join(assertdir, path.Base(t.Name()))

		opts := generate.ExecOpts{
			FileHandlers: func() []generate.FileHandler {
				metas := generate.MetaHandlers()
				result := make([]generate.FileHandler, 0, len(metas))
				for _, handler := range metas {
					result = append(result, handler(metadata))
				}
				return result
			}(),
			EndDelim:   ">>",
			StartDelim: "<<",
			ForceAll:   true,
		}

		// Act
		err := exec(ctx, clog.Noop(), cfs.OS(), srcdir, destdir, metadata, opts)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
	}
}
