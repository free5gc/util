package validator

import (
	"regexp"
	"strconv"
	"strings"
)

// IsValidSnssai checks if the string is a valid S-NSSAI
// Format: SST (e.g. "1") or SST-SD (e.g. "1-000001")
func IsValidSnssai(snssai string) bool {
	parts := strings.Split(snssai, "-")
	if len(parts) == 1 {
		return isValidSst(parts[0])
	}
	if len(parts) == 2 {
		return isValidSst(parts[0]) && isValidSd(parts[1])
	}
	return false
}

func isValidSst(sst string) bool {
	// SST is 0-255
	val, err := strconv.Atoi(sst)
	if err != nil {
		return false
	}
	return val >= 0 && val <= 255
}

func isValidSd(sd string) bool {
	// SD is 6 hex digits
	match, err := regexp.MatchString(`^[0-9a-fA-F]{6}$`, sd)
	if err != nil {
		return false
	}
	return match
}
