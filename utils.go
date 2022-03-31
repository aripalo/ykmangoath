package ykmangoath

import "strings"

// getLines splits multiline strings to arrays of strings
func getLines(output string) []string {

	lines := []string{}

	for _, line := range strings.Split(output, "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}
