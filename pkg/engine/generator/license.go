package generator

import (
	"errors"
	"fmt"
	"os"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

const (
	// GitLabURL is the default GitLab API URL.
	GitLabURL = "https://gitlab.com/api/v4"

	// GitHubURL is the default GitHub API URL.
	GitHubURL = "https://api.github.com"
)

// FileLicense represents the target filename for the generated project LICENSE.
const FileLicense = "LICENSE"

// LicenseOptions represents the arguments required to generate a LICENSE file.
type LicenseOptions struct {
	// Client is the GitLab client used to fetch the license template.
	Client *gitlab.Client

	// License is the license key to fetch the license template.
	License string

	// Maintainer is the maintainer's full name to add in LICENSE file.
	//
	// It's optional, not provided a placeholder will be used.
	Maintainer *string

	// Project is the project name to fetch the license template.
	//
	// It's optional, not provided a placeholder will be used.
	Project *string
}

// ErrNoClient is returned when input HTTP client is nil.
//
// This error can be returned from DownloadLicense and DownloadGitignore.
var ErrNoClient = errors.New("no client provided")

// DownloadLicense generates the LICENSE file.
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func GeneratorLicense(ctx context.Context, destdir string, c config) error {
//		client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
//			gitlab.WithBaseURL(generator.GitLabURL),
//			gitlab.WithHTTPClient(cleanhttp.DefaultClient()),
//			gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
//		if err != nil {
//			// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
//			// and currently WithBaseURL with fixed URL
//			// and WithoutRetries won't throw errors
//			// but in any case err must be handled in case it evolves or other options are added
//			engine.GetLogger().Warnf("failed to initialize GitLab client, skipping license generation: %v", err)
//			return nil
//		}
//
//		dest := filepath.Join(destdir, generator.FileLicense)
//		if config.License == nil {
//			engine.GetLogger().Infof("skipping license generation, configuration doesn't have 'license' key")
//			if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
//				return fmt.Errorf("remove '%s': %w", generator.FileLicense, err)
//			}
//			return nil
//		}
//
//		if files.Exists(dest) {
//			engine.GetLogger().Infof("not generating '%s' since it already exists", generator.FileLicense)
//			return nil
//		}
//		engine.GetLogger().Infof("license detected, configuration has 'license' key")
//
//		opts := generator.LicenseOptions{
//			Client:     client,
//			License:    config.License,
//			Maintainer: &config.Maintainer,
//			Project:    &config.Project,
//		}
//		return generator.DownloadLicense(destdir, opts)
//	}
func DownloadLicense(out string, opts LicenseOptions) error {
	// validate args
	if opts.Client == nil {
		return ErrNoClient
	}

	// fetch license template
	template, _, err := opts.Client.LicenseTemplates.GetLicenseTemplate(opts.License,
		&gitlab.GetLicenseTemplateOptions{
			Project:  opts.Project,
			Fullname: opts.Maintainer,
		})
	if err != nil {
		return fmt.Errorf("get license template '%s': %w", opts.License, err)
	}

	// write license template
	if err := os.WriteFile(out, []byte(template.Content), files.RwRR); err != nil {
		return fmt.Errorf("write license file: %w", err)
	}
	return nil
}
