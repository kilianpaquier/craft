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
	// Netlify is the option to deploy on netlify.
	Netlify string = "netlify"
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

const (
	// GhRelease is the value for github release with gh-release, yet another github action to create github releases.
	//
	// With this action, in any case, the token used is the github-token mode.
	GhRelease string = "gh-release"

	// ReleaseDrafter is the value for github release with release-drafter instead of semantic-release.
	//
	// With this action, in any case, the token used is the github-token mode.
	ReleaseDrafter string = "release-drafter"

	// SemanticRelease is the value for github / gitlab release with semantic-release.
	SemanticRelease string = "semantic-release"
)

// CIOptions returns the slice with all availables CI options.
func CIOptions() []string {
	return []string{CodeCov, CodeQL, Dependabot, Netlify, Pages, Renovate, Sonar}
}
