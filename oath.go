package ykmangoath

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// deviceSerialPattern ensures the given Yubikey device serial is either empty string or at least 8 digits
var deviceSerialPattern = regexp.MustCompile(`^$|^\d{8,}$`)

// OathAccounts represents a the main functionality of Yubikey OATH accounts.
type OathAccounts struct {
	passwordPrompt func(ctx context.Context) (string, error)
	ctx            context.Context
	serial         string
	password       string
}

// New defines a new instance of OathAccounts.
func New(ctx context.Context, serial string) (OathAccounts, error) {
	oa := OathAccounts{ctx: ctx}
	result := deviceSerialPattern.FindString(serial)
	if result == "" {
		return oa, fmt.Errorf("%w: %v", ErrDeviceSerial, serial)
	}
	oa.serial = serial
	return oa, nil
}

// GetSerial returns the currently configured Yubikey device serial.
func (oa *OathAccounts) GetSerial() string {
	return oa.serial
}

// IsPasswordProtected checks whether the OATH application is password protected
func (oa *OathAccounts) IsPasswordProtected() bool {
	queryOptions := ykmanOptions{
		serial:   oa.serial,
		password: "",
		args:     []string{"list"},
	}

	_, err := executeYkman(oa.ctx, queryOptions)
	return err == ErrOathAccountPasswordProtected
}

// SetPassword directly configures the Yubikey OATH application password.
// Mutually exclusive with SetPasswordPrompt.
func (oa *OathAccounts) SetPassword(password string) error {
	if oa.passwordPrompt != nil {
		return fmt.Errorf("%w: password prompt already set", ErrPasswordSetup)
	}
	oa.password = password
	return nil
}

// SetPasswordPrompt configures a function that will be called upon if/when
// the Yubikey OATH application password is required.
// Mutually exclusive with SetPassword.
func (oa *OathAccounts) SetPasswordPrompt(prompt func(ctx context.Context) (string, error)) error {
	if oa.password != "" {
		return fmt.Errorf("%w: password already set", ErrPromptSetup)
	}
	oa.passwordPrompt = prompt
	return nil
}

// GetPassword returns the password that successfully unlocked the Yubikey OATH application.
func (oa *OathAccounts) GetPassword() (string, error) {
	if oa.password == "" {
		return "", ErrNoPassword
	}
	return oa.password, nil
}

// List returns a list of configured accounts in the Yubikey OATH application.
func (oa *OathAccounts) List() ([]string, error) {
	queryOptions := ykmanOptions{
		serial:   oa.serial,
		password: oa.password,
		args:     []string{"list"},
	}

	oa.ensurePrompt()
	output, err := executeYkmanWithPrompt(oa.ctx, queryOptions, oa.passwordPrompt)

	if err != nil {
		return nil, err
	}

	return getLines(output), nil
}

// Code requests a Time-based one-time password (TOPT) 6-digit code for given
// account (such as "<issuer>:<name>") from Yubikey OATH application.
func (oa *OathAccounts) Code(account string) (string, error) {
	queryOptions := ykmanOptions{
		serial:   oa.serial,
		password: oa.password,
		args:     []string{"code", "--single", account},
	}

	oa.ensurePrompt()
	output, err := executeYkmanWithPrompt(oa.ctx, queryOptions, oa.passwordPrompt)

	if err != nil {
		return output, err
	}

	return parseCode(output)
}

// ensurePrompt checks if prompt is not configured by user and then assigns
// a simple function as prompt that returns oa.password.
func (oa *OathAccounts) ensurePrompt() {
	if oa.passwordPrompt == nil {
		oa.passwordPrompt = func(ctx context.Context) (string, error) { return oa.password, nil }
	}
}

// parseCode retrieves the generated 6 digit OATH TOPT code from output
func parseCode(output string) (string, error) {
	result := yubikeyTokenFindPattern.FindString(strings.TrimSpace(output))
	if result == "" {
		return "", ErrOathAccountCodeParseFailed
	}

	return result, nil
}

// yubikeyTokenFindPattern describes the regexp that will match OATH TOPT MFA token code from Yubikey
var yubikeyTokenFindPattern = regexp.MustCompile(`\d{6}\d*$`)
