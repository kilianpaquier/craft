// Code generated by go-builder-generator (https://github.com/kilianpaquier/go-builder-generator). DO NOT EDIT.

package tests

import (
	"github.com/kilianpaquier/craft/internal/models"
)

// GenerateConfigBuilder is an alias of GenerateConfig to build GenerateConfig with builder-pattern.
type GenerateConfigBuilder models.GenerateConfig

// NewGenerateConfigBuilder creates a new GenerateConfigBuilder.
func NewGenerateConfigBuilder() *GenerateConfigBuilder {
	return &GenerateConfigBuilder{}
}

// Copy reassigns the builder struct (behind pointer) to a new pointer and returns it.
func (b *GenerateConfigBuilder) Copy() *GenerateConfigBuilder {
	c := *b
	return &c
}

// Build returns built GenerateConfig.
func (b *GenerateConfigBuilder) Build() *models.GenerateConfig {
	c := (models.GenerateConfig)(*b)
	return &c
}

// SetClis sets GenerateConfig's Clis.
func (b *GenerateConfigBuilder) SetClis(clis map[string]struct{}) *GenerateConfigBuilder {
	b.Clis = clis
	return b
}

// SetCraftConfig sets GenerateConfig's CraftConfig.
func (b *GenerateConfigBuilder) SetCraftConfig(craftConfig models.CraftConfig) *GenerateConfigBuilder {
	b.CraftConfig = craftConfig
	return b
}

// SetCrons sets GenerateConfig's Crons.
func (b *GenerateConfigBuilder) SetCrons(crons map[string]struct{}) *GenerateConfigBuilder {
	b.Crons = crons
	return b
}

// SetJobs sets GenerateConfig's Jobs.
func (b *GenerateConfigBuilder) SetJobs(jobs map[string]struct{}) *GenerateConfigBuilder {
	b.Jobs = jobs
	return b
}

// SetModuleName sets GenerateConfig's ModuleName.
func (b *GenerateConfigBuilder) SetModuleName(moduleName string) *GenerateConfigBuilder {
	b.ModuleName = moduleName
	return b
}

// SetModuleVersion sets GenerateConfig's ModuleVersion.
func (b *GenerateConfigBuilder) SetModuleVersion(moduleVersion string) *GenerateConfigBuilder {
	b.ModuleVersion = moduleVersion
	return b
}

// SetOptions sets GenerateConfig's Options.
func (b *GenerateConfigBuilder) SetOptions(options models.GenerateOptions) *GenerateConfigBuilder {
	b.Options = options
	return b
}

// SetProjectName sets GenerateConfig's ProjectName.
func (b *GenerateConfigBuilder) SetProjectName(projectName string) *GenerateConfigBuilder {
	b.ProjectName = projectName
	return b
}

// SetWorkers sets GenerateConfig's Workers.
func (b *GenerateConfigBuilder) SetWorkers(workers map[string]struct{}) *GenerateConfigBuilder {
	b.Workers = workers
	return b
}
