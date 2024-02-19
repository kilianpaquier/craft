package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/internal/models"
)

const gitlabURL = "https://gitlab.com/api/v4"

type license struct {
	GitlabClient *gitlab.Client
}

var _ plugin = &license{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (plugin *license) Detect(ctx context.Context, config *models.GenerateConfig) bool {
	if config.License != nil {
		client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
			gitlab.WithBaseURL(gitlabURL),
			gitlab.WithoutRetries(),
		)
		if err != nil {
			logrus.WithContext(ctx).
				WithError(err).
				Warn("failed to initialize gitlab client in license plugin, skipping license retrieval")
			return false
		}

		plugin.GitlabClient = client
		return true
	}
	return false
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *license) Execute(ctx context.Context, config models.GenerateConfig, _ filesystem.FS) error {
	dest := filepath.Join(config.Options.DestinationDir, models.License)

	// don't fetch template is force on file or force all isn't activated
	if !config.Options.ForceAll && filesystem.Exists(dest) && !slices.Contains(config.Options.Force, models.License) {
		return nil
	}

	// fetch license template
	license, _, err := plugin.GitlabClient.LicenseTemplates.GetLicenseTemplate(*config.License, &gitlab.GetLicenseTemplateOptions{
		Fullname: &config.Maintainers[0].Name,
		Project:  &config.ProjectName,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to retrieve license from gitlab: %w", err)
	}

	// write license template
	if err := os.WriteFile(filepath.Join(config.Options.DestinationDir, models.License), []byte(license.Content), filesystem.RwRR); err != nil {
		return fmt.Errorf("failed to write license: %w", err)
	}
	return nil
}

// Name returns the plugin name.
func (*license) Name() string {
	return "license"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*license) Remove(_ context.Context, config models.GenerateConfig) error {
	if err := os.Remove(filepath.Join(config.Options.DestinationDir, models.License)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove LICENSE file: %w", err)
	}
	return nil
}

// Type returns the type of given plugin.
func (*license) Type() pluginType {
	return secondary
}
