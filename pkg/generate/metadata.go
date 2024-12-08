package generate

import "github.com/kilianpaquier/craft/pkg/craft"

// Metadata represents all properties available for enrichment during repository parsing.
//
// Updated properties will be used during generation to determine if a specific file or part of a file must be generated.
type Metadata struct {
	craft.Configuration

	// Languages is a map of languages name with its specificities.
	Languages map[string]any `json:"-"`

	// ProjectHost represents the host where the project is hosted.
	//
	// As craft only handles git, it would be an host like github.com, gitlab.com, bitbucket.org, etc.
	//
	// Of course it can also be a private host like github.company.com.
	//
	// It will depend on the git origin URL or for golang the host of go.mod module name.
	ProjectHost string `json:"projectHost,omitempty"`

	// ProjectName is the project name being generated.
	// By default with Run function, it will be the base path of ParseRemote's subpath result following OriginURL result.
	ProjectName string `json:"projectName,omitempty"`

	// ProjectPath is the project path.
	// By default with Run function, it will be the subpath in ParseRemote result.
	ProjectPath string `json:"projectPath,omitempty"`

	// Binaries is the total number of binaries / executables parsed during craft execution.
	//
	// It's especially used for golang generation (with workers, cronjob, jobs, etc.)
	// but also in nodejs generation in case a "main" property is present in package.json.
	Binaries uint8 `json:"-"`

	// Clis is a map of CLI names without value (empty struct).
	Clis map[string]struct{} `json:"-"`

	// Crons is a map of cronjob names without value (empty struct).
	Crons map[string]struct{} `json:"crons,omitempty"`

	// Jobs is a map of job names without value (empty struct).
	Jobs map[string]struct{} `json:"jobs,omitempty"`

	// Workers is a map of workers names without value (empty struct).
	Workers map[string]struct{} `json:"workers,omitempty"`
}
