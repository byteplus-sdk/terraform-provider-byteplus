package byteplus

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import "github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"

var (
	// ErrMissingRegion is an error that is returned if region configuration is
	// not found.
	ErrMissingRegion = bytepluserr.New("MissingRegion", "could not find region configuration", nil)

	// ErrMissingEndpoint is an error that is returned if an endpoint cannot be
	// resolved for a service.
	ErrMissingEndpoint = bytepluserr.New("MissingEndpoint", "'Endpoint' configuration is required for this service", nil)
)
