package serviceuri

import (
	"net"

	"github.com/free5gc/openapi/models"
)

// Addressing holds the addressing attributes an NF registers to the NRF.
type Addressing struct {
	FQDN          string
	IPv4Addresses []string
	IPv6Addresses []string
}

// NewAddressing maps the configured SBI register address (TS 29.510 6.1.6.2.2):
// an IPv4 or IPv6 literal goes to ipv4Addresses / ipv6Addresses, anything else
// (e.g. "amf.free5gc.org") is an FQDN.
func NewAddressing(registerAddr string) Addressing {
	ip := net.ParseIP(registerAddr)
	switch {
	case ip == nil:
		return Addressing{FQDN: registerAddr}
	case ip.To4() != nil:
		return Addressing{IPv4Addresses: []string{registerAddr}}
	default:
		return Addressing{IPv6Addresses: []string{registerAddr}}
	}
}

// IpEndPoints returns the NFService ipEndPoints for port (TS 29.510 6.1.6.2.5).
// With an FQDN only the port is registered.
func (a Addressing) IpEndPoints(port int32) []models.Nrf_NFMgmt_IpEndPoint {
	ep := models.Nrf_NFMgmt_IpEndPoint{Transport: models.Nrf_NFMgmt_TransportProtocol_TCP, Port: port}
	if len(a.IPv4Addresses) != 0 {
		ep.Ipv4Address = a.IPv4Addresses[0]
	}
	if len(a.IPv6Addresses) != 0 {
		ep.Ipv6Address = a.IPv6Addresses[0]
	}
	return []models.Nrf_NFMgmt_IpEndPoint{ep}
}
