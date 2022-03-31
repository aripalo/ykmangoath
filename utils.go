package ykmangoath

import "strings"

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
