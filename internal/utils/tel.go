package utils

import "regexp"

var (
	thaiMobileRegex = regexp.MustCompile(`^(06|08|09)[0-9]{8}$`)
)

func IsThaiMobile(tel string) bool {
	return thaiMobileRegex.MatchString(tel)
}
