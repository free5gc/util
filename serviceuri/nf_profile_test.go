package serviceuri

import (
	"errors"
	"strings"
	"testing"

	"github.com/free5gc/openapi/models"
)

const testServiceName = models.Nrf_NFMgmt_ServiceName_NAMF_COMM

func registered(s models.Nrf_NFDisc_NFService) models.Nrf_NFDisc_NFService {
	s.ServiceName = testServiceName
	s.NfServiceStatus = models.Nrf_NFMgmt_NFServiceStatus_REGISTERED
	return s
}

func ipv4Service(ip string, port int32) models.Nrf_NFDisc_NFService {
	return registered(models.Nrf_NFDisc_NFService{
		Scheme:      models.UriScheme_HTTP,
		IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Ipv4Address: ip, Port: port}},
	})
}

func TestResolveNFServiceURISelection(t *testing.T) {
	const want = "http://10.101.200.28:8000"

	valid := ipv4Service("10.101.200.28", 8000)
	wrongName := valid
	wrongName.ServiceName = models.Nrf_NFMgmt_ServiceName_NUDR_DR
	notRegistered := valid
	notRegistered.NfServiceStatus = models.Nrf_NFMgmt_NFServiceStatus_SUSPENDED
	noAddress := registered(models.Nrf_NFDisc_NFService{Scheme: models.UriScheme_HTTP})

	tests := []struct {
		name     string
		services []models.Nrf_NFDisc_NFService
		want     string
		wantErr  error
	}{
		{"No_services", nil, "", ErrServiceNotFound},
		{"Wrong_service_name", []models.Nrf_NFDisc_NFService{wrongName}, "", ErrServiceNotFound},
		{"Not_registered", []models.Nrf_NFDisc_NFService{notRegistered}, "", ErrServiceNotFound},
		{"No_usable_address", []models.Nrf_NFDisc_NFService{noAddress}, "", ErrAddressUnavailable},
		{"Skip_wrong_name", []models.Nrf_NFDisc_NFService{wrongName, valid}, want, nil},
		{"Skip_not_registered", []models.Nrf_NFDisc_NFService{notRegistered, valid}, want, nil},
		{"Skip_no_usable_address", []models.Nrf_NFDisc_NFService{noAddress, valid}, want, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &models.Nrf_NFDisc_NFProfile{NfServices: tt.services}
			got, err := ResolveNFServiceURI(profile, testServiceName, Options{})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil || got.APIRoot != tt.want {
				t.Fatalf("got %q, error %v; want %q", got.APIRoot, err, tt.want)
			}
		})
	}
}

func TestResolveNFServiceURIStandardFields(t *testing.T) {
	tests := []struct {
		name    string
		profile models.Nrf_NFDisc_NFProfile
		want    string
	}{
		{
			name: "HTTP_service_FQDN_default_port",
			profile: models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
				registered(models.Nrf_NFDisc_NFService{Scheme: models.UriScheme_HTTP, Fqdn: "amf.example.org"}),
			}},
			want: "http://amf.example.org:80",
		},
		{
			name: "HTTPS_service_FQDN_port_and_prefix",
			profile: models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
				registered(models.Nrf_NFDisc_NFService{
					Scheme:      models.UriScheme_HTTPS,
					Fqdn:        "amf.example.org",
					ApiPrefix:   "/sbi",
					IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Port: 8443}},
				}),
			}},
			want: "https://amf.example.org:8443/sbi",
		},
		{
			// free5gc/free5gc#1173: OAI UDR registers a profile-level FQDN and a
			// service-level IP endpoint. With HTTP the service IP takes precedence.
			name: "OAI_profile_FQDN_with_service_IP",
			profile: models.Nrf_NFDisc_NFProfile{
				Fqdn:       "oai-udr",
				NfServices: []models.Nrf_NFDisc_NFService{ipv4Service("192.168.70.136", 8080)},
			},
			want: "http://192.168.70.136:8080",
		},
		{
			name: "Port_only_with_profile_FQDN",
			profile: models.Nrf_NFDisc_NFProfile{
				Fqdn: "profile.example.org",
				NfServices: []models.Nrf_NFDisc_NFService{
					registered(models.Nrf_NFDisc_NFService{
						Scheme:      models.UriScheme_HTTP,
						IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Port: 8080}},
					}),
				},
			},
			want: "http://profile.example.org:8080",
		},
		{
			name: "HTTPS_profile_FQDN_default_port",
			profile: models.Nrf_NFDisc_NFProfile{
				Fqdn: "profile.example.org",
				NfServices: []models.Nrf_NFDisc_NFService{
					registered(models.Nrf_NFDisc_NFService{Scheme: models.UriScheme_HTTPS}),
				},
			},
			want: "https://profile.example.org:443",
		},
		{
			name: "HTTP_IPv6_with_prefix",
			profile: models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
				registered(models.Nrf_NFDisc_NFService{
					Scheme:      models.UriScheme_HTTP,
					ApiPrefix:   "/api",
					IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Ipv6Address: "2001:db8::1", Port: 8000}},
				}),
			}},
			want: "http://[2001:db8::1]:8000/api",
		},
		{
			name: "NFServiceList_only",
			profile: models.Nrf_NFDisc_NFProfile{NfServiceList: map[string]models.Nrf_NFDisc_NFService{
				"comm": registered(models.Nrf_NFDisc_NFService{Scheme: models.UriScheme_HTTP, Fqdn: "amf.example.org"}),
			}},
			want: "http://amf.example.org:80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveNFServiceURI(&tt.profile, testServiceName, Options{})
			if err != nil || got.APIRoot != tt.want {
				t.Fatalf("got %q, error %v; want %q", got.APIRoot, err, tt.want)
			}
		})
	}
}

func TestResolveNFServiceURIInvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		scheme models.UriScheme
		prefix string
		fqdn   string
	}{
		{"Missing_scheme", "", "/x", "amf.example.org"},
		{"Bad_scheme", "ftp", "/x", "amf.example.org"},
		{"Relative_prefix", models.UriScheme_HTTP, "bad", "amf.example.org"},
		{"Query_prefix", models.UriScheme_HTTP, "/x?admin=1", "amf.example.org"},
		{"Authority_prefix", models.UriScheme_HTTP, "//evil.example.org", "amf.example.org"},
		{"Absolute_prefix", models.UriScheme_HTTP, "http://evil.example.org", "amf.example.org"},
		{"Missing_FQDN_for_HTTPS", models.UriScheme_HTTPS, "/x", ""},
		{"URL_as_FQDN", models.UriScheme_HTTP, "/x", "http://evil.example.org"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
				registered(models.Nrf_NFDisc_NFService{Scheme: tt.scheme, Fqdn: tt.fqdn, ApiPrefix: tt.prefix}),
			}}
			if got, err := ResolveNFServiceURI(profile, testServiceName, Options{}); err == nil {
				t.Fatalf("accepted invalid input: %q", got.APIRoot)
			}
		})
	}
}

func TestResolveNFServiceURIHTTPSNeverUsesBareIP(t *testing.T) {
	profile := &models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
		registered(models.Nrf_NFDisc_NFService{
			Scheme:      models.UriScheme_HTTPS,
			IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Ipv4Address: "192.0.2.10", Port: 8443}},
		}),
	}}
	got, err := ResolveNFServiceURI(profile, testServiceName, Options{})
	if err == nil || strings.Contains(got.APIRoot, "192.0.2.10") {
		t.Fatalf("HTTPS must not fall back to a bare IP: got %q, error %v", got.APIRoot, err)
	}
}

func TestResolveNFServiceURILegacyAbsolutePrefix(t *testing.T) {
	profile := &models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{
		registered(models.Nrf_NFDisc_NFService{
			Scheme:      models.UriScheme_HTTP,
			ApiPrefix:   "http://amf.free5gc.org:8000",
			IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Ipv4Address: "192.0.2.11", Port: 8000}},
		}),
	}}

	if _, err := ResolveNFServiceURI(profile, testServiceName, Options{}); err == nil {
		t.Fatal("strict mode must reject a non-standard absolute apiPrefix")
	}

	got, err := ResolveNFServiceURI(profile, testServiceName, Options{AllowLegacyAbsoluteAPIPrefix: true})
	if err != nil || got.APIRoot != "http://amf.free5gc.org:8000" || !got.UsedLegacyAbsoluteAPIPrefix {
		t.Fatalf("legacy mode returned %+v, error %v", got, err)
	}
}

func TestResolveNFServiceURINilProfile(t *testing.T) {
	if _, err := ResolveNFServiceURI(nil, testServiceName, Options{}); err == nil {
		t.Fatal("nil profile must return an error")
	}
}

// fixRandom makes candidate ordering deterministic: no shuffling and the
// random number is always r.
func fixRandom(t *testing.T, r int64) {
	t.Helper()
	origRand, origShuffle := randInt63n, shuffle
	randInt63n = func(n int64) int64 {
		if r >= n {
			return n - 1
		}
		return r
	}
	shuffle = func(int, func(i, j int)) {}
	t.Cleanup(func() { randInt63n, shuffle = origRand, origShuffle })
}

