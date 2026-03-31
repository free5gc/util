package validator

import "regexp"

var (
	plmnIdRegex = regexp.MustCompile(`^[0-9]{5,6}$`)
	mccRegex    = regexp.MustCompile(`^[0-9]{3}$`)
	mncRegex    = regexp.MustCompile(`^[0-9]{2,3}$`)
)

// IsValidPlmnId checks if the string is a valid PLMN ID (MCC+MNC)
// Format: 5 or 6 digits
func IsValidPlmnId(plmnId string) bool {
	return plmnIdRegex.MatchString(plmnId)
}

// IsValidPlmnIdParts checks whether MCC and MNC are individually valid and
// form a valid PLMN ID when concatenated.
func IsValidPlmnIdParts(mcc, mnc string) bool {
	if !mccRegex.MatchString(mcc) {
		return false
	}
	if !mncRegex.MatchString(mnc) {
		return false
	}

	return IsValidPlmnId(mcc + mnc)
}
