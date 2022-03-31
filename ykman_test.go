package ykmangoath

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineYkmanArgs(t *testing.T) {

	tests := []struct {
		name     string
		input    ykmanOptions
		received []string
	}{
		{
			name:     "empty options",
			input:    ykmanOptions{},
			received: []string{"oath", "accounts"},
		},
		{
			name:     "with serial",
			input:    ykmanOptions{serial: "12345678"},
			received: []string{"--device", "12345678", "oath", "accounts"},
		},
		{
			name:     "list accounts",
			input:    ykmanOptions{args: []string{"list"}},
			received: []string{"oath", "accounts", "list"},
		},
		{
			name:     "code for account",
			input:    ykmanOptions{args: []string{"code", "--single", "Amazon Web Services:john.doe@example"}},
			received: []string{"oath", "accounts", "code", "--single", "Amazon Web Services:john.doe@example"},
		},
		{
			name:     "code for account with all options",
			input:    ykmanOptions{serial: "12345678", args: []string{"code", "--single", "Amazon Web Services:john.doe@example"}},
			received: []string{"--device", "12345678", "oath", "accounts", "code", "--single", "Amazon Web Services:john.doe@example"},
		},
	}

	for index, test := range tests {

		name := fmt.Sprintf("case #%d - %s", index, test.name)

		t.Run(name, func(t *testing.T) {
			received := defineYkmanArgs(test.input)
			assert.Equal(t, test.received, received)
		})
	}
}

func TestProcessYkmanErrors(t *testing.T) {

	genericErr := errors.New("just for testing")

	tests := []struct {
		name         string
		err          error
		outputStderr string
		password     string
		received     error
	}{
		{
			name: "ykman not found",
			err: &exec.Error{
				Err: exec.ErrNotFound,
			},
			outputStderr: "",
			password:     "",
			received:     ErrCommandNotFound,
		},
		// TODO SIGINT
		// TODO SIGTREM
		{
			name:         "yubikey not found",
			err:          genericErr,
			outputStderr: "Failed connecting to the YubiKey",
			password:     "",
			received:     ErrDeviceNotFound,
		},
		{
			name:         "yubikey removed while in-use",
			err:          genericErr,
			outputStderr: "Failed to transmit with protocol",
			password:     "",
			received:     ErrDeviceRemoved,
		},
		{
			name:         "yubikey timeout",
			err:          genericErr,
			outputStderr: "Touch account timed out!",
			password:     "",
			received:     ErrDeviceTimeout,
		},
		{
			name:         "password not provided",
			err:          genericErr,
			outputStderr: "Authentication to the YubiKey failed. Wrong password?",
			password:     "",
			received:     ErrOathAccountPasswordProtected,
		},
		{
			name:         "invalid password",
			err:          genericErr,
			outputStderr: "Authentication to the YubiKey failed. Wrong password?",
			password:     "foo",
			received:     ErrOathAccountPasswordIncorrect,
		},
		{
			name:         "invalid account",
			err:          genericErr,
			outputStderr: "No matching account found.",
			password:     "",
			received:     ErrOathAccountNotFound,
		},
		{
			name:         "catch all",
			err:          genericErr,
			outputStderr: "",
			password:     "",
			received:     genericErr,
		},
	}

	for index, test := range tests {

		name := fmt.Sprintf("case #%d - %s", index, test.name)

		t.Run(name, func(t *testing.T) {
			received := processYkmanErrors(test.err, test.outputStderr, test.password)
			assert.Equal(t, test.received, received)
		})
	}
}
