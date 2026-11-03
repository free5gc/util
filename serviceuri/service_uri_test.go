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
		{
			name: "HTTP_IP_with_path",
			service: Service{
				Scheme:    "http",
				APIPrefix: "/operator/edge",
				Endpoints: []Endpoint{{IPv4Address: "10.101.200.28", Port: 8000}},
			},
			want: "http://10.101.200.28:8000/operator/edge",
		},
		{
			name: "HTTPS_FQDN_with_path",
			service: Service{
				Scheme:    "https",
				FQDN:      "amf.example.org",
				APIPrefix: "/sbi",
				Endpoints: []Endpoint{{Port: 8443}},
			},
			want: "https://amf.example.org:8443/sbi",
		},
		{
			name: "HTTP_FQDN_with_path",
			service: Service{
				Scheme:    "http",
				FQDN:      "amf.example.org",
				APIPrefix: "/operator/edge",
				Endpoints: []Endpoint{{Port: 8000}},
			},
			want: "http://amf.example.org:8000/operator/edge",
		},
		{
			name: "HTTP_IP_no_prefix",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "10.101.200.28", Port: 8000}},
			},
			want: "http://10.101.200.28:8000",
		},
		{
			name: "Legacy_absolute_allowed",
			service: Service{
				Scheme:    "http",
				APIPrefix: "http://10.101.200.18:8000",
			},
			opts:   Options{AllowLegacyAbsoluteAPIPrefix: true},
			want:   "http://10.101.200.18:8000",
			legacy: true,
		},
		{
			name: "Legacy_absolute_rejected",
			service: Service{
				Scheme:    "http",
				APIPrefix: "http://10.101.200.18:8000",
			},
			errorPart: "non-standard",
		},
		{
			name: "Bad_scheme",
			service: Service{
				Scheme:    "ftp",
				Endpoints: []Endpoint{{IPv4Address: "127.0.0.1"}},
			},
			errorPart: "scheme",
		},
		{
			name: "Invalid_path",
			service: Service{
				Scheme:    "http",
				APIPrefix: "relative",
				Endpoints: []Endpoint{{IPv4Address: "127.0.0.1"}},
			},
			errorPart: "apiPrefix",
		},
		{
			name: "Bad_IP",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "999.1.2.3", Port: 8000}},
			},
			errorPart: "IPv4",
		},
		{
			name: "HTTPS_requires_FQDN",
			service: Service{
				Scheme:    "https",
				Endpoints: []Endpoint{{IPv4Address: "127.0.0.1", Port: 8000}},
			},
			errorPart: "FQDN",
		},
		{
			name:    "Profile_FQDN_fallback",
			profile: Profile{FQDN: "amf.example.org"},
			service: Service{Scheme: "https"},
			want:    "https://amf.example.org:443",
		},
		{
			name: "IPv6",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv6Address: "2001:db8::1", Port: 8080}},
			},
			want: "http://[2001:db8::1]:8080",
		},
		{
			// TS 29.510 NFService.ipEndPoints: with HTTP, service IP addresses take
			// precedence over the NFProfile FQDN (e.g. OAI UDR registration).
			name:    "Service_IP_precedence_over_profile_FQDN",
			profile: Profile{FQDN: "oai-udr"},
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "192.168.70.136", Port: 8080}},
			},
			want: "http://192.168.70.136:8080",
		},
		{
			name:    "Service_FQDN_precedence_over_profile_FQDN",
			profile: Profile{FQDN: "profile.example.org"},
			service: Service{Scheme: "https", FQDN: "svc.example.org"},
			want:    "https://svc.example.org:443",
		},
		{
			name:    "Profile_IPv4_fallback",
			profile: Profile{IPv4Addresses: []string{"10.0.0.1"}},
			service: Service{Scheme: "http", Endpoints: []Endpoint{{Port: 8000}}},
			want:    "http://10.0.0.1:8000",
		},
		{
			name:    "Profile_IPv6_fallback",
			profile: Profile{IPv6Addresses: []string{"2001:db8::2"}},
			service: Service{Scheme: "http"},
			want:    "http://[2001:db8::2]:80",
		},
		{
			name:      "Profile_invalid_IPv4",
			profile:   Profile{IPv4Addresses: []string{"999.1.2.3"}},
			service:   Service{Scheme: "http"},
			errorPart: "NFProfile IPv4",
		},
		{
			name:      "Profile_hostname_as_IPv4",
			profile:   Profile{IPv4Addresses: []string{"amf.free5gc.org"}},
			service:   Service{Scheme: "http"},
			errorPart: "NFProfile IPv4",
		},
		{
			name:      "Profile_IPv4_in_IPv6_list",
			profile:   Profile{IPv6Addresses: []string{"10.0.0.1"}},
			service:   Service{Scheme: "http"},
			errorPart: "NFProfile IPv6",
		},
		{
			name: "Ambiguous_FQDN_ports",
			service: Service{
				Scheme:    "https",
				FQDN:      "amf.example.org",
				Endpoints: []Endpoint{{Port: 8001}, {Port: 8002}},
			},
			errorPart: "ambiguous",
		},
		{
			name: "FQDN_port_out_of_range",
			service: Service{
				Scheme:    "https",
				FQDN:      "amf.example.org",
				Endpoints: []Endpoint{{Port: 70000}},
			},
			errorPart: "port",
		},
		{
			name:      "Profile_FQDN_port_out_of_range",
			profile:   Profile{FQDN: "amf.example.org"},
			service:   Service{Scheme: "http", Endpoints: []Endpoint{{Port: -1}}},
			errorPart: "port",
		},
		{
			name: "Endpoint_with_both_families",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", IPv6Address: "2001:db8::1"}},
			},
			errorPart: "both",
		},
		{
			name: "Endpoint_port_out_of_range",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", Port: 70000}},
			},
			errorPart: "port",
		},
		{
			name: "InterPLMN_requires_FQDN",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", Port: 8000}},
			},
			opts:      Options{InterPLMN: true},
			errorPart: "FQDN",
		},
		{
			// TS 29.510 6.1.6.2.3 NOTE 1: inter-PLMN signalling uses the
			// interPlmnFqdn of the NF Service Profile, if any.
			name:    "InterPLMN_service_interPlmnFqdn",
			profile: Profile{FQDN: "amf.local", InterPlmnFQDN: "amf.profile.mnc001.mcc001.3gppnetwork.org"},
			service: Service{
				Scheme:        "https",
				FQDN:          "svc.local",
				InterPlmnFQDN: "svc.mnc001.mcc001.3gppnetwork.org",
				Endpoints:     []Endpoint{{IPv4Address: "10.0.0.1", Port: 8443}},
			},
			opts: Options{InterPLMN: true},
			want: "https://svc.mnc001.mcc001.3gppnetwork.org:8443",
		},
		{
			name:    "InterPLMN_profile_interPlmnFqdn",
			profile: Profile{FQDN: "amf.local", InterPlmnFQDN: "amf.mnc001.mcc001.3gppnetwork.org"},
			service: Service{Scheme: "https", FQDN: "svc.local"},
			opts:    Options{InterPLMN: true},
			want:    "https://amf.mnc001.mcc001.3gppnetwork.org:443",
		},
		{
			name:    "InterPLMN_falls_back_to_fqdn",
			profile: Profile{FQDN: "amf.mnc001.mcc001.3gppnetwork.org"},
			service: Service{Scheme: "https"},
			opts:    Options{InterPLMN: true},
			want:    "https://amf.mnc001.mcc001.3gppnetwork.org:443",
		},
		{
			name:    "IntraPLMN_ignores_interPlmnFqdn",
			profile: Profile{FQDN: "amf.local", InterPlmnFQDN: "amf.mnc001.mcc001.3gppnetwork.org"},
			service: Service{Scheme: "https"},
			want:    "https://amf.local:443",
		},
		{
			name: "IndirectSCP_requires_FQDN",
			service: Service{
				Scheme:    "http",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", Port: 8000}},
			},
			opts:      Options{IndirectSCP: true},
			errorPart: "FQDN",
		},
		{
			name: "Legacy_scheme_mismatch",
			service: Service{
				Scheme:    "https",
				APIPrefix: "http://10.0.0.1:8000",
			},
			opts:      Options{AllowLegacyAbsoluteAPIPrefix: true},
			errorPart: "invalid legacy",
		},
		{
			name: "Trailing_slash_trimmed",
			service: Service{
				Scheme:    "http",
				APIPrefix: "/sbi/",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", Port: 8000}},
			},
			want: "http://10.0.0.1:8000/sbi",
		},
		{
			name: "APIPrefix_with_space",
			service: Service{
				Scheme:    "http",
				APIPrefix: "/a b",
				Endpoints: []Endpoint{{IPv4Address: "10.0.0.1", Port: 8000}},
			},
			errorPart: "apiPrefix",
		},
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

// TS 29.510 NOTE 7: when several addresses are registered, the consumer selects
// one of them randomly. Every result must be one of the registered addresses.
func TestBuildSelectsRegisteredAddress(t *testing.T) {
	service := Service{
		Scheme: "http",
		Endpoints: []Endpoint{
			{IPv4Address: "10.0.0.1", Port: 8000},
			{IPv4Address: "10.0.0.2", Port: 8000},
		},
	}
	allowed := map[string]bool{
		"http://10.0.0.1:8000": true,
		"http://10.0.0.2:8000": true,
	}
	for i := 0; i < 20; i++ {
		got, err := Build(Profile{}, service, Options{})
		if err != nil {
			t.Fatalf("error=%v", err)
		}
		if !allowed[got.APIRoot] {
			t.Fatalf("unexpected APIRoot %q", got.APIRoot)
		}
	}
}

func TestErrAddressUnavailable(t *testing.T) {
	_, err := Build(Profile{}, Service{Scheme: "http"}, Options{})
	if !errors.Is(err, ErrAddressUnavailable) {
		t.Fatalf("error=%v", err)
	}
}
