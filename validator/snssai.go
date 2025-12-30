package validator

import (
	"regexp"
)

var snssaiRegex = regexp.MustCompile(`^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(-[0-9a-fA-F]{6})?$`)

// IsValidSnssai checks if the string is a valid S-NSSAI
// Format: SST (e.g. "1") or SST-SD (e.g. "1-000001")
func IsValidSnssai(snssai string) bool {
	return snssaiRegex.MatchString(snssai)
}
