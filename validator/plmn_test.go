package validator

import "testing"

func TestIsValidPlmnId(t *testing.T) {
	tests := []struct {
		name   string
		plmnId string
		want   bool
	}{
		{"Valid 5 digits", "20893", true},
		{"Valid 6 digits", "208930", true},
		{"Invalid short", "1234", false},
		{"Invalid long", "1234567", false},
		{"Invalid alpha", "2089a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPlmnId(tt.plmnId); got != tt.want {
				t.Errorf("IsValidPlmnId(%v) = %v, want %v", tt.plmnId, got, tt.want)
			}
		})
	}
}
