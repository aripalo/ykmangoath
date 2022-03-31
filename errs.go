package ykmangoath

import "errors"

var (
	// ErrCommandNotFound indicates ykman executable not found from $PATH
	ErrCommandNotFound = errors.New("ykman command not found")

	// ErrCommandInterrupted indicates ykman process was killed by a signal (such as SIGINT or SIGTERM)
	ErrCommandInterrupted = errors.New("ykman command interrupted")

	// ErrDeviceNotFound indicates Yubikey device is not connected
	ErrDeviceNotFound = errors.New("yubikey device not found")

	// ErrDeviceRemoved indicates Yubikey device was removed during the operation
	ErrDeviceRemoved = errors.New("yubikey device removed")

	// ErrDeviceTimeout indicates Yubikey device timed out, usually because it was not touched
	ErrDeviceTimeout = errors.New("yubikey device timeout")

	// ErrOathAccountPasswordProtected indicates that the OATH application requires a password
	ErrOathAccountPasswordProtected = errors.New("oath application is password protected")

	// ErrOathAccountPasswordIncorrect indicates a wrong password for OATH application
	ErrOathAccountPasswordIncorrect = errors.New("oath application password is incorrect")

	// ErrOathAccountNotFound indicates that the given OATH account name does not exists
	ErrOathAccountNotFound = errors.New("oath account not found")

	// ErrOathAccountCodeParseFailed indicates that valid OATH account code was not received
	ErrOathAccountCodeParseFailed = errors.New("oath account code could not be parsed")

	// ErrPasswordNotAllowedWithPrompt indicates that password should not be passed when using prompt
	ErrPasswordNotAllowedWithPrompt = errors.New("password string not allowed when using prompt")
)
