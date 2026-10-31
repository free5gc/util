package serviceuri

import (
	"reflect"
	"testing"

	"github.com/free5gc/openapi/models"
)

func TestNewAddressing(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want Addressing
		ep   models.Nrf_NFMgmt_IpEndPoint
	}{
		{
			name: "IPv4",
			addr: "10.0.0.1",
			want: Addressing{IPv4Addresses: []string{"10.0.0.1"}},
			ep:   models.Nrf_NFMgmt_IpEndPoint{Ipv4Address: "10.0.0.1", Port: 8000},
		},
		{
			name: "IPv6",
			addr: "2001:db8::1",
			want: Addressing{IPv6Addresses: []string{"2001:db8::1"}},
			ep:   models.Nrf_NFMgmt_IpEndPoint{Ipv6Address: "2001:db8::1", Port: 8000},
		},
		{
			name: "FQDN",
			addr: "amf.free5gc.org",
			want: Addressing{FQDN: "amf.free5gc.org"},
			ep:   models.Nrf_NFMgmt_IpEndPoint{Port: 8000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAddressing(tt.addr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewAddressing(%q) = %+v, want %+v", tt.addr, got, tt.want)
			}
			tt.ep.Transport = models.Nrf_NFMgmt_TransportProtocol_TCP
			eps := got.IpEndPoints(8000)
			if len(eps) != 1 || eps[0] != tt.ep {
				t.Fatalf("IpEndPoints = %+v, want [%+v]", eps, tt.ep)
			}
		})
	}
}

// A registered profile must resolve back to a usable API root, including the
// docker-compose case where the register address is a hostname.
func TestAddressingRoundTrip(t *testing.T) {
	for addr, want := range map[string]string{
		"10.0.0.1":        "http://10.0.0.1:8000",
		"amf.free5gc.org": "http://amf.free5gc.org:8000",
	} {
		a := NewAddressing(addr)
		profile := &models.Nrf_NFDisc_NFProfile{
			Fqdn:          a.FQDN,
			Ipv4Addresses: a.IPv4Addresses,
			NfServices: []models.Nrf_NFDisc_NFService{
				registered(models.Nrf_NFDisc_NFService{
					Scheme:      models.UriScheme_HTTP,
					Fqdn:        a.FQDN,
					IpEndPoints: a.IpEndPoints(8000),
				}),
			},
		}
		got, err := ResolveNFServiceURI(profile, testServiceName, Options{})
		if err != nil || got != want {
			t.Fatalf("%s: got %q, error %v; want %q", addr, got, err, want)
		}
	}
}
