// Code generated by go-builder-generator (https://github.com/kilianpaquier/go-builder-generator). DO NOT EDIT.

package tests

import (
	"github.com/kilianpaquier/craft/internal/models"
)

// GenerateOptionsBuilder represents GenerateOptions's builder.
type GenerateOptionsBuilder struct {
	build models.GenerateOptions
}

// NewGenerateOptionsBuilder creates a new GenerateOptionsBuilder.
func NewGenerateOptionsBuilder() *GenerateOptionsBuilder {
	return &GenerateOptionsBuilder{}
}

// Copy reassigns the builder struct (behind pointer) to a new pointer and returns it.
func (b *GenerateOptionsBuilder) Copy() *GenerateOptionsBuilder {
	return &GenerateOptionsBuilder{b.build}
}

// Build returns built GenerateOptions.
func (b *GenerateOptionsBuilder) Build() *models.GenerateOptions {
	result := b.build
	return &result
}

// DestinationDir sets GenerateOptions's DestinationDir.
func (b *GenerateOptionsBuilder) DestinationDir(destinationDir string) *GenerateOptionsBuilder {
	b.build.DestinationDir = destinationDir
	return b
}

// EndDelim sets GenerateOptions's EndDelim.
func (b *GenerateOptionsBuilder) EndDelim(endDelim string) *GenerateOptionsBuilder {
	b.build.EndDelim = endDelim
	return b
}

// Force sets GenerateOptions's Force.
func (b *GenerateOptionsBuilder) Force(force ...string) *GenerateOptionsBuilder {
	b.build.Force = append(b.build.Force, force...)
	return b
}

// ForceAll sets GenerateOptions's ForceAll.
func (b *GenerateOptionsBuilder) ForceAll(forceAll bool) *GenerateOptionsBuilder {
	b.build.ForceAll = forceAll
	return b
}

// StartDelim sets GenerateOptions's StartDelim.
func (b *GenerateOptionsBuilder) StartDelim(startDelim string) *GenerateOptionsBuilder {
	b.build.StartDelim = startDelim
	return b
}

// TemplatesDir sets GenerateOptions's TemplatesDir.
func (b *GenerateOptionsBuilder) TemplatesDir(templatesDir string) *GenerateOptionsBuilder {
	b.build.TemplatesDir = templatesDir
	return b
}
