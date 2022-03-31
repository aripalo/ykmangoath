package ykmangoath

import (
	"context"
)

func List(ctx context.Context, options Options) (error, []string) {

	queryOptions := queryOptions{
		deviceID: options.DeviceID,
		password: options.Password,
		args:     []string{"list"},
	}

	err, output := performQuery(ctx, queryOptions)

	if err != nil {
		return err, nil
	}

	return nil, getLines(output)
}

func ListWithPasswordPrompt(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (error, string),
	options Options,
) (error, []string) {

	err, result, _ := ListWithPasswordPromptAndCache(ctx, passwordPrompt, options)
	return err, result
}

func ListWithPasswordPromptAndCache(
	ctx context.Context,
	passwordPrompt func(ctx context.Context) (error, string),
	options Options,
) (error, []string, string) {

	if options.Password != "" {
		return ErrPasswordNotAllowedWithPrompt, nil, ""
	}

	err, result := List(ctx, Options{DeviceID: options.DeviceID})

	if err != ErrOathAccountPasswordProtected {
		return err, result, ""
	}

	err, password := passwordPrompt(ctx)
	if err != nil {
		return err, nil, ""
	}

	err, result = List(ctx, Options{DeviceID: options.DeviceID, Password: password})

	return err, result, password
}
