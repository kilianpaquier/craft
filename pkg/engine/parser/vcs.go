package parser

import "strings"

const (
	// Bitbucket is just the bitbucket constant.
	Bitbucket = "bitbucket"
	// Gitea is just the gitea constant.
	Gitea = "gitea"
	// GitHub is just the github constant.
	GitHub = "github"
	// GitLab is just the gitlab constant.
	GitLab = "gitlab"
)

// VCS structure contains all properties related to VCS (Version Control System).
type VCS struct {
	// Name is the version control system name.
	Name string `json:"name,omitempty" yaml:"-"`

	// Platform represents the vcs platform hosting the project.
	Platform string `json:"-" yaml:"platform,omitempty"`

	// ProjectHost represents the host where the project is hosted.
	//
	// It will depend on the vcs origin URL
	// or for golang the host of go.mod module name.
	ProjectHost string `json:"projectHost,omitempty" yaml:"-"`

	// ProjectName is the project name being generated.
	// By default with Run function, it will be the base path of ParseRemote's subpath result following OriginURL result.
	ProjectName string `json:"projectName,omitempty" yaml:"-"`

	// ProjectPath is the project path.
	// By default with Run function, it will be the subpath in ParseRemote result.
	ProjectPath string `json:"projectPath,omitempty" yaml:"-"`
}

// parsePlatform returns the platform name associated to input host.
func parsePlatform(host string) (string, bool) {
	matchers := map[string][]string{
		Bitbucket: {"bb", Bitbucket, "stash"},
		Gitea:     {Gitea},
		GitHub:    {GitHub, "gh"},
		GitLab:    {GitLab, "gl"},
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
