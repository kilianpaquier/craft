package detectgen

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/internal/models"
)

// GitlabURL is the default gitlab API URL.
const GitlabURL = "https://gitlab.com/api/v4"

// detectLicense handles the detection of license option in craft configuration.
// It also initializes a gitlab client to retrieve the appropriate license in returned slice of GenerateFunc.
func detectLicense(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	if config.License == nil {
		return []GenerateFunc{removeLicense}
	}
	log := logrus.WithContext(ctx)

	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(GitlabURL), gitlab.WithoutRetries())
	if err != nil {
		// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
		// and currently WithBaseURL with fixed URL
		// and WithoutRetries won't throw errors
		// but in any case err must be handled in case it evolves or other options are added
		log.WithError(err).
			Warn("failed to initialize gitlab client in license detection, skipping license generation")
		return nil
	}

	log.Infof("license detected, %s has license key", models.CraftFile)
	return []GenerateFunc{downloadLicense(client)}
}

// downloadLicense returns the GenerateFunc to download the appropriate license file from gitlab API.
func downloadLicense(client *gitlab.Client) GenerateFunc {
	return func(ctx context.Context, config models.GenerateConfig, _ filesystem.FS) error {
		dest := filepath.Join(config.Options.DestinationDir, models.License)

		// don't fetch template is force on file or force all isn't activated
		if !config.Options.ForceAll && filesystem.Exists(dest) && !slices.Contains(config.Options.Force, models.License) {
			return nil
		}

		// fetch license template
		options := &gitlab.GetLicenseTemplateOptions{
			Fullname: &config.Maintainers[0].Name,
			Project:  &config.ProjectName,
		}
		license, _, err := client.LicenseTemplates.GetLicenseTemplate(*config.License, options, gitlab.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("license template retrieval: %w", err)
		}

		// remove file before rewritting it (in case rights changed)
		if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("delete file: %w", err)
		}

		// write license template
		if err := os.WriteFile(dest, []byte(license.Content), filesystem.RwRR); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		return nil
	}
}

// removeLicense deletes from config provided destination directory the license file.
func removeLicense(_ context.Context, config models.GenerateConfig, _ filesystem.FS) error {
	dest := filepath.Join(config.Options.DestinationDir, models.License)
	if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
