package openapi

import (
	ae "github.com/piyushkumar96/app-error"
)

var (
	// ErrRouteNotFound is returned when the request path/method does not match any OpenAPI route.
	ErrRouteNotFound = ae.GetCustomErr(
		"ERR_OPENAPI_1001",
		"route not found",
		false)

	// ErrValidationUnexpectedType is returned when the validator returns an error that is not openapi3.MultiError.
	ErrValidationUnexpectedType = ae.GetCustomErr(
		"ERR_OPENAPI_1002",
		"validation error is not of type openapi3 multi error",
		false)

	// ErrValidationFailed is the base error for request validation failures; message is replaced with the actual validation details.
	ErrValidationFailed = ae.GetCustomErr(
		"ERR_OPENAPI_1003",
		"openapi request validation failed",
		false)

	// ErrLoadSpec is returned when the OpenAPI spec file cannot be loaded.
	ErrLoadSpec = ae.GetCustomErr(
		"ERR_OPENAPI_1004",
		"failed to load open api specification",
		false)

	// ErrCreateRouter is returned when the OpenAPI router cannot be created from the spec.
	ErrCreateRouter = ae.GetCustomErr(
		"ERR_OPENAPI_1005",
		"failed to create router for open api validation",
		false)

	// ErrResponseValidationFailed is returned when the response does not match the OpenAPI spec (logged; response already sent).
	ErrResponseValidationFailed = ae.GetCustomErr(
		"ERR_OPENAPI_1006",
		"openapi response validation failed",
		false)
)
