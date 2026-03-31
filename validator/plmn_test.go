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

func TestIsValidPlmnIdParts(t *testing.T) {
	tests := []struct {
		name string
		mcc  string
		mnc  string
		want bool
	}{
		{"Valid 3+2", "208", "93", true},
		{"Valid 3+3", "208", "930", true},
		{"Invalid short MCC", "20", "893", false},
		{"Invalid short MNC", "208", "9", false},
		{"Invalid alpha MCC", "2a8", "93", false},
		{"Invalid alpha MNC", "208", "9a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPlmnIdParts(tt.mcc, tt.mnc); got != tt.want {
				t.Errorf("IsValidPlmnIdParts(%v, %v) = %v, want %v", tt.mcc, tt.mnc, got, tt.want)
			}
		})
	}
}
