package parser

// Executables represents a collection of executables (clis, cronjobs, jobs, workers).
type Executables struct {
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
func (g Executables) Binaries() int {
	return len(g.Clis) + len(g.Crons) + len(g.Jobs) + len(g.Workers)
}

// SetCLI sets a CLI with its name.
func (g *Executables) SetCLI(name string) {
	if g.Clis == nil {
		g.Clis = map[string]struct{}{}
	}
	g.Clis[name] = struct{}{}
}

// SetCron sets a cronjob with its name.
func (g *Executables) SetCron(name string) {
	if g.Crons == nil {
		g.Crons = map[string]struct{}{}
	}
	g.Crons[name] = struct{}{}
}

// SetJob sets a job with its name.
func (g *Executables) SetJob(name string) {
	if g.Jobs == nil {
		g.Jobs = map[string]struct{}{}
	}
	g.Jobs[name] = struct{}{}
}

// SetWorker sets a worker with its name.
func (g *Executables) SetWorker(name string) {
	if g.Workers == nil {
		g.Workers = map[string]struct{}{}
	}
	g.Workers[name] = struct{}{}
}
