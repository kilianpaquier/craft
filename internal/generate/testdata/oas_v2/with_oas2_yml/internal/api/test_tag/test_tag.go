// This file is safe to edit. Once it exists it will not be overwritten.

package test_tag

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/kilianpaquier/craft/restapi/operations/test_tag"
)

// TestTag - GET /unsecured/test-tag.
//
// tests that a tag operation is generated like it would be expected.
func TestTag(params test_tag.TestTagParams) middleware.Responder {
	_ = params.HTTPRequest.Context() // request context
	return middleware.NotImplemented("not implemented")
}
