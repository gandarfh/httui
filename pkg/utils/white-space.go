package utils

import (
	"fmt"
	"strings"
)

func AddWhiteSpace(value string, size, maxsize int) string {
	if len(value) == 0 {
		value = "-"
	}
	value = Truncate(value, maxsize)

	s := strings.Repeat(" ", size)
	s = s[len(value):]
	s = fmt.Sprint(value, s)

	return s
}
