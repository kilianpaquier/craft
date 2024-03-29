// Code generated by go-swagger; DO NOT EDIT.

package test_tag

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewTestTagParams creates a new TestTagParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewTestTagParams() *TestTagParams {
	return &TestTagParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewTestTagParamsWithTimeout creates a new TestTagParams object
// with the ability to set a timeout on a request.
func NewTestTagParamsWithTimeout(timeout time.Duration) *TestTagParams {
	return &TestTagParams{
		timeout: timeout,
	}
}

// NewTestTagParamsWithContext creates a new TestTagParams object
// with the ability to set a context for a request.
func NewTestTagParamsWithContext(ctx context.Context) *TestTagParams {
	return &TestTagParams{
		Context: ctx,
	}
}

// NewTestTagParamsWithHTTPClient creates a new TestTagParams object
// with the ability to set a custom HTTPClient for a request.
func NewTestTagParamsWithHTTPClient(client *http.Client) *TestTagParams {
	return &TestTagParams{
		HTTPClient: client,
	}
}

/*
TestTagParams contains all the parameters to send to the API endpoint

	for the test tag operation.

	Typically these are written to a http.Request.
*/
type TestTagParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the test tag params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *TestTagParams) WithDefaults() *TestTagParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the test tag params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *TestTagParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the test tag params
func (o *TestTagParams) WithTimeout(timeout time.Duration) *TestTagParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the test tag params
func (o *TestTagParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the test tag params
func (o *TestTagParams) WithContext(ctx context.Context) *TestTagParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the test tag params
func (o *TestTagParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the test tag params
func (o *TestTagParams) WithHTTPClient(client *http.Client) *TestTagParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the test tag params
func (o *TestTagParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *TestTagParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
