package generate

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// OriginURL returns input directory git config --get remote.origin.url.
func OriginURL(destdir string) (string, error) {
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

// ParseRemote returns the current repository host and path to repository on the given host's platform.
func ParseRemote(rawRemote string) (host, path string) {
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

// ParsePlatform returns the platform name associated to input host.
func ParsePlatform(host string) (string, bool) {
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
