package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestValidatePackageJSON(t *testing.T) {
	t.Run("error_join_all", func(t *testing.T) {
		// Arrange
		jsonfile := parser.PackageJSON{}

		// Act
		err := jsonfile.Validate()

		// Assert
		assert.ErrorIs(t, err, parser.ErrMissingPackageName)
		assert.ErrorIs(t, err, parser.ErrInvalidPackageManager)
	})

	t.Run("error_invalid_package_manager", func(t *testing.T) {
		// Arrange
		jsonfile := parser.PackageJSON{
			Name:           "test",
			PackageManager: "invalid",
		}

		// Act
		err := jsonfile.Validate()

		// Assert
		assert.ErrorIs(t, err, parser.ErrInvalidPackageManager)
	})

	t.Run("success_valid_file", func(t *testing.T) {
		// Arrange
		jsonfile := parser.PackageJSON{
			Name:           "test",
			PackageManager: "npm@9.0.0",
		}

		// Act
		err := jsonfile.Validate()

		// Assert
		assert.NoError(t, err)
	})
}
