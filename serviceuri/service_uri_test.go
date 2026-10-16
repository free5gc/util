package serviceuri

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildTable(t *testing.T) {
	tests := []struct {
		name      string
		profile   Profile
		service   Service
		opts      Options
		want      string
		legacy    bool
		errorPart string
	}{
		{name: "HTTP_IP_with_path", service: Service{Scheme: "http", APIPrefix: "/operator/edge", Endpoints: []Endpoint{{IPv4Address: "10.101.200.28", Port: 8000}}}, want: "http://10.101.200.28:8000/operator/edge"},
		{name: "HTTPS_FQDN_with_path", service: Service{Scheme: "https", FQDN: "amf.example.org", APIPrefix: "/sbi", Endpoints: []Endpoint{{Port: 8443}}}, want: "https://amf.example.org:8443/sbi"},
		{name: "HTTP_FQDN_with_path", service: Service{Scheme: "http", FQDN: "amf.example.org", APIPrefix: "/operator/edge", Endpoints: []Endpoint{{Port: 8000}}}, want: "http://amf.example.org:8000/operator/edge"},
		{name: "HTTP_IP_no_prefix", service: Service{Scheme: "http", Endpoints: []Endpoint{{IPv4Address: "10.101.200.28", Port: 8000}}}, want: "http://10.101.200.28:8000"},
		{name: "Legacy_absolute_allowed", service: Service{Scheme: "http", APIPrefix: "http://10.101.200.18:8000"}, opts: Options{AllowLegacyAbsoluteAPIPrefix: true}, want: "http://10.101.200.18:8000", legacy: true},
		{name: "Legacy_absolute_rejected", service: Service{Scheme: "http", APIPrefix: "http://10.101.200.18:8000"}, errorPart: "non-standard"},
		{name: "Bad_scheme", service: Service{Scheme: "ftp", Endpoints: []Endpoint{{IPv4Address: "127.0.0.1"}}}, errorPart: "scheme"},
		{name: "Invalid_path", service: Service{Scheme: "http", APIPrefix: "relative", Endpoints: []Endpoint{{IPv4Address: "127.0.0.1"}}}, errorPart: "apiPrefix"},
		{name: "Bad_IP", service: Service{Scheme: "http", Endpoints: []Endpoint{{IPv4Address: "999.1.2.3", Port: 8000}}}, errorPart: "IPv4"},
		{name: "HTTPS_requires_FQDN", service: Service{Scheme: "https", Endpoints: []Endpoint{{IPv4Address: "127.0.0.1", Port: 8000}}}, errorPart: "FQDN"},
		{name: "Profile_FQDN_fallback", profile: Profile{FQDN: "amf.example.org"}, service: Service{Scheme: "https"}, want: "https://amf.example.org:443"},
		{name: "IPv6", service: Service{Scheme: "http", Endpoints: []Endpoint{{IPv6Address: "2001:db8::1", Port: 8080}}}, want: "http://[2001:db8::1]:8080"},
		{name: "No_address", service: Service{Scheme: "http"}, errorPart: "no usable"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Build(tc.profile, tc.service, tc.opts)
			if tc.errorPart != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errorPart) {
					t.Fatalf("error=%v; expected %q", err, tc.errorPart)
				}
				return
			}
			if err != nil {
				t.Fatalf("error=%v", err)
			}
			if got.APIRoot != tc.want || got.UsedLegacyAbsoluteAPIPrefix != tc.legacy {
				t.Fatalf("got=%+v wantRoot=%q legacy=%v", got, tc.want, tc.legacy)
			}
		})
	}
}

func TestErrAddressUnavailable(t *testing.T) {
	_, err := Build(Profile{}, Service{Scheme: "http"}, Options{})
	if !errors.Is(err, ErrAddressUnavailable) {
		t.Fatalf("error=%v", err)
	}
}
