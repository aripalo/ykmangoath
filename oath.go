package ykmangoath

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

// deviceSerialPattern ensures the given Yubikey device serial is either empty string or at least 8 digits
var deviceSerialPattern = regexp.MustCompile(`^$|^\d{8,}$`)

// OathAccounts represents a the main functionality of Yubikey OATH accounts.
type OathAccounts struct {
	passwordPrompt func(ctx context.Context) (string, error)
	ctx            context.Context
	serial         string
	password       string
	ykman          Ykman
}

// New defines a new instance of OathAccounts.
func New(ctx context.Context, serial string) (OathAccounts, error) {
	oa := OathAccounts{ctx: ctx, serial: serial}
	oa.ykman = *NewYkman(ctx, oa.serial)
	result := deviceSerialPattern.FindString(oa.serial)
	if result == "" {
		return oa, fmt.Errorf("%w: %v", ErrDeviceSerial, oa.serial)
	}
	oa.serial = serial
	return oa, nil
}

// GetSerial returns the currently configured Yubikey device serial.
func (oa *OathAccounts) GetSerial() string {
	return oa.serial
}

// IsAvailable checks whether the Yubikey device is connected & available
func (oa *OathAccounts) IsAvailable() bool {
	_, err := oa.ykman.Execute([]string{"info"})
	return err == nil
}

// IsPasswordProtected checks whether the OATH application is password protected
func (oa *OathAccounts) IsPasswordProtected() bool {
	_, err := oa.ykman.Execute([]string{"oath", "accounts", "list"})
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
	if oa.ykman.password == "" {
		return "", ErrNoPassword
	}
	return oa.ykman.password, nil
}

// List returns a list of configured accounts in the Yubikey OATH application.
func (oa *OathAccounts) List() ([]string, error) {
	oa.ensurePrompt()
	output, err := oa.ykman.ExecuteWithPrompt([]string{"oath", "accounts", "list"}, oa.passwordPrompt)

	if err != nil {
		return nil, err
	}

	return getLines(output), nil
}

// HasAccount returns a boolean indicating if the device has the given account
// configured in its OATH application.
func (oa *OathAccounts) HasAccount(account string) (bool, error) {
	accounts, err := oa.List()
	if err != nil {
		return false, err
	}

	return slices.Contains(accounts, account), err
}

// Code requests a Time-based one-time password (TOTP) 6-digit code for given
// account (such as "<issuer>:<name>") from Yubikey OATH application.
func (oa *OathAccounts) Code(account string) (string, error) {
	oa.ensurePrompt()
	output, err := oa.ykman.ExecuteWithPrompt([]string{"oath", "accounts", "code", "--single", account}, oa.passwordPrompt)

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

// parseCode retrieves the generated 6 digit OATH TOTP code from output
func parseCode(output string) (string, error) {
	result := yubikeyTokenFindPattern.FindString(strings.TrimSpace(output))
	if result == "" {
		return "", ErrOathAccountCodeParseFailed
	}

	return result, nil
}

// yubikeyTokenFindPattern describes the regexp that will match OATH TOTP MFA token code from Yubikey
var yubikeyTokenFindPattern = regexp.MustCompile(`\d{6}\d*$`)
