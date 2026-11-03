package validator

import "testing"

func TestIsValidPort(t *testing.T) {
	tests := []struct {
		name string
		port int32
		want bool
	}{
		{"Unspecified (0)", 0, true},
		{"Valid Port", 8000, true},
		{"Max Port", 65535, true},
		{"Negative", -1, false},
		{"Too Large", 65536, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPort(tt.port); got != tt.want {
				t.Errorf("IsValidPort(%d) = %v, want %v", tt.port, got, tt.want)
			}
		})
	}
}

func TestValidateIpEndPoint(t *testing.T) {
	tests := []struct {
		name    string
		ipv4    string
		ipv6    string
		port    int32
		wantErr bool
	}{
		{"Valid IPv4", "10.101.200.28", "", 8000, false},
		{"Valid IPv6", "", "2001:db8::1", 8080, false},
		{"Port Only", "", "", 8443, false},
		{"Both IPv4 and IPv6", "10.0.0.1", "2001:db8::1", 8000, true},
		{"Invalid IPv4", "999.1.2.3", "", 8000, true},
		{"Hostname as IPv4", "amf.free5gc.org", "", 8000, true},
		{"IPv6 in IPv4 field", "2001:db8::1", "", 8000, true},
		{"IPv4 in IPv6 field", "", "10.0.0.1", 8000, true},
		{"Invalid Port", "10.0.0.1", "", 70000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIpEndPoint(tt.ipv4, tt.ipv6, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIpEndPoint(%q, %q, %d) error = %v, wantErr %v",
					tt.ipv4, tt.ipv6, tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestValidateApiPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		wantErr bool
	}{
		{"Empty", "", false},
		{"Single Segment", "/sbi", false},
		{"Multiple Segments", "/operator/edge", false},
		{"Relative Path", "relative", true},
		{"Absolute URL", "http://10.0.0.1:8000", true},
		{"Network Path", "//host/path", true},
		{"Query", "/sbi?x=1", true},
		{"Fragment", "/sbi#x", true},
		{"Backslash", "/sbi\\x", true},
		{"Control Character", "/sbi\x00", true},
		{"Invalid Escape", "/sbi%zz", true},
		{"Space", "/a b", true},
		{"Unescaped Angle Bracket", "/a<b", true},
		{"Escaped Space", "/a%20b", false},
		{"Path Parameter", "/sbi;v=1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateApiPrefix(tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateApiPrefix(%q) error = %v, wantErr %v", tt.prefix, err, tt.wantErr)
			}
		})
	}
}
