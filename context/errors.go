package context

import (
	ae "github.com/piyushkumar96/app-error"
)

var (
	// Common error codes
	OnRequestFailure = ae.GetCustomErr(
		"ERR_CTX_1000",
		"request failed",
		false)
	OnResponseWriteError = ae.GetCustomErr(
		"ERR_CTX_1001",
		"error while writing response",
		false)
)