func instance(
	id string, priority, capacity int32, services ...models.Nrf_NFDisc_NFService,
) models.Nrf_NFDisc_NFProfile {
	return models.Nrf_NFDisc_NFProfile{
		NfInstanceId: id,
		Priority:     priority,
		Capacity:     capacity,
		NfServices:   services,
	}
}

func withPriority(s models.Nrf_NFDisc_NFService, priority, capacity int32) models.Nrf_NFDisc_NFService {
	s.Priority = priority
	s.Capacity = capacity
	return s
}

// TS 29.510 6.1.6.2.2 NOTE 4: lower priority values are selected first.
func TestResolveFromSearchResultPriority(t *testing.T) {
	fixRandom(t, 0)
	instances := []models.Nrf_NFDisc_NFProfile{
		instance("low", 20, 0, ipv4Service("10.0.0.20", 8000)),
		instance("high", 10, 0, ipv4Service("10.0.0.10", 8000)),
	}
	profile, got, err := ResolveFromSearchResult(instances, testServiceName, Options{})
	if err != nil || profile.NfInstanceId != "high" || got.APIRoot != "http://10.0.0.10:8000" {
		t.Fatalf("selected %v %q, error %v; want instance %q", profile, got.APIRoot, err, "high")
	}
}

// TS 29.510 6.1.6.2.2: NFService priority has precedence over NFProfile priority.
func TestResolveFromSearchResultServicePriorityPrecedence(t *testing.T) {
	fixRandom(t, 0)
	instances := []models.Nrf_NFDisc_NFProfile{
		instance("profile-high", 1, 0, withPriority(ipv4Service("10.0.0.1", 8000), 50, 0)),
		instance("profile-low", 100, 0, withPriority(ipv4Service("10.0.0.2", 8000), 5, 0)),
	}
	profile, _, err := ResolveFromSearchResult(instances, testServiceName, Options{})
	if err != nil || profile.NfInstanceId != "profile-low" {
		t.Fatalf("selected %v, error %v; want instance %q", profile, err, "profile-low")
	}
}

// A higher priority candidate that cannot be resolved falls back to the next one.
func TestResolveFromSearchResultPriorityFallback(t *testing.T) {
	fixRandom(t, 0)
	broken := registered(models.Nrf_NFDisc_NFService{Scheme: models.UriScheme_HTTPS})
	instances := []models.Nrf_NFDisc_NFProfile{
		instance("broken", 1, 0, broken),
		instance("backup", 2, 0, ipv4Service("10.0.0.2", 8000)),
	}
	profile, _, err := ResolveFromSearchResult(instances, testServiceName, Options{})
	if err != nil || profile.NfInstanceId != "backup" {
		t.Fatalf("selected %v, error %v; want instance %q", profile, err, "backup")
	}
}

// RFC 2782 weight selection within the same priority, used for capacity.
func TestOrderCandidatesCapacity(t *testing.T) {
	cands := func() []candidate {
		return []candidate{
			{service: models.Nrf_NFDisc_NFService{ServiceInstanceId: "w0"}, capacity: 0},
			{service: models.Nrf_NFDisc_NFService{ServiceInstanceId: "w10"}, capacity: 10},
			{service: models.Nrf_NFDisc_NFService{ServiceInstanceId: "w30"}, capacity: 30},
		}
	}
	// Running sums in RFC 2782 order (weight 0 first): w0=0, w10=10, w30=40.
	tests := []struct {
		r     int64
		first string
	}{
		{0, "w0"},
		{1, "w10"},
		{10, "w10"},
		{11, "w30"},
		{40, "w30"},
	}
	for _, tt := range tests {
		fixRandom(t, tt.r)
		got := orderCandidates(cands())
		if len(got) != 3 || got[0].service.ServiceInstanceId != tt.first {
			t.Fatalf("r=%d: first = %q, want %q", tt.r, got[0].service.ServiceInstanceId, tt.first)
		}
	}
}

func TestOrderCandidatesKeepsAllEntries(t *testing.T) {
	in := []candidate{
		{priority: 2, capacity: 5},
		{priority: 1, capacity: 0},
		{priority: 1, capacity: 7},
		{priority: 3, capacity: 0},
	}
	got := orderCandidates(in)
	if len(got) != len(in) {
		t.Fatalf("got %d candidates, want %d", len(got), len(in))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].priority > got[i].priority {
			t.Fatalf("candidates not ordered by priority: %+v", got)
		}
	}
}

