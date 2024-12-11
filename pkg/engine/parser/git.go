package parser

import (
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Git reads the input destdir directory remote.origin.url
// to retrieve various project information (git host, project name, etc.).
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func ParserGit(ctx context.Context, destdir string, c *config) error {
//		vcs, err := parser.Git(destdir)
//		if err != nil {
//			engine.GetLogger().Warnf("failed to retrieve git vcs configuration: %v", err)
//			return nil // a repository may not be a git repository
//		}
//		engine.GetLogger().Infof("git repository detected")
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

// gitOriginURL returns input directory remote origin by using go-git.
//
// Two errors may be checked with errors.Is when one is returned:
//   - git.ErrRepositoryNotExists
//   - git.ErrRemoteNotFound
func gitOriginURL(destdir string) (string, error) {
	repository, err := git.PlainOpenWithOptions(destdir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", fmt.Errorf("open repository: %w", err)
	}
	origin, err := repository.Remote("origin")
	if err != nil {
		return "", fmt.Errorf("get remote 'origin': %w", err)
	}
	if len(origin.Config().URLs) == 0 {
		return "", fmt.Errorf("no URL associated to remote: %w", git.ErrRemoteNotFound)
	}
	return origin.Config().URLs[0], nil
}

// gitParseRemote returns the current repository host and path to repository on the given host's platform.
func gitParseRemote(rawRemote string) (_, _ string) {
	if rawRemote == "" {
		return "", ""
	}

	originURL := strings.TrimSuffix(rawRemote, "\n")
	originURL = strings.TrimSuffix(originURL, ".git")

	// handle ssh remotes
	if after, ok := strings.CutPrefix(originURL, "git@"); ok {
		originURL := after
		host, subpath, _ := strings.Cut(originURL, ":")
		return host, subpath
	}

	// handle web url remotes
	originURL = strings.TrimPrefix(originURL, "http://")
	originURL = strings.TrimPrefix(originURL, "https://")
	host, subpath, _ := strings.Cut(originURL, "/")
	return host, subpath
}
