package ykmangoath

import (
	"context"
	"regexp"
	"strings"
)

// yubikeyTokenFindPattern describes the regexp that will match OATH TOPT MFA token code from Yubikey
var yubikeyTokenFindPattern = regexp.MustCompile(`\d{6}\d*$`)

// Code generates a OATH TOPT code from the Yubikey
func Code(ctx context.Context, account string, options Options) (string, error) {

	queryOptions := queryOptions{
		deviceID: options.DeviceID,
		password: options.Password,
		args:     []string{"code", "--single", account},
	}

	output, err := performQuery(ctx, queryOptions)

	if err != nil {
		return output, err
	}

	return parseCode(output)
}

// CodeWithPasswordPrompt generates a OATH TOPT code from the Yubikey with a password prompt support
func CodeWithPasswordPrompt(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (string, error),
	account string,
	options Options,
) (string, error) {
	result, _, err := CodeWithPasswordPromptAndCache(ctx, passwordPrompt, account, options)
	return result, err
}

// CodeWithPasswordPromptAndCache generates a OATH TOPT code from the Yubikey with a password prompt support and also returns the password which succesfully unlocked the OATH application for caching purposes
func CodeWithPasswordPromptAndCache(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (string, error),
	account string,
	options Options,
) (string, string, error) {

	if options.Password != "" {
		return "", "", ErrPasswordNotAllowedWithPrompt
	}

	result, err := Code(ctx, account, Options{DeviceID: options.DeviceID})

	if err != ErrOathAccountPasswordProtected {
		return result, "", err
	}

	password, err := passwordPrompt(ctx)
	if err != nil {
		return "", "", err
	}

	result, err = Code(ctx, account, Options{DeviceID: options.DeviceID, Password: password})
	return result, password, err
}

// parseCode retrieves the generated 6 digit OATH TOPT code from output
func parseCode(output string) (string, error) {
	result := yubikeyTokenFindPattern.FindString(strings.TrimSpace(output))
	if result == "" {
		return "", ErrOathAccountCodeParseFailed
	}

	return result, nil
}
