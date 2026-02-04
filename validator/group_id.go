package validator

import (
	"fmt"
	"regexp"
)

// TS 23.003 19.7.3 & TS 29.571 5.4.3.3
// External Group Identity format validation
// Expected pattern: prefix-mcc(3 digits)-mnc(2-3 digits)-localGroupId
var groupIdRegex = regexp.MustCompile(`^[^-]+-[0-9]{3}-[0-9]{2,3}-.+$`)

// ValidateGroupIdFormat validates the external group identity format
// Expected pattern: prefix-mcc(3 digits)-mnc(2-3 digits)-localGroupId
func ValidateGroupIdFormat(groupId string) error {
	if !groupIdRegex.MatchString(groupId) {
		return fmt.Errorf("invalid groupId format: expected pattern 'prefix-mcc-mnc-localGroupId', got '%s'", groupId)
	}
	return nil
}
