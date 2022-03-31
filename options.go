package ykmangoath

type Options struct {

	// DeviceID chooses which Yubikey device is used, can be omitted if only one
	DeviceID string

	// Password for Yubikey OATH application, can be omitted if passwordless oath
	Password string
}
