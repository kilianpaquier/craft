// Code generated by go-swagger; DO NOT EDIT.

package test_tag

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kilianpaquier/craft/models"
)

// TestTagReader is a Reader for the TestTag structure.
type TestTagReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *TestTagReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewTestTagNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewTestTagDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewTestTagNoContent creates a TestTagNoContent with default headers values
func NewTestTagNoContent() *TestTagNoContent {
	return &TestTagNoContent{}
}

/*
TestTagNoContent describes a response with status code 204, with default header values.

success response test tag result
*/
type TestTagNoContent struct {
}

// IsSuccess returns true when this test tag no content response has a 2xx status code
func (o *TestTagNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this test tag no content response has a 3xx status code
func (o *TestTagNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this test tag no content response has a 4xx status code
func (o *TestTagNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this test tag no content response has a 5xx status code
func (o *TestTagNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this test tag no content response a status code equal to that given
func (o *TestTagNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the test tag no content response
func (o *TestTagNoContent) Code() int {
	return 204
}

func (o *TestTagNoContent) Error() string {
	return fmt.Sprintf("[GET /unsecured/test-tag][%d] testTagNoContent ", 204)
}

func (o *TestTagNoContent) String() string {
	return fmt.Sprintf("[GET /unsecured/test-tag][%d] testTagNoContent ", 204)
}

func (o *TestTagNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewTestTagDefault creates a TestTagDefault with default headers values
func NewTestTagDefault(code int) *TestTagDefault {
	return &TestTagDefault{
		_statusCode: code,
	}
}

/*
TestTagDefault describes a response with status code -1, with default header values.

default error response
*/
type TestTagDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this test tag default response has a 2xx status code
func (o *TestTagDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this test tag default response has a 3xx status code
func (o *TestTagDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this test tag default response has a 4xx status code
func (o *TestTagDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this test tag default response has a 5xx status code
func (o *TestTagDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this test tag default response a status code equal to that given
func (o *TestTagDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the test tag default response
func (o *TestTagDefault) Code() int {
	return o._statusCode
}

func (o *TestTagDefault) Error() string {
	return fmt.Sprintf("[GET /unsecured/test-tag][%d] testTag default  %+v", o._statusCode, o.Payload)
}

func (o *TestTagDefault) String() string {
	return fmt.Sprintf("[GET /unsecured/test-tag][%d] testTag default  %+v", o._statusCode, o.Payload)
}

func (o *TestTagDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *TestTagDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
