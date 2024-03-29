package generate

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/internal/models"
)

// getRemoteURL returns current directory execution git config remote.origin.url.
func getRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve remote url: %w", err)
	}
	return string(out), nil
}

// parseRemote returns the current repository host and path to repository on the given host's platform.
func parseRemote(rawRemote string) (host, path string) {
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
	return lo.FindKeyBy(map[string][]string{
		models.Bitbucket: {"bb", models.Bitbucket, "stash"},
		models.Gitea:     {models.Gitea},
		models.Github:    {models.Github, "gh"},
		models.Gitlab:    {models.Gitlab, "gl"},
	}, func(_ string, searches []string) bool {
		return slices.ContainsFunc(searches, func(search string) bool { return strings.Contains(host, search) })
	})
}
