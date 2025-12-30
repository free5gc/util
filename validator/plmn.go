package validator

import "regexp"

var plmnIdRegex = regexp.MustCompile(`^[0-9]{5,6}$`)

// IsValidPlmnId checks if the string is a valid PLMN ID (MCC+MNC)
// Format: 5 or 6 digits
func IsValidPlmnId(plmnId string) bool {
	return plmnIdRegex.MatchString(plmnId)
}
