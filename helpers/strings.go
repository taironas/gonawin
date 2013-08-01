package helpers

import (
	"strings"
)

// TrimLower returns a lower case slice of the string s, with all leading and trailing white space removed, as defined by Unicode.
func TrimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}