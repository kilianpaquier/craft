package craft

const (
	// File is the craft configuration file name.
	File = ".craft"

	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// Gocmd represents the cmd folder where go main.go should be placed according to go layout.
	Gocmd = "cmd"
	// Gomod represents the go.mod filename.
	Gomod = "go.mod"
	// PackageJSON represents the package.json filename.
	PackageJSON = "package.json"

	// License represents the target filename for the generated project LICENSE.
	License = "LICENSE"
)

const (
	// Bitbucket is just the bitbucket constant.
	Bitbucket = "bitbucket"
	// Gitea is just the gitea constant.
	Gitea = "gitea"
	// Github is just the github constant.
	Github = "github"
	// Gitlab is just the gitlab constant.
	Gitlab = "gitlab"
)

const (
	// CodeCov is the codecov option for CI tuning.
	CodeCov string = "codecov"
	// CodeQL is the codeql option for CI tuning.
	CodeQL string = "codeql"
	// Dependabot is the dependabot option for CI tuning.
	Dependabot string = "dependabot"
	// Pages is the pages option for CI tuning.
	Pages string = "pages"
	// Renovate is the renovate option for CI tuning.
	Renovate string = "renovate"
	// Sonar is the sonar option for CI tuning.
	Sonar string = "sonar"
)

const (
	// GithubApps is the value for github release mode with a github app.
	GithubApps string = "github-apps"
	// GithubToken is the value for github release mode with a github token.
	GithubToken string = "github-token"
	// PersonalToken is the value for github release mode with a personal token (PAT).
	PersonalToken string = "personal-token"
)

// CIOptions returns the slice with all availables CI options.
func CIOptions() []string {
	return []string{CodeCov, CodeQL, Dependabot, Pages, Renovate, Sonar}
}
