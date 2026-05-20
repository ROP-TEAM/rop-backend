package utils

import "strings"

func IsSpacedString(a *string) bool {
	return a == nil || strings.TrimSpace(*a) == ""
}
