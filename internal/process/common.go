package process

import "strings"

const (
	trimChars = ` "`
)

// clean removes unwanted characters from the input string to prepare it
// for parsing as URL.
func clean(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, trimChars), trimChars)
}
