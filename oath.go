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
	o := OathAccounts{ctx: ctx}
	result := deviceSerialPattern.FindString(serial)
	if result == "" {
		return o, fmt.Errorf("%w: %v", ErrDeviceSerial, serial)
	}
	o.serial = serial
	return o, nil
}

// GetSerial returns the currently configured Yubikey device serial.
func (o *OathAccounts) GetSerial() string {
	return o.serial
}

// SetPassword directly configures the Yubikey OATH application password.
// Mutually exclusive with SetPasswordPrompt.
func (o *OathAccounts) SetPassword(password string) error {
	if o.passwordPrompt != nil {
		return fmt.Errorf("%w: password prompt already set", ErrPasswordSetup)
	}
	o.password = password
	return nil
}

// SetPasswordPrompt configures a function that will be called upon if/when
// the Yubikey OATH application password is required.
// Mutually exclusive with SetPassword.
func (o *OathAccounts) SetPasswordPrompt(prompt func(ctx context.Context) (string, error)) error {
	if o.password != "" {
		return fmt.Errorf("%w: password already set", ErrPromptSetup)
	}
	o.passwordPrompt = prompt
	return nil
}

// GetPassword returns the password that successfully unlocked the Yubikey OATH application.
func (o *OathAccounts) GetPassword() (string, error) {
	if o.password == "" {
		return "", ErrNoPassword
	}
	return o.password, nil
}

// List returns a list of configured accounts in the Yubikey OATH application.
func (o *OathAccounts) List() ([]string, error) {
	queryOptions := ykmanOptions{
		serial:   o.serial,
		password: o.password,
		args:     []string{"list"},
	}

	o.ensurePrompt()
	output, err := executeYkmanWithPrompt(o.ctx, queryOptions, o.passwordPrompt)

	if err != nil {
		return nil, err
	}

	return getLines(output), nil
}

// Code requests a Time-based one-time password (TOPT) 6-digit code for given
// account (such as "<issuer>:<name>") from Yubikey OATH application.
func (o *OathAccounts) Code(account string) (string, error) {
	queryOptions := ykmanOptions{
		serial:   o.serial,
		password: o.password,
		args:     []string{"code", "--single", account},
	}

	o.ensurePrompt()
	output, err := executeYkmanWithPrompt(o.ctx, queryOptions, o.passwordPrompt)

	if err != nil {
		return output, err
	}

	return parseCode(output)
}

// ensurePrompt checks if prompt is not configured by user and then assigns
// a simple function as prompt that returns o.password.
func (o *OathAccounts) ensurePrompt() {
	if o.passwordPrompt == nil {
		o.passwordPrompt = func(ctx context.Context) (string, error) { return o.password, nil }
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
