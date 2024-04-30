// Code generated by go-builder-generator (https://github.com/kilianpaquier/go-builder-generator). DO NOT EDIT.

package detectgen

import (
	"fmt"

	"github.com/kilianpaquier/craft/internal/generate/filehandler"
	"github.com/kilianpaquier/craft/internal/models"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
)

// DirGenerateBuilder represents DirGenerate's builder.
type DirGenerateBuilder struct {
	build DirGenerate
}

// NewDirGenerateBuilder creates a new DirGenerateBuilder.
func NewDirGenerateBuilder() *DirGenerateBuilder {
	return &DirGenerateBuilder{}
}

// Copy reassigns the builder struct (behind pointer) to a new pointer and returns it.
func (b *DirGenerateBuilder) Copy() *DirGenerateBuilder {
	return &DirGenerateBuilder{b.build}
}

// Build returns built DirGenerate.
func (b *DirGenerateBuilder) Build() (*DirGenerate, error) {
	result := b.build
	if err := result.Validate(); err != nil {
		return nil, fmt.Errorf("validation of 'DirGenerate''s struct: %w", err)
	}
	return &result, nil
}

// SetConfig sets DirGenerate's Config.
func (b *DirGenerateBuilder) SetConfig(config models.GenerateConfig) *DirGenerateBuilder {
	b.build.Config = config
	return b
}

// SetData sets DirGenerate's Data.
func (b *DirGenerateBuilder) SetData(data any) *DirGenerateBuilder {
	b.build.Data = data
	return b
}

// SetFileHandlers sets DirGenerate's FileHandlers.
func (b *DirGenerateBuilder) SetFileHandlers(fileHandlers []filehandler.Handler) *DirGenerateBuilder {
	b.build.FileHandlers = fileHandlers
	return b
}

// SetFS sets DirGenerate's FS.
func (b *DirGenerateBuilder) SetFS(fs filesystem.FS) *DirGenerateBuilder {
	b.build.FS = fs
	return b
}

// SetName sets DirGenerate's Name.
func (b *DirGenerateBuilder) SetName(name GenerateName) *DirGenerateBuilder {
	b.build.Name = name
	return b
}