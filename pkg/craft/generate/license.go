package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-cleanhttp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserLicense generates the license file for the project.
func ParserLicense(ctx context.Context, destdir string, config *craft.Config) error {
	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(parser.GitLabURL),
		gitlab.WithHTTPClient(cleanhttp.DefaultClient()),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)))
	if err != nil {
		// should never happen since it's gitlab.ClientOptionFunc that are throwing errors
		// and currently WithBaseURL with fixed URL
		// and WithoutRetries won't throw errors
		// but in any case err must be handled in case it evolves or other options are added
		engine.GetLogger(ctx).Warnf("failed to initialize GitLab client, skipping license generation: %v", err)
		return nil
	}

	dest := filepath.Join(destdir, parser.FileLicense)
	if config.License == nil {
		engine.GetLogger(ctx).Infof("skipping license generation, configuration doesn't have 'license' key")
		if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove '%s': %w", parser.FileLicense, err)
		}
		return nil
	}

	if files.Exists(dest) {
		engine.GetLogger(ctx).Infof("not generating '%s' since it already exists", parser.FileLicense)
		return nil
	}
	engine.GetLogger(ctx).Infof("license detected, configuration has 'license' key")

	opts := parser.LicenseOptions{
		Client:     client,
		License:    *config.License,
		Maintainer: &config.Maintainers[0].Name,
		Project:    &config.ProjectName,
	}
	return parser.DownloadLicense(dest, opts)
}

var _ engine.Parser[craft.Config] = ParserLicense // ensure interface is implemented
