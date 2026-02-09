package authentication

import (
	ae "github.com/piyushkumar96/app-error"
)

var (
	// Common error codes
	ErrUnauthorized = ae.GetCustomErr(
		"ERR_AUTH_001",
		"user is not authorized to access this resource",
		false)
)
