package generate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestParserGit(t *testing.T) {
	ctx := context.Background()

	t.Run("success_no_vcs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		config := craft.Config{}

		// Act
		err := generate.ParserGit(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Zero(t, config)
	})

	t.Run("success_vcs", func(t *testing.T) {
		// Arrange
		expected := craft.Config{
			VCS: parser.VCS{
				Platform:    parser.GitHub,
				ProjectHost: "github.com",
				ProjectName: "craft",
				ProjectPath: "kilianpaquier/craft",
			},
		}
		config := craft.Config{}

		// Act
		err := generate.ParserGit(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
