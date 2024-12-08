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
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

const (
	// GitLabURL is the default GitLab API URL.
	GitLabURL = "https://gitlab.com/api/v4"

	// GitHubURL is the default GitHub API URL.
	GitHubURL = "https://api.github.com"
)

// httpClient is the default http.Client to make requests overs HTTP(s) in craft.
var httpClient = cleanhttp.DefaultClient()

// SetHTTPClient sets the global craft http.Client only if the input one is not nil.
func SetHTTPClient(client *http.Client) {
	if client != nil {
		httpClient = client
	}
}

// License generates the LICENSE file in case input configuration asks for a LICENSE file.
func License(ctx context.Context, destdir string, metadata *generate.Metadata) error {
	dest := filepath.Join(destdir, craft.License)
	if metadata.License == nil {
		if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove '%s': %w", craft.License, err)
		}
		return nil
	}
	// log.Infof("license detected, %s has license key", craft.File)

	// don't do anything if the LICENSE file already exists
	if cfs.Exists(dest) {
		return nil
	}

	// initialize gitlab client
	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(GitLabURL),
		gitlab.WithHTTPClient(httpClient),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
	if err != nil {
		// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
		// and currently WithBaseURL with fixed URL
		// and WithoutRetries won't throw errors
		// but in any case err must be handled in case it evolves or other options are added
		// log.Warnf("failed to initialize gitlab client in license detection, skipping license generation: %s", err.Error())
		return nil
	}

	// fetch license template
	options := &gitlab.GetLicenseTemplateOptions{
		Fullname: &metadata.Maintainers[0].Name,
		Project:  &metadata.ProjectName,
	}
	license, _, err := client.LicenseTemplates.GetLicenseTemplate(*metadata.License, options, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("get license template '%s': %w", *metadata.License, err)
	}

	// write license template
	if err := os.WriteFile(dest, []byte(license.Content), cfs.RwRR); err != nil {
		return fmt.Errorf("write license file: %w", err)
	}
	return nil
}

var _ generate.Parser = License // ensure interface is implemented
