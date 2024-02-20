package models

const (
	// CraftFile is the craft configuration file name.
	CraftFile = ".craft"

	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// ShellExtension is the extension for executable files.
	ShellExtension = ".sh"

	// Dockerfile represents the target filename for the generated project Dockerfile.
	Dockerfile = "Dockerfile"
	// Dockerignore represents the target filename for the generated project .dockerignore.
	Dockerignore = ".dockerignore"
	// Github is just the github constant.
	Github = "github"
	// GithubCI represents the folder name for github repository configuration files.
	GithubCI = ".github"
	// GithubWorkflows represents the folder name for github actions files.
	GithubWorkflows = "workflows"
	// Gitignore represents the target filename for the generated project .gitignore.
	Gitignore = ".gitignore"
	// Gitlab is just the gitlab constant.
	Gitlab = "gitlab"
	// GitlabCI represents the target filename for the generated project .gitlab-ci.yml.
	GitlabCI = ".gitlab-ci.yml"
	// GoCmd represents the cmd folder where go main.go should be placed according to go layout.
	GoCmd = "cmd"
	// GolangCI represents the target filename for the generated project golangci-lint.
	GolangCI = ".golangci.yml"
	// GoMain represents the common main filename.
	GoMain = "main.go"
	// GoMod represents the go.mod filename.
	GoMod = "go.mod"
	// Goreleaser represents the target filename for the generated project goreleaser.yml.
	Goreleaser = ".goreleaser.yml"
	// Launcher represents the target filename for the generated project launcher.sh.
	Launcher = "launcher.sh"
	// License represents the target filename for the generated project LICENSE.
	License = "LICENSE"
	// Makefile represents the target filename for the generated project Makefile.
	Makefile = "Makefile"
	// PythonInit represents the init filename used by python to define a folder as a module.
	PythonInit = "__init__.py"
	// PythonMain represents the main filename (commonly) used in python development to identify a module with a main function.
	PythonMain = "__main__.py"
	// PythonRequirements represents the filename for python projects requirements.
	PythonRequirements = "requirements.txt"
	// Readme represents the target filename for the generated project README.md.
	Readme = "README.md"
	// Releaserc represents the target filename for semantic release configuration file.
	Releaserc = ".releaserc.yml"
	// SonarProperties represents the target filename for the generated project sonar properties.
	SonarProperties = "sonar.properties"
	// SwaggerFile represents the target filename for the generated project api.yml.
	SwaggerFile = "api.yml"

	// ChartDir represents the directory name where helm chart will be initialized and maintained.
	ChartDir = "chart"
	// ChartFile represents the filename of helm chart file.
	ChartFile = "Chart.yaml"
	// ValuesFile represents the filename of helm default values file.
	ValuesFile = "values.yaml"
	// Helpers represents the filename of helm templates helper file.
	Helpers = "_helpers.tpl"
)
