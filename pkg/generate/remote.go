package generate

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// OriginURL returns input directory git config --get remote.origin.url.
func OriginURL(destdir string) (out []byte, err error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = destdir

	out, err = cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("failed to retrieve remote url: %w", err)
	}

	// we want to return the real output alongside the error if there was one, that's why return properties are named
	// by doing that we're allowing the potential print of the terminal output in case there was an error, to potentially get more information from git
	return
}

// ParseRemote returns the current repository host and path to repository on the given host's platform.
func ParseRemote(rawRemote string) (host, path string) {
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
	return lo.FindKeyBy(map[string][]string{
		craft.Bitbucket: {"bb", craft.Bitbucket, "stash"},
		craft.Gitea:     {craft.Gitea},
		craft.Github:    {craft.Github, "gh"},
		craft.Gitlab:    {craft.Gitlab, "gl"},
	}, func(_ string, searches []string) bool {
		return slices.ContainsFunc(searches, func(search string) bool { return strings.Contains(host, search) })
	})
}
