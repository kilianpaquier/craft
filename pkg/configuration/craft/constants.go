package craft

// File is the craft configuration file name.
const File = ".craft"

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
