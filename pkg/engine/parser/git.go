package parser

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

// Git reads the input destdir directory remote.origin.url
// to retrieve various project information (git host, project name, etc.).
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func Git(ctx context.Context, destdir string, c *config) error {
//		vcs, err := parser.Git(destdir)
//		if err != nil {
//			engine.GetLogger(ctx).Warnf("failed to retrieve git vcs configuration: %v", err)
//			return nil // a repository may not be a git repository
//		}
//		engine.GetLogger(ctx).Infof("git repository detected")
//		// do something with vcs (e.g. update config since it's a pointer)
//		return nil
//	}
func Git(destdir string) (VCS, error) {
	rawRemote, err := gitOriginURL(destdir)
	if err != nil {
		return VCS{}, fmt.Errorf("git origin URL: %w", err)
	}
	host, subpath := gitParseRemote(rawRemote)
	platform, _ := parsePlatform(host)
	return VCS{
		ProjectHost: host,
		ProjectName: path.Base(subpath),
		ProjectPath: subpath,
		Platform:    platform,
	}, nil
}

// gitOriginURL returns input directory git config --get remote.origin.url.
func gitOriginURL(destdir string) (string, error) {
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

// gitParseRemote returns the current repository host and path to repository on the given host's platform.
func gitParseRemote(rawRemote string) (_, _ string) {
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
