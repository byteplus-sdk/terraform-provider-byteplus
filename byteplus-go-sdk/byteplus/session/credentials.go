package session

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"fmt"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/credentials"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/credentials/processcreds"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
)

func resolveCredentials(cfg *byteplus.Config,
	envCfg envConfig, sharedCfg sharedConfig,
	handlers request.Handlers,
	sessOpts Options,
) (*credentials.Credentials, error) {

	switch {
	case len(sessOpts.Profile) != 0:
		// User explicitly provided an Profile in the session's configuration
		// so load that profile from shared config first.
		// Github(byteplus/byteplus-go-sdk#2727)
		return resolveCredsFromProfile(cfg, envCfg, sharedCfg, handlers, sessOpts)

	case envCfg.Creds.HasKeys():
		// Environment credentials
		return credentials.NewStaticCredentialsFromCreds(envCfg.Creds), nil

	default:
		// Fallback to the "default" credential resolution chain.
		return resolveCredsFromProfile(cfg, envCfg, sharedCfg, handlers, sessOpts)
	}
}

func resolveCredsFromProfile(cfg *byteplus.Config,
	envCfg envConfig, sharedCfg sharedConfig,
	handlers request.Handlers,
	sessOpts Options,
) (creds *credentials.Credentials, err error) {

	switch {
	case sharedCfg.SourceProfile != nil:
		// Assume IAM role with credentials source from a different profile.
		creds, err = resolveCredsFromProfile(cfg, envCfg,
			*sharedCfg.SourceProfile, handlers, sessOpts,
		)

	case sharedCfg.Creds.HasKeys():
		// Static Credentials from Shared Config/Credentials file.
		creds = credentials.NewStaticCredentialsFromCreds(
			sharedCfg.Creds,
		)

	case len(sharedCfg.CredentialProcess) != 0:
		// Get credentials from CredentialProcess
		creds = processcreds.NewCredentials(sharedCfg.CredentialProcess)

	case len(sharedCfg.CredentialSource) != 0:
		creds, err = resolveCredsFromSource(cfg, envCfg,
			sharedCfg, handlers, sessOpts,
		)

	default:
		// Fallback to default credentials provider, include mock errors for
		// the credential chain so user can identify why credentials failed to
		// be retrieved.
		creds = credentials.NewCredentials(&credentials.ChainProvider{
			VerboseErrors: byteplus.BoolValue(cfg.CredentialsChainVerboseErrors),
			Providers: []credentials.Provider{
				&credProviderError{
					Err: bytepluserr.New("EnvAccessKeyNotFound",
						"failed to find credentials in the environment.", nil),
				},
				&credProviderError{
					Err: bytepluserr.New("SharedCredsLoad",
						fmt.Sprintf("failed to load profile, %s.", envCfg.Profile), nil),
				},
			},
		})
	}
	if err != nil {
		return nil, err
	}

	return creds, nil
}

// valid credential source values
const (
	credSourceEc2Metadata  = "Ec2InstanceMetadata"
	credSourceEnvironment  = "Environment"
	credSourceECSContainer = "EcsContainer"
)

func resolveCredsFromSource(cfg *byteplus.Config,
	envCfg envConfig, sharedCfg sharedConfig,
	handlers request.Handlers,
	sessOpts Options,
) (creds *credentials.Credentials, err error) {

	switch sharedCfg.CredentialSource {

	case credSourceEnvironment:
		creds = credentials.NewStaticCredentialsFromCreds(envCfg.Creds)

	default:
		return nil, ErrSharedConfigInvalidCredSource
	}

	return creds, nil
}

type credProviderError struct {
	Err error
}

func (c credProviderError) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, c.Err
}
func (c credProviderError) IsExpired() bool {
	return true
}
