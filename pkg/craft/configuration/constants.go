package craft

// File is the craft configuration file name.
const File = ".craft"

const (
	// Chart can be given in craft exclusions ('exclude' key) to avoid generating an Helm chart.
	//
	// By default, a Helm chart is generated since a project could satisfies one of the following possibilities:
	//	- the project is just an Helm chart for another product
	// 	- the project is a product with an Helm chart (with one or multiple resources, cronjobs, job, worker, etc)
	Chart string = "chart"

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
