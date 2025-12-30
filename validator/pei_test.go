package validator

import "testing"

func TestIsValidPei(t *testing.T) {
	tests := []struct {
		name string
		pei  string
		want bool
	}{
		{"Valid IMEI raw", "123456789012345", true},
		{"Valid IMEI prefix", "imei-123456789012345", true},
		{"Valid IMEISV raw", "1234567890123456", true},
		{"Valid IMEISV prefix", "imeisv-1234567890123456", true},
		{"Invalid Short", "123", false},
		{"Invalid Alpha", "12345678901234a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPei(tt.pei); got != tt.want {
				t.Errorf("IsValidPei(%v) = %v, want %v", tt.pei, got, tt.want)
			}
		})
	}
}
