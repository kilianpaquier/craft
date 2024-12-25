package parser

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

const (
	// GitLabURL is the default GitLab API URL.
	GitLabURL = "https://gitlab.com/api/v4"

	// GitHubURL is the default GitHub API URL.
	GitHubURL = "https://api.github.com"
)

// FileLicense represents the target filename for the generated project LICENSE.
const FileLicense = "LICENSE"

// License generates the LICENSE file in case input configuration asks for a LICENSE file.
func License(ctx context.Context, destdir string, config *craft.Config) error {
	dest := filepath.Join(destdir, FileLicense)
	if config.License == nil {
		if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove '%s': %w", FileLicense, err)
		}
		return nil
	}
	generate.GetLogger(ctx).Infof("license detected, %s has license key", craft.File)

	// don't do anything if the LICENSE file already exists
	if cfs.Exists(dest) {
		return nil
	}

	// initialize gitlab client
	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(GitLabURL),
		gitlab.WithHTTPClient(GetHTTPClient(ctx)),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
	if err != nil {
		// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
		// and currently WithBaseURL with fixed URL
		// and WithoutRetries won't throw errors
		// but in any case err must be handled in case it evolves or other options are added
		generate.GetLogger(ctx).Warnf("failed to initialize gitlab client in license detection, skipping license generation: %s", err.Error())
		return nil
	}

	// fetch license template
	options := &gitlab.GetLicenseTemplateOptions{
		Fullname: &config.Maintainers[0].Name,
		Project:  &config.ProjectName,
	}
	license, _, err := client.LicenseTemplates.GetLicenseTemplate(*config.License, options, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("get license template '%s': %w", *config.License, err)
	}

	// write license template
	if err := os.WriteFile(dest, []byte(license.Content), cfs.RwRR); err != nil {
		return fmt.Errorf("write license file: %w", err)
	}
	return nil
}

var _ generate.Parser[craft.Config] = License // ensure interface is implemented

type httpClientKeyType string

// HTTPClientKey is the context key for the http client.
//
// It is used during License parsing in case a custom http client needs to be provided.
// Example:
//
//	ctx := context.WithValue(context.Background(), parser.HTTPClientKey, &http.Client{})
const HTTPClientKey httpClientKeyType = "http_client"

// GetHTTPClient returns the context http client.
//
// By default it will be cleanhttp.DefaultClient, but it can be set following this example:
//
//	ctx := context.WithValue(context.Background(), parser.HTTPClientKey, &http.Client{})
func GetHTTPClient(ctx context.Context) *http.Client {
	client, ok := ctx.Value(HTTPClientKey).(*http.Client)
	if !ok {
		return cleanhttp.DefaultClient()
	}
	return client
}
