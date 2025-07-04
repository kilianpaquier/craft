package craft

// File is the craft configuration file name.
const File = ".craft"

const (
	// Goreleaser can be given in craft exclusions ('exclude' key) to avoid generating a .goreleaser.yml file.
	//
	// By default, if a given project is a Go project,
	// and a cmd CLI is defined (cmd/<some useful CLI name>)
	// a .goreleaser.yml file is generated.
	//
	// As such, it's unnecessary to specify this property when your project isn't a Go one.
	Goreleaser string = "goreleaser"

	// Makefile can be given in craft exclusions ('exclude' key) to avoid generating a Makefile
	// and additional Makefiles in scripts/mk/*.mk.
	//
	// When crafting a Node project, it's unnecessary to specify this property since no Makefile will be generated anyway.
	// It's because Node projects contain all their scripts in package.json.
	Makefile string = "makefile"

	// PreCommit can be given in craft exclusions ('exclude' key) to avoid generating pre-commit files and associated Continuous Integration.
	PreCommit string = "pre-commit"

	// Shell can be given in craft exclusions ('exclude' key)
	// to avoid generating shell (check / test / pre-commit) Continuous Integration.
	Shell string = "shell"
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
	// Kubernetes is the deployment name for kubernetes
	// (doesn't matter the provider, kubernetes has its own communication interface - i.e. kubectl and/or helm).
	Kubernetes string = "kubernetes"
	// Netlify is the deployment name for netlify.
	Netlify string = "netlify"
	// Pages is the deployment name for pages (GitLab or GitHub depending on CI name).
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
	// HelmAuto is the constant indicating that Helm chart publication should be made automatically.
	HelmAuto string = "auto"
	// HelmManual is the constant indicating that Helm chart publication should be made manually.
	HelmManual string = "manual"
	// HelmNone is the constant indicating that Helm chart publication should not be made.
	HelmNone string = "none"
)
