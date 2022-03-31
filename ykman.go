package ykmangoath

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// ykmanOptions controls the ykman operation performed
type ykmanOptions struct {
	serial   string
	password string
	args     []string
}

func executeYkmanWithPrompt(ctx context.Context, options ykmanOptions, prompt func(ctx context.Context) (string, error)) (string, error) {
	result, err := executeYkman(ctx, options)
	if err != ErrOathAccountPasswordProtected {
		return result, err
	}

	password, err := prompt(ctx)
	if err != nil {
		return "", err
	}

	options.password = password

	return executeYkman(ctx, options)
}

// executeYkman executes ykman with given options and handles most common errors
func executeYkman(ctx context.Context, options ykmanOptions) (string, error) {

	args := defineYkmanArgs(options)

	// define the ykman command to be run
	cmd := exec.CommandContext(ctx, "ykman", args...)

	// redirect stdout & stderr into byte buffer
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// execute the ykman command
	err := cmd.Run()

	err = processYkmanErrors(err, errb.String(), options.password)

	// finally return the ykman output
	return outb.String(), err
}

func defineYkmanArgs(options ykmanOptions) []string {
	args := []string{}

	// only apply device argument if an id is given
	if options.serial != "" {
		args = append(args, "--device", options.serial)
	}

	// setup oath application arguments
	args = append(args, "oath", "accounts")
	args = append(args, options.args...)

	if options.password != "" {
		args = append(args, "--password", options.password)
	}

	return args
}

func processYkmanErrors(err error, outputStderr string, password string) error {
	if err != nil {

		// check for ykman process existance
		if execErr, ok := err.(*exec.Error); ok {
			if execErr.Err == exec.ErrNotFound {
				return ErrCommandNotFound
			}
		}

		// check for ykman process interuption
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus).Signal()
			if status == syscall.SIGINT {
				fmt.Println("SIGINT")
				return ErrCommandInterrupted
			}
			if status == syscall.SIGTERM {
				fmt.Println("SIGTERM")
				return ErrCommandInterrupted
			}
		}

		// check for yubikey device connection
		if strings.Contains(outputStderr, "Failed connecting to the YubiKey") {
			return ErrDeviceNotFound
		}

		// check for yubikey device removal
		if strings.Contains(outputStderr, "Failed to transmit with protocol") {
			return ErrDeviceRemoved
		}

		// check for yubikey device timeout
		if strings.Contains(outputStderr, "Touch account timed out!") {
			return ErrDeviceTimeout
		}

		// check for oath password protection
		if strings.Contains(outputStderr, "Authentication to the YubiKey failed. Wrong password?") {
			if password == "" {
				return ErrOathAccountPasswordProtected
			} else {
				return ErrOathAccountPasswordIncorrect
			}
		}

		// check for oath account mismatch
		if strings.Contains(outputStderr, "No matching account found.") {
			return ErrOathAccountNotFound
		}

		// catch-all error
		return err
	}

	return nil
}
