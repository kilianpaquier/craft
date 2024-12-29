package craft

// ConfigFiles structure contains all properties related
// to files layout in a given project (languages, clis, cronjobs, jobs, workers, etc.).
//
// It's used in a composition with craft.Configuration.
//
// However, in case of a custom configuration (since generate and initialize packages can handle this),
// ConfigFiles can be reused in another composition.
type ConfigFiles struct {
	// Languages is a map of languages name with its specificities.
	Languages map[string]any `json:"-" yaml:"-"`

	// Clis is a map of CLI names without value (empty struct).
	Clis map[string]struct{} `json:"-" yaml:"-"`

	// Crons is a map of cronjob names without value (empty struct).
	Crons map[string]struct{} `json:"crons,omitempty" yaml:"-"`

	// Jobs is a map of job names without value (empty struct).
	Jobs map[string]struct{} `json:"jobs,omitempty" yaml:"-"`

	// Workers is a map of workers names without value (empty struct).
	Workers map[string]struct{} `json:"workers,omitempty" yaml:"-"`
}

// Binaries returns the sum of all executables (clis, cronjobs, jobs, workers).
func (g ConfigFiles) Binaries() int {
	return len(g.Clis) + len(g.Crons) + len(g.Jobs) + len(g.Workers)
}

// SetLanguage sets a language with its specificities.
func (g *ConfigFiles) SetLanguage(name string, value any) {
	if g.Languages == nil {
		g.Languages = map[string]any{}
	}
	g.Languages[name] = value
}

// SetCLI sets a CLI with its name.
func (g *ConfigFiles) SetCLI(name string) {
	if g.Clis == nil {
		g.Clis = map[string]struct{}{}
	}
	g.Clis[name] = struct{}{}
}

// SetCron sets a cronjob with its name.
func (g *ConfigFiles) SetCron(name string) {
	if g.Crons == nil {
		g.Crons = map[string]struct{}{}
	}
	g.Crons[name] = struct{}{}
}

// SetJob sets a job with its name.
func (g *ConfigFiles) SetJob(name string) {
	if g.Jobs == nil {
		g.Jobs = map[string]struct{}{}
	}
	g.Jobs[name] = struct{}{}
}

// SetWorker sets a worker with its name.
func (g *ConfigFiles) SetWorker(name string) {
	if g.Workers == nil {
		g.Workers = map[string]struct{}{}
	}
	g.Workers[name] = struct{}{}
}
