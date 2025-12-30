package validator

import "regexp"

// IsValidPlmnId checks if the string is a valid PLMN ID (MCC+MNC)
// Format: 5 or 6 digits
func IsValidPlmnId(plmnId string) bool {
	match, err := regexp.MatchString(`^\d{5,6}$`, plmnId)
	if err != nil {
		return false
	}
	return match
}
