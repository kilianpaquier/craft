package craft_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	schemas "github.com/kilianpaquier/craft/.schemas"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestValidate(t *testing.T) {
	readSchema := func(out any) error { return files.ReadJSON(schemas.Craft, out, schemas.ReadFile) }
	readFile := func(t *testing.T) func(out any) error {
		t.Helper()
		src := filepath.Join("..", "..", "..", "testdata", t.Name()+craft.File)
		return func(out any) error { return files.ReadYAML(src, out, os.ReadFile) }
	}

	t.Run("invalid_bot", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/bot': value must be 'renovate'`, err.Error())
	})

	t.Run("dependabot_no_auth", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/maintenance': must not be provided`, err.Error())
	})

	t.Run("renovate_auth_required", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth': missing property 'maintenance'`, err.Error())
	})

	t.Run("auth_release_no_auth", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/release': must not be provided`, err.Error())
	})

	t.Run("release_auth_required", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth': missing property 'release'`, err.Error())
	})

	t.Run("release_gitlab_no_auth", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/release': must not be provided`, err.Error())
	})

	t.Run("empty", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("gitlab_no_bot", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("gitlab_renovate", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("github_release", func(t *testing.T) {
		// Act
		err := files.Validate(readSchema, readFile(t))

		// Assert
		assert.NoError(t, err)
	})
}
