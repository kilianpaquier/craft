package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// GitlabURL is the default gitlab API URL.
const GitlabURL = "https://gitlab.com/api/v4"

// DetectLicense handles the detection of license option in craft configuration.
// It also initializes a gitlab client to retrieve the appropriate license in returned slice of GenerateFunc.
func DetectLicense(ctx context.Context, _ string, metadata *Metadata) ([]ExecFunc, error) {
	if metadata.License == nil {
		return []ExecFunc{removeLicense}, nil
	}

	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(GitlabURL),
		gitlab.WithHTTPClient(cleanhttp.DefaultClient()),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
	if err != nil {
		// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
		// and currently WithBaseURL with fixed URL
		// and WithoutRetries won't throw errors
		// but in any case err must be handled in case it evolves or other options are added
		log.Warnf("failed to initialize gitlab client in license detection, skipping license generation: %s", err.Error())
		return nil, nil
	}

	log.Infof("license detected, %s has license key", craft.File)
	return []ExecFunc{downloadLicense(client)}, nil
}

var _ DetectFunc = DetectLicense // ensure interface is implemented

// downloadLicense returns the GenerateFunc to download the appropriate license file from gitlab API.
func downloadLicense(client *gitlab.Client) ExecFunc {
	return func(ctx context.Context, _ cfs.FS, _, destdir string, metadata Metadata, opts ExecOpts) error {
		dest := filepath.Join(destdir, craft.License)

		// don't fetch template is force on file or force all isn't activated
		if !opts.ForceAll && cfs.Exists(dest) && !slices.Contains(opts.Force, craft.License) {
			return nil
		}

		// fetch license template
		options := &gitlab.GetLicenseTemplateOptions{
			Fullname: &metadata.Maintainers[0].Name,
			Project:  &metadata.ProjectName,
		}
		license, _, err := client.LicenseTemplates.GetLicenseTemplate(*metadata.License, options, gitlab.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("license template retrieval: %w", err)
		}

		// remove file before rewritting it (in case rights changed)
		if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("delete file: %w", err)
		}

		// write license template
		if err := os.WriteFile(dest, []byte(license.Content), cfs.RwRR); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		return nil
	}
}

// removeLicense deletes the license file in input destdir.
func removeLicense(_ context.Context, _ cfs.FS, _, destdir string, _ Metadata, _ ExecOpts) error {
	if err := os.Remove(filepath.Join(destdir, craft.License)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
