package ykmangoath

import "errors"

// ErrDeviceSerial indicates incorrect Yubikey Device Serial Number.
// The serial number must be at least 8 digits long string.
var ErrDeviceSerial = errors.New("invalid device serial")

// ErrPasswordSetup means password string could not be configured and
// is usually returned when trying to set both password and password prompt
// which are mutually exclusive.
var ErrPasswordSetup = errors.New("cannot set password")

// ErrPromptSetup means password prompt could not be configured and
// is usually returned when trying to set both password and password prompt
// which are mutually exclusive.
var ErrPromptSetup = errors.New("cannot set password prompt")

// ErrNoPassword is returned when trying to retrieve a the password
// but it has not been resolved by password prompt (yet).
var ErrNoPassword = errors.New("password not set")

// ErrCommandNotFound indicates ykman executable not found from $PATH.
// To resolve the error, user should install the ykman CLI tool and try again.
var ErrCommandNotFound = errors.New("ykman command not found")

// ErrCommandInterrupted indicates ykman process was killed by an external
// signal (such as SIGINT or SIGTERM).
var ErrCommandInterrupted = errors.New("ykman command interrupted")

// ErrDeviceNotFound indicates Yubikey device is not connected to host machine.
var ErrDeviceNotFound = errors.New("yubikey device not found")

// ErrDeviceRemoved means the Yubikey device was physically removed
// during the operation from the host machine.
var ErrDeviceRemoved = errors.New("yubikey device removed")

// ErrDeviceTimeout indicates Yubikey device timed out,
// usually because it was not touched during ykman's timeout period.
var ErrDeviceTimeout = errors.New("yubikey device timeout")

// ErrOathAccountPasswordProtected indicates that the OATH application requires a password.
var ErrOathAccountPasswordProtected = errors.New("oath application is password protected")

// ErrOathAccountPasswordIncorrect indicates a wrong password for OATH application.
var ErrOathAccountPasswordIncorrect = errors.New("oath application password is incorrect")

// ErrOathAccountNotFound indicates that the given OATH account name does not exists.
var ErrOathAccountNotFound = errors.New("oath account not found")

// ErrOathAccountCodeParseFailed indicates that valid OATH account code was not received.
var ErrOathAccountCodeParseFailed = errors.New("oath account code could not be parsed")

// ErrPasswordNotAllowedWithPrompt indicates that password should not be passed when using prompt.
var ErrPasswordNotAllowedWithPrompt = errors.New("password string not allowed when using prompt")
