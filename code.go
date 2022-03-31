package ykmangoath

import (
	"context"
	"regexp"
	"strings"
)

// yubikeyTokenFindPattern describes the regexp that will match OATH TOPT MFA token code from Yubikey
var yubikeyTokenFindPattern = regexp.MustCompile(`\d{6}\d*$`)

// Code generates a OATH TOPT code from the Yubikey
func Code(ctx context.Context, account string, options Options) (error, string) {

	queryOptions := queryOptions{
		deviceID: options.DeviceID,
		password: options.Password,
		args:     []string{"code", "--single", account},
	}

	err, output := performQuery(ctx, queryOptions)

	if err != nil {
		return err, output
	}

	return parseCode(output)
}

// CodeWithPasswordPrompt generates a OATH TOPT code from the Yubikey with a password prompt support
func CodeWithPasswordPrompt(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (error, string),
	account string,
	options Options,
) (error, string) {
	err, result, _ := CodeWithPasswordPromptAndCache(ctx, passwordPrompt, account, options)
	return err, result
}

// CodeWithPasswordPromptAndCache generates a OATH TOPT code from the Yubikey with a password prompt support and also returns the password which succesfully unlocked the OATH application for caching purposes
func CodeWithPasswordPromptAndCache(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (error, string),
	account string,
	options Options,
) (error, string, string) {

	if options.Password != "" {
		return ErrPasswordNotAllowedWithPrompt, "", ""
	}

	err, result := Code(ctx, account, Options{DeviceID: options.DeviceID})

	if err != ErrOathAccountPasswordProtected {
		return err, result, ""
	}

	err, password := passwordPrompt(ctx)
	if err != nil {
		return err, "", ""
	}

	err, result = Code(ctx, account, Options{DeviceID: options.DeviceID, Password: password})
	return err, result, password
}

// parseCode retrieves the generated 6 digit OATH TOPT code from output
func parseCode(output string) (error, string) {
	result := yubikeyTokenFindPattern.FindString(strings.TrimSpace(output))
	if result == "" {
		return ErrOathAccountCodeParseFailed, ""
	}

	return nil, result
}
