package shareddefaults

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"os"
	"path/filepath"
	"runtime"
)

// SharedCredentialsFilename returns the SDK's default file path
// for the shared credentials file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.byteplus/credentials
//   - Windows: %USERPROFILE%\.byteplus\credentials
func SharedCredentialsFilename() string {
	return filepath.Join(UserHomeDir(), ".byteplus", "credentials")
}

// SharedConfigFilename returns the SDK's default file path for
// the shared config file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.byteplus/config
//   - Windows: %USERPROFILE%\.byteplus\config
func SharedConfigFilename() string {
	return filepath.Join(UserHomeDir(), ".byteplus", "config")
}

// UserHomeDir returns the home directory for the user the process is
// running under.
func UserHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
