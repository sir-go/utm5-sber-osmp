package utm

import (
	"strings"
	"unicode"
)

func HasDigitsOnly(s string) bool {
	if len(strings.TrimSpace(s)) < 1 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
