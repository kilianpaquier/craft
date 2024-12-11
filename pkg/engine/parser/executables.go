package parser

// Executables represents a collection of executables (clis, cronjobs, jobs, workers).
type Executables struct {
	// Clis is a map of CLI names without value (empty struct).
	Clis map[string]struct{} `json:"clis,omitempty"`

	// Crons is a map of cronjob names without value (empty struct).
	Crons map[string]struct{} `json:"crons,omitempty"`

	// Jobs is a map of job names without value (empty struct).
	Jobs map[string]struct{} `json:"jobs,omitempty"`

	// Workers is a map of workers names without value (empty struct).
	Workers map[string]struct{} `json:"workers,omitempty"`
}

// Binaries returns the sum of all executables (clis, cronjobs, jobs, workers).
func (g Executables) Binaries() int {
	return len(g.Clis) + len(g.Crons) + len(g.Jobs) + len(g.Workers)
}

// AddCLI sets a CLI with its name.
// In case a CLI with the name same already exists, it is replaced.
func (g *Executables) AddCLI(name string) {
	if g.Clis == nil {
		g.Clis = map[string]struct{}{}
	}
	g.Clis[name] = struct{}{}
}

// AddCron sets a cronjob with its name.
// In case a cronjob with the name same already exists, it is replaced.
func (g *Executables) AddCron(name string) {
	if g.Crons == nil {
		g.Crons = map[string]struct{}{}
	}
	g.Crons[name] = struct{}{}
}

// AddJob sets a job with its name.
// In case a job with the name same already exists, it is replaced.
func (g *Executables) AddJob(name string) {
	if g.Jobs == nil {
		g.Jobs = map[string]struct{}{}
	}
	g.Jobs[name] = struct{}{}
}

// AddWorker sets a worker with its name.
// In case a worker with the name same already exists, it is replaced.
func (g *Executables) AddWorker(name string) {
	if g.Workers == nil {
		g.Workers = map[string]struct{}{}
	}
	g.Workers[name] = struct{}{}
}
