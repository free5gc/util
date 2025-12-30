package validator

import "regexp"

// IsValidPlmnId checks if the string is a valid PLMN ID (MCC+MNC)
// Format: 5 or 6 digits
func IsValidPlmnId(plmnId string) bool {
	match, _ := regexp.MatchString(`^\d{5,6}$`, plmnId)
	return match
}
