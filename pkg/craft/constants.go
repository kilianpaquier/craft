package craft

const (
	// File is the craft configuration file name.
	File = ".craft"

	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// PartExtension is the extension for templates files' subparts.
	//
	// It must be used with TmplExtension
	// and as such files with only templates parts (define) can be created.
	PartExtension = ".part"

	// PatchExtension is the extension for templates files patches.
	//
	// It will be used in the future to patch altered files by users to follow updates with less generation issues.
	PatchExtension = ".patch"

	// Gocmd represents the cmd folder where go main.go should be placed according to go layout.
	Gocmd = "cmd"
	// Gomod represents the go.mod filename.
	Gomod = "go.mod"
	// PackageJSON represents package.json filename.
	PackageJSON = "package.json"

	// License represents the target filename for the generated project LICENSE.
	License = "LICENSE"
)

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

const (
	// CodeCov is the codecov option for CI tuning.
	CodeCov string = "codecov"
	// CodeQL is the codeql option for CI tuning.
	CodeQL string = "codeql"
	// Labeler is the auto labeling option for CI tuning.
	Labeler string = "labeler"
	// Sonar is the sonar option for CI tuning.
	Sonar string = "sonar"
)

const (
	// Netlify is the static name to deploy on netlify (only available on github actions).
	Netlify string = "netlify"
	// Pages is the static name for pages deployment.
	Pages string = "pages"
)

const (
	// Dependabot is the dependabot updater name for CI maintenance configuration.
	Dependabot string = "dependabot"
	// Renovate is the renovate updater name for CI maintenance configuration.
	Renovate string = "renovate"
)

const (
	// GitHubApp is the value for github release mode with a github app.
	GitHubApp string = "github-app"
	// GitHubToken is the value for github release mode with a github token.
	GitHubToken string = "github-token"
	// PersonalToken is the value for github release mode with a personal token (PAT).
	PersonalToken string = "personal-token"
)

const (
	// Mendio is the value for maintenance mode with renovate and mend.io (meaning no self-hosted renovate).
	Mendio string = "mend.io"
)

// SemanticRelease is the value for github / gitlab release with semantic-release.
const SemanticRelease string = "semantic-release"

// CIOptions returns the slice with all availables CI options.
func CIOptions() []string {
	return []string{CodeCov, CodeQL, Labeler, Sonar}
}
