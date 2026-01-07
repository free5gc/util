package validator

import "testing"

func TestIsValidSuci(t *testing.T) {
	tests := []struct {
		name string
		suci string
		want bool
	}{
		{"Valid SUCI", "suci-0-208-93-0-0-0-0000000003", true},
		{"Valid SUCI with hex", "suci-0-208-93-A-1-FF-AABBCC", true},
		{"Valid SUCI NAI", "nai-user@example.com", true},
		{"Invalid SUCI prefix", "imsi-208930000000003", false},
		{"Invalid SUCI format", "suci-0-208-93-0-0-0", false},
		{"Invalid SUCI non-hex output", "suci-0-208-93-0-0-0-ZZZ", false},
		{"Invalid SUPI Type", "suci-8-208-93-0-0-0-0000000003", false},
		{"Invalid NAI", "nai-userexample.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSuci(tt.suci); got != tt.want {
				t.Errorf("IsValidSuci(%v) = %v, want %v", tt.suci, got, tt.want)
			}
		})
	}
}
