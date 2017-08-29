package utils

import "strings"

// ParseKeyPair parses a key=value pair.
func ParseKeyPair(p []string) map[string]string {
	params := map[string]string{}
	for _, i := range p {
		parts := strings.Split(i, "=")
		if len(parts) != 2 {
			continue
		}
		params[parts[0]] = parts[1]
	}
	return params
}