// TS 29.510 6.1.6.2.2 NOTE 15: the NRF returns either nfServiceList or the
// deprecated nfServices; nfServiceList is used when present.
func TestResolveNFServiceURIPrefersNFServiceList(t *testing.T) {
	profile := &models.Nrf_NFDisc_NFProfile{
		NfServices: []models.Nrf_NFDisc_NFService{ipv4Service("10.0.0.1", 8000)},
		NfServiceList: map[string]models.Nrf_NFDisc_NFService{
			"svc-1": ipv4Service("10.0.0.2", 8000),
		},
	}
	got, err := ResolveNFServiceURI(profile, testServiceName, Options{})
	if err != nil || got.APIRoot != "http://10.0.0.2:8000" {
		t.Fatalf("got %q, error %v", got.APIRoot, err)
	}
}

// A higher priority service of the same NF that cannot be resolved falls back
// to the next service.
func TestResolveNFServiceURIPriorityFallback(t *testing.T) {
	fixRandom(t, 0)
	broken := withPriority(registered(models.Nrf_NFDisc_NFService{Scheme: "ftp", Fqdn: "amf.example.org"}), 1, 0)
	backup := withPriority(ipv4Service("10.0.0.2", 8000), 2, 0)
	profile := &models.Nrf_NFDisc_NFProfile{NfServices: []models.Nrf_NFDisc_NFService{backup, broken}}
	got, err := ResolveNFServiceURI(profile, testServiceName, Options{})
	if err != nil || got.APIRoot != "http://10.0.0.2:8000" {
		t.Fatalf("got %q, error %v", got.APIRoot, err)
	}
}

func TestResolveNFServiceURIInterPLMN(t *testing.T) {
	profile := &models.Nrf_NFDisc_NFProfile{
		Fqdn:          "udm.local",
		InterPlmnFqdn: "udm.5gc.mnc001.mcc001.3gppnetwork.org",
		NfServices: []models.Nrf_NFDisc_NFService{
			registered(models.Nrf_NFDisc_NFService{
				Scheme:      models.UriScheme_HTTPS,
				IpEndPoints: []models.Nrf_NFMgmt_IpEndPoint{{Ipv4Address: "10.0.0.1", Port: 443}},
			}),
		},
	}
	got, err := ResolveNFServiceURI(profile, testServiceName, Options{InterPLMN: true})
	if err != nil || got.APIRoot != "https://udm.5gc.mnc001.mcc001.3gppnetwork.org:443" {
		t.Fatalf("got %q, error %v", got.APIRoot, err)
	}
}

func TestResolveFromSearchResult(t *testing.T) {
	withoutService := models.Nrf_NFDisc_NFProfile{NfInstanceId: "no-service"}
	invalid := models.Nrf_NFDisc_NFProfile{
		NfInstanceId: "invalid",
		NfServices: []models.Nrf_NFDisc_NFService{
			registered(models.Nrf_NFDisc_NFService{Scheme: "ftp", Fqdn: "amf.example.org"}),
		},
	}
	valid := models.Nrf_NFDisc_NFProfile{
		NfInstanceId: "valid",
		NfServices:   []models.Nrf_NFDisc_NFService{ipv4Service("10.0.0.1", 8000)},
	}

	t.Run("Skips_unusable_instances", func(t *testing.T) {
		instances := []models.Nrf_NFDisc_NFProfile{withoutService, invalid, valid}
		profile, got, err := ResolveFromSearchResult(instances, testServiceName, Options{})
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if profile == nil || profile.NfInstanceId != "valid" {
			t.Fatalf("selected profile = %+v, want instance %q", profile, "valid")
		}
		if profile != &instances[2] {
			t.Fatal("returned profile must point into the given slice")
		}
		if got.APIRoot != "http://10.0.0.1:8000" {
			t.Fatalf("APIRoot = %q", got.APIRoot)
		}
	})

	t.Run("Empty_result", func(t *testing.T) {
		if _, _, err := ResolveFromSearchResult(nil, testServiceName, Options{}); !errors.Is(err, ErrServiceNotFound) {
			t.Fatalf("error = %v, want ErrServiceNotFound", err)
		}
	})

	t.Run("No_usable_instance", func(t *testing.T) {
		instances := []models.Nrf_NFDisc_NFProfile{withoutService, invalid}
		if _, _, err := ResolveFromSearchResult(instances, testServiceName, Options{}); err == nil {
			t.Fatal("expected an error when no instance is usable")
		}
	})
}
