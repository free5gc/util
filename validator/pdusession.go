package validator

import "strconv"

// TS 29.571 5.4.2 valid PDU Session ID is an integer in the range 0 to 255
func IsValidPduSessionID(pduSessionID string) bool {
	pduSessionIDInt, err := strconv.Atoi(pduSessionID)
	return err == nil && pduSessionIDInt >= 0 && pduSessionIDInt < 256
}

// TS 24.501 9.11.3.57 (PSI 0-15)
func IsPduSessionIdInPsiRange(id int32) bool {
	return id >= 0 && id < 16
}
