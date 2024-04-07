// This file is safe to edit. Once it exists it will not be overwritten.

package test_tag_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	test_tag_api "github.com/kilianpaquier/craft/internal/api/test_tag"
	"github.com/kilianpaquier/craft/restapi/operations/test_tag"
)

func TestTestTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		params := test_tag.TestTagParams{
			HTTPRequest: &http.Request{},
		}
		response := test_tag.NewTestTagDefault(http.StatusOK)

		// Act
		responder := test_tag_api.TestTag(params)

		// Assert
		assert.Equal(t, response, responder)
	})
}
