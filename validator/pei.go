package validator

import (
	"regexp"
	"strings"
)

// PEI types
const (
	PeiTypeImei    = "IMEI"
	PeiTypeImeisv  = "IMEISV"
	PeiTypeUnknown = "UNKNOWN"
)

var (
	peiImeiRegex   = regexp.MustCompile(`^imei-[0-9]{15}$`)
	peiImeisvRegex = regexp.MustCompile(`^imeisv-[0-9]{16}$`)
)

// IsValidPei checks if the string is a valid PEI (IMEI or IMEISV)
func IsValidPei(pei string) bool {
	if strings.HasPrefix(pei, "imei-") {
		return peiImeiRegex.MatchString(pei)
	}
	if strings.HasPrefix(pei, "imeisv-") {
		return peiImeisvRegex.MatchString(pei)
	}
	return false
}
