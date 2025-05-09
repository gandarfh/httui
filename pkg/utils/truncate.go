package utils

func Truncate(text string, maxLen int) string {
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}

	return text
}
