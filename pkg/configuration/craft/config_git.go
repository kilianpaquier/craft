package craft

// GitConfig structure contains all properties related
// to how git configuration is read in a given project.
//
// It's used in a composition with craft.Configuration.
//
// However, in case of a custom configuration (since generate and initialize packages can handle this),
// GitConfig can be reused in another composition.
type GitConfig struct {
	// Platform represents the git platform hosting the project.
	//
	// On the first generation run (with parser.Git), it will be set.
	// However, it's possible to override it manually in the .craft file.
	Platform string `json:"-" yaml:"platform,omitempty"`

	// ProjectHost represents the host where the project is hosted.
	//
	// As craft only handles git, it would be an host like github.com, gitlab.com, bitbucket.org, etc.
	//
	// Of course it can also be a private host like github.company.com.
	//
	// It will depend on the git origin URL or for golang the host of go.mod module name.
	ProjectHost string `json:"projectHost,omitempty" yaml:"-"`

	// ProjectName is the project name being generated.
	// By default with Run function, it will be the base path of ParseRemote's subpath result following OriginURL result.
	ProjectName string `json:"projectName,omitempty" yaml:"-"`

	// ProjectPath is the project path.
	// By default with Run function, it will be the subpath in ParseRemote result.
	ProjectPath string `json:"projectPath,omitempty" yaml:"-"`
}
