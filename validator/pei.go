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

// IsValidPei checks if the string is a valid PEI (IMEI or IMEISV)
// It supports raw digits or "imei-"/"imeisv-" prefix
func IsValidPei(pei string) bool {
	if strings.HasPrefix(pei, "imei-") {
		return IsValidImei(pei[5:])
	}
	if strings.HasPrefix(pei, "imeisv-") {
		return IsValidImeisv(pei[7:])
	}
	// Try to guess based on length if no prefix
	if IsValidImei(pei) || IsValidImeisv(pei) {
		return true
	}
	return false
}

// IsValidImei checks if the string is a valid IMEI (15 digits)
func IsValidImei(imei string) bool {
	match, err := regexp.MatchString(`^\d{15}$`, imei)
	if err != nil {
		return false
	}
	return match
}

// IsValidImeisv checks if the string is a valid IMEISV (16 digits)
func IsValidImeisv(imeisv string) bool {
	match, err := regexp.MatchString(`^\d{16}$`, imeisv)
	if err != nil {
		return false
	}
	return match
}
