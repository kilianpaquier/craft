// Code generated by go-builder-generator (https://github.com/kilianpaquier/go-builder-generator). DO NOT EDIT.

package tests

import (
	"github.com/kilianpaquier/craft/internal/models"
)

// CIBuilder represents CI's builder.
type CIBuilder struct {
	build models.CI
}

// NewCIBuilder creates a new CIBuilder.
func NewCIBuilder() *CIBuilder {
	return &CIBuilder{}
}

// Copy reassigns the builder struct (behind pointer) to a new pointer and returns it.
func (b *CIBuilder) Copy() *CIBuilder {
	return &CIBuilder{b.build}
}

// Build returns built CI.
func (b *CIBuilder) Build() *models.CI {
	result := b.build
	return &result
}

// AutoRelease sets CI's AutoRelease.
func (b *CIBuilder) AutoRelease(autoRelease bool) *CIBuilder {
	b.build.AutoRelease = autoRelease
	return b
}

// Name sets CI's Name.
func (b *CIBuilder) Name(name string) *CIBuilder {
	b.build.Name = name
	return b
}

// Options sets CI's Options.
func (b *CIBuilder) Options(options ...string) *CIBuilder {
	b.build.Options = append(b.build.Options, options...)
	return b
}
