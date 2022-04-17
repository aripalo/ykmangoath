package ykmangoath

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type Ykman struct {
	ctx      context.Context
	serial   string
	password string
}

func NewYkman(ctx context.Context, serial string) *Ykman {
	return &Ykman{ctx: ctx, serial: serial}
}

func (y *Ykman) Execute(args []string) (string, error) {
	// only apply device argument if an id is given
	if y.serial != "" {
		args = append(args, "--device", y.serial)
	}

	// define the ykman command to be run
	cmd := exec.CommandContext(y.ctx, "ykman", args...)

	// in case a password is provided, provide it to ykman via stdin
	// it's better to pass it in via stdin as it will fail on empty string immediately
	// if the oath is password protected
	var b bytes.Buffer
	b.Write([]byte(fmt.Sprintf("%s\n", y.password)))
	cmd.Stdin = &b

	// redirect stdout & stderr into byte buffer
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// execute the ykman command
	err := cmd.Run()

	err = processYkmanErrors(err, errb.String(), y.password)

	// finally return the ykman output
	return outb.String(), err
}

func (y *Ykman) ExecuteWithPrompt(args []string, prompt func(ctx context.Context) (string, error)) (string, error) {
	result, err := y.Execute(args)
	if err != ErrOathAccountPasswordProtected {
		return result, err
	}

	password, err := prompt(y.ctx)
	if err != nil {
		return "", err
	}

	y.password = password

	return y.Execute(args)
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
		if strings.Contains(outputStderr, "Failed to open device for communication") {
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
