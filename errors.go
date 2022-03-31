package ykmangoath

import "errors"

var (
	ErrDeviceSerial  = errors.New("invalid device serial")
	ErrPasswordSetup = errors.New("cannot set password")
	ErrPromptSetup   = errors.New("cannot set password prompt")
	ErrNoPassword    = errors.New("password not set")
)
