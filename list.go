package ykmangoath

import (
	"context"
)

// ListWithPasswordPromptAndCache lists available OATH TOPT accounts configured in the Yubikey
func List(ctx context.Context, options Options) ([]string, error) {

	queryOptions := queryOptions{
		deviceID: options.DeviceID,
		password: options.Password,
		args:     []string{"list"},
	}

	output, err := performQuery(ctx, queryOptions)

	if err != nil {
		return nil, err
	}

	return getLines(output), nil
}

// ListWithPasswordPromptAndCache lists available OATH TOPT accounts configured in the Yubikey with a password prompt support
func ListWithPasswordPrompt(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (string, error),
	options Options,
) ([]string, error) {

	result, _, err := ListWithPasswordPromptAndCache(ctx, passwordPrompt, options)
	return result, err
}

// ListWithPasswordPromptAndCache lists available OATH TOPT accounts configured in the Yubikey with a password prompt support and also returns the password which succesfully unlocked the OATH application for caching purposes
func ListWithPasswordPromptAndCache(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (string, error),
	options Options,
) ([]string, string, error) {

	if options.Password != "" {
		return nil, "", ErrPasswordNotAllowedWithPrompt
	}

	result, err := List(ctx, Options{DeviceID: options.DeviceID})

	if err != ErrOathAccountPasswordProtected {
		return result, "", err
	}

	password, err := passwordPrompt(ctx)
	if err != nil {
		return nil, "", err
	}

	result, err = List(ctx, Options{DeviceID: options.DeviceID, Password: password})

	return result, password, err
}
