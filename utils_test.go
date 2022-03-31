package ykmangoath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetLines is based on https://github.com/joshdk/ykmango/blob/master/utils_test.go
func TestGetLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		received []string
	}{
		{
			name:     "empty body",
			received: []string{},
		},
		{
			name:     "single space",
			input:    " ",
			received: []string{},
		},
		{
			name:     "multiple spaces",
			input:    "   ",
			received: []string{},
		},
		{
			name:     "single tab",
			input:    "\t",
			received: []string{},
		},
		{
			name:     "multiple tabs",
			input:    "\t\t\t",
			received: []string{},
		},
		{
			name:     "mixed whitespace",
			input:    "\t \t \t ",
			received: []string{},
		},
		{
			name:  "single word",
			input: "alice",
			received: []string{
				"alice",
			},
		},
		{
			name:  "single word surrounded with whitespace",
			input: "\t \t \t alice\t \t \t ",
			received: []string{
				"alice",
			},
		},
		{
			name:  "multiple words",
			input: "alice bob carol",
			received: []string{
				"alice bob carol",
			},
		},
		{
			name:  "multiple words surrounded with whitespace",
			input: "\t \t \t alice bob carol\t \t \t ",
			received: []string{
				"alice bob carol",
			},
		},
		{
			name:  "single line ending with a newline",
			input: "alice bob carol\n",
			received: []string{
				"alice bob carol",
			},
		},
		{
			name: "multiple body",
			input: `
				alice bob carol
				dave eve fred
			    grant henry ida
			`,
			received: []string{
				"alice bob carol",
				"dave eve fred",
				"grant henry ida",
			},
		},
		{
			name: "multiple body some blank",
			input: `
				alice bob carol
				dave eve fred
			    grant henry ida
			`,
			received: []string{
				"alice bob carol",
				"dave eve fred",
				"grant henry ida",
			},
		},
		{
			name: "multiple blank body",
			input: `
			`,
			received: []string{},
		},
	}

	for index, test := range tests {

		name := fmt.Sprintf("case #%d - %s", index, test.name)

		t.Run(name, func(t *testing.T) {
			received := getLines(test.input)
			assert.Equal(t, test.received, received)
		})
	}

}
