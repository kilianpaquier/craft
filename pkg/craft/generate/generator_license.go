package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/generator"
)

// GeneratorLicense generates the license file for the project.
func GeneratorLicense(httpClient *http.Client) func(ctx context.Context, destdir string, config craft.Config) error {
	if httpClient == nil {
		httpClient = http.DefaultClient //nolint:revive
	}
	return func(ctx context.Context, destdir string, config craft.Config) error {
		client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
			gitlab.WithBaseURL(generator.GitLabURL),
			gitlab.WithHTTPClient(httpClient),
			gitlab.WithoutRetries(),
			gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
		if err != nil {
			// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
			// and currently WithBaseURL with fixed URL
			// and WithoutRetries won't throw errors
			// but in any case err must be handled in case it evolves or other options are added
			engine.GetLogger().Warnf("failed to initialize GitLab client, skipping license generation: %v", err)
			return nil
		}

		dest := filepath.Join(destdir, generator.FileLicense)
		if config.License == "" {
			engine.GetLogger().Infof("skipping license generation, configuration doesn't have 'license' key")
			if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("remove '%s': %w", generator.FileLicense, err)
			}
			return nil
		}

		if files.Exists(dest) {
			engine.GetLogger().Infof("not generating '%s' since it already exists", generator.FileLicense)
			return nil
		}
		engine.GetLogger().Infof("license detected, configuration has 'license' key")

		opts := generator.LicenseOptions{
			Client:  client,
			License: config.License,
			Maintainer: func() *string {
				var zero string
				if len(config.Maintainers) == 0 {
					return &zero
				}
				if config.Maintainers[0] == nil {
					return &zero
				}
				return &config.Maintainers[0].Name
			}(),
			Project: &config.ProjectName,
		}
		if err := generator.DownloadLicense(dest, opts); err != nil {
			return fmt.Errorf("download license: %w", err)
		}
		return nil
	}
}

var _ engine.Generator[craft.Config] = GeneratorLicense(nil) // ensure interface is implemented
