package ykmangoath

import (
	"fmt"
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
			name:     "with password",
			input:    ykmanOptions{password: "p4ssword"},
			received: []string{"oath", "accounts", "--password", "p4ssword"},
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
			input:    ykmanOptions{serial: "12345678", password: "p4ssword", args: []string{"code", "--single", "Amazon Web Services:john.doe@example"}},
			received: []string{"--device", "12345678", "oath", "accounts", "code", "--single", "Amazon Web Services:john.doe@example", "--password", "p4ssword"},
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
