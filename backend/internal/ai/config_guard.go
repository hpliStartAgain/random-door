package ai

import "strings"

func isPlaceholderValue(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return true
	}
	for _, marker := range []string{"replace_me", "your-", "your_", "change_me"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}
