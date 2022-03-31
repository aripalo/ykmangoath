package ykmangoath

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type queryOptions struct {
	deviceID string
	password string
	args     []string
}

func performQuery(ctx context.Context, options queryOptions) (error, string) {

	args := []string{}

	if options.deviceID != "" {
		args = append(args, "--device", options.deviceID)
	}

	args = append(args, "oath", "accounts")
	args = append(args, options.args...)

	cmd := exec.CommandContext(ctx, "ykman", args...)

	// In case a password is provided, provide it to ykman via stdin
	var b bytes.Buffer
	b.Write([]byte(fmt.Sprintf("%s\n", options.password)))
	cmd.Stdin = &b

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	outputErr := errb.String()
	output := outb.String()

	//combined, err := cmd.CombinedOutput()
	//output := string(combined)

	if err != nil {

		// check for ykman process existance
		if execErr, ok := err.(*exec.Error); ok {
			if execErr.Err == exec.ErrNotFound {
				return ErrCommandNotFound, ""
			}
		}

		// check for ykman process interuption
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus).Signal()
			if status == syscall.SIGINT {
				fmt.Println("SIGINT")
				return ErrCommandInterrupted, ""
			}
			if status == syscall.SIGTERM {
				fmt.Println("SIGTERM")
				return ErrCommandInterrupted, ""
			}
		}

		// check for yubikey device connection
		if strings.Contains(outputErr, "Failed connecting to the YubiKey") {
			return ErrDeviceNotFound, ""
		}

		// check for yubikey device removal
		if strings.Contains(outputErr, "Failed to transmit with protocol") {
			return ErrDeviceRemoved, ""
		}

		// check for yubikey device timeout
		if strings.Contains(outputErr, "Touch account timed out!") {
			return ErrDeviceTimeout, ""
		}

		// check for oath password protection
		if strings.Contains(outputErr, "Authentication to the YubiKey failed. Wrong password?") {
			if options.password == "" {
				return ErrOathAccountPasswordProtected, ""
			} else {
				return ErrOathAccountPasswordIncorrect, ""
			}
		}

		// check for oath account mismatch
		if strings.Contains(outputErr, "No matching account found.") {
			return ErrOathAccountNotFound, ""
		}

		// catch-all error
		return err, ""
	}

	return err, output
}
