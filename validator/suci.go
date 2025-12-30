package validator

import (
	"regexp"
	"strings"
)

// TS 23.003 28.7.3
// SUCI format validation
// suci-<MCC>-<MNC>-<Routing Indicator>-<Protection Scheme>-<Home Network Public Key Id>-<Scheme Output>
var suciRegex = regexp.MustCompile(`^suci-[0-9]{3}-[0-9]{2,3}-[a-fA-F0-9]{1,4}-[a-fA-F0-9]{1,2}-[a-fA-F0-9]{1,2}-[a-fA-F0-9]+$`)

// TS 23.003 28.7.2
// NAI format validation (SUCI can also be in NAI format)
var suciNaiRegex = regexp.MustCompile(`^nai-.+@.+$`)

// IsValidSuci checks if the given SUCI is valid
func IsValidSuci(suci string) bool {
	if strings.HasPrefix(suci, "suci-") {
		return suciRegex.MatchString(suci)
	}
	if strings.HasPrefix(suci, "nai-") {
		return suciNaiRegex.MatchString(suci)
	}
	return false
}
