package credentials

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"os"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
)

// EnvProviderName provides a name of Env provider
const EnvProviderName = "EnvProvider"

var (
	// ErrAccessKeyIDNotFound is returned when the byteplus Access Key ID can't be
	// found in the process's environment.
	ErrAccessKeyIDNotFound = bytepluserr.New("EnvAccessKeyNotFound", "BYTEPLUS_ACCESS_KEY_ID or BYTEPLUS_ACCESS_KEY not found in environment", nil)

	// ErrSecretAccessKeyNotFound is returned when the byteplus Secret Access Key
	// can't be found in the process's environment.
	ErrSecretAccessKeyNotFound = bytepluserr.New("EnvSecretNotFound", "BYTEPLUS_SECRET_ACCESS_KEY or BYTEPLUS_SECRET_KEY not found in environment", nil)
)

// A EnvProvider retrieves credentials from the environment variables of the
// running process. Environment credentials never expire.
//
// Environment variables used:
//
// * Access Key ID:     BYTEPLUS_ACCESS_KEY_ID or BYTEPLUS_ACCESS_KEY
//
// * Secret Access Key: BYTEPLUS_SECRET_ACCESS_KEY or BYTEPLUS_SECRET_KEY
type EnvProvider struct {
	retrieved bool
}

// NewEnvCredentials returns a pointer to a new Credentials object
// wrapping the environment variable provider.
func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (Value, error) {
	e.retrieved = false

	id := os.Getenv("BYTEPLUS_ACCESS_KEY_ID")
	if id == "" {
		id = os.Getenv("BYTEPLUS_ACCESS_KEY")
	}

	secret := os.Getenv("BYTEPLUS_SECRET_ACCESS_KEY")
	if secret == "" {
		secret = os.Getenv("BYTEPLUS_SECRET_KEY")
	}

	if id == "" {
		return Value{ProviderName: EnvProviderName}, ErrAccessKeyIDNotFound
	}

	if secret == "" {
		return Value{ProviderName: EnvProviderName}, ErrSecretAccessKeyNotFound
	}

	e.retrieved = true
	return Value{
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    os.Getenv("BYTEPLUS_SESSION_TOKEN"),
		ProviderName:    EnvProviderName,
	}, nil
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}
