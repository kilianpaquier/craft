package parser

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Git reads the input destdir directory remote.origin.url to retrieve various project information (git host, project name, etc.).
func Git(ctx context.Context, destdir string, config *craft.Config) error {
	rawRemote, err := originURL(destdir)
	if err != nil {
		generate.GetLogger(ctx).Warnf("failed to retrieve git remote.origin.url: %s", err.Error())
		return nil
	}
	generate.GetLogger(ctx).Infof("git repository detected")

	host, subpath := parseRemote(rawRemote)
	if config.Platform == "" {
		config.Platform, _ = parsePlatform(host)
	}

	config.ProjectHost = host
	config.ProjectName = path.Base(subpath)
	config.ProjectPath = subpath
	return nil
}

var _ generate.Parser[craft.Config] = Git // ensure interface is implemented

// originURL returns input directory git config --get remote.origin.url.
func originURL(destdir string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = destdir

	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("retrieve remote url with response '%s': %w", string(out), err)
		}
		return "", fmt.Errorf("retrieve remote url: %w", err)
	}
	return string(out), nil
}

// parseRemote returns the current repository host and path to repository on the given host's platform.
func parseRemote(rawRemote string) (_, _ string) {
	if rawRemote == "" {
		return "", ""
	}

	originURL := strings.TrimSuffix(rawRemote, "\n")
	originURL = strings.TrimSuffix(originURL, ".git")

	// handle ssh remotes
	if strings.HasPrefix(originURL, "git@") {
		originURL := strings.TrimPrefix(originURL, "git@")
		host, subpath, _ := strings.Cut(originURL, ":")
		return host, subpath
	}

	// handle web url remotes
	originURL = strings.TrimPrefix(originURL, "http://")
	originURL = strings.TrimPrefix(originURL, "https://")
	host, subpath, _ := strings.Cut(originURL, "/")
	return host, subpath
}

// parsePlatform returns the platform name associated to input host.
func parsePlatform(host string) (string, bool) {
	matchers := map[string][]string{
		craft.Bitbucket: {"bb", craft.Bitbucket, "stash"},
		craft.Gitea:     {craft.Gitea},
		craft.GitHub:    {craft.GitHub, "gh"},
		craft.GitLab:    {craft.GitLab, "gl"},
	}

	for platform, candidates := range matchers {
		for _, candidate := range candidates {
			if strings.Contains(host, candidate) {
				return platform, true
			}
		}
	}
	return "", false
}
