// This file is safe to edit. Once it exists it will not be overwritten.

package monitoring_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	monitoring_api "github.com/kilianpaquier/craft/internal/api/monitoring"
	"github.com/kilianpaquier/craft/restapi/operations/monitoring"
)

func TestUnsecuredPing(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		params := monitoring.UnsecuredPingParams{
			HTTPRequest: &http.Request{},
		}
		response := monitoring.NewUnsecuredPingDefault(http.StatusOK)

		// Act
		responder := monitoring_api.UnsecuredPing(params)

		// Assert
		assert.Equal(t, response, responder)
	})
}
