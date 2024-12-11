package craft_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/craft"
)

func TestValidate(t *testing.T) {
	src := func(t *testing.T) string {
		t.Helper()
		return filepath.Join("..", "..", "testdata", t.Name()+craft.File)
	}

	t.Run("invalid_bot", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/bot': value must be 'renovate'`, err.Error())
	})

	t.Run("dependabot_no_auth", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/maintenance': must not be provided`, err.Error())
	})

	t.Run("renovate_auth_required", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth': missing property 'maintenance'`, err.Error())
	})

	t.Run("auth_release_no_auth", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/release': must not be provided`, err.Error())
	})

	t.Run("release_auth_required", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth': missing property 'release'`, err.Error())
	})

	t.Run("release_gitlab_no_auth", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.Equal(t, `validate schema:
- at '/ci/auth/release': must not be provided`, err.Error())
	})

	t.Run("empty", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("gitlab_no_bot", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("gitlab_renovate", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.NoError(t, err)
	})

	t.Run("github_release", func(t *testing.T) {
		// Act
		err := craft.Validate(src(t))

		// Assert
		assert.NoError(t, err)
	})
}
