package utils

import "strings"

// GetModuleAndOperation extracts module and operation from a URL path
func GetModuleAndOperation(path string) (module, operation string) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "unknown", "unknown"
	}

	// Remove empty strings
	var validParts []string
	for _, part := range parts {
		if part != "" {
			validParts = append(validParts, part)
		}
	}

	if len(validParts) == 0 {
		return "unknown", "unknown"
	}

	// First part as module
	module = validParts[0]

	// Remaining parts as operation
	if len(validParts) > 1 {
		operation = strings.Join(validParts[1:], "_")
	} else {
		operation = "unknown"
	}

	return module, operation
}
