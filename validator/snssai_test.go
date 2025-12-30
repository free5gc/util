package validator

import "testing"

func TestIsValidSnssai(t *testing.T) {
	tests := []struct {
		name   string
		snssai string
		want   bool
	}{
		{"Valid SST only", "1", true},
		{"Valid SST-SD", "1-000001", true},
		{"Valid SST-SD Upper", "1-AABBCC", true},
		{"Invalid SST range", "256", false},
		{"Invalid SST alpha", "a", false},
		{"Invalid SD length", "1-001", false},
		{"Invalid SD alpha", "1-GGGGGG", false},
		{"Invalid Format", "1-2-3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSnssai(tt.snssai); got != tt.want {
				t.Errorf("IsValidSnssai(%v) = %v, want %v", tt.snssai, got, tt.want)
			}
		})
	}
}
