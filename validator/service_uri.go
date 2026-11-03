package validator

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// IsValidPort checks the TCP port range. 0 means "not specified".
func IsValidPort(port int32) bool {
	return port >= 0 && port <= 65535
}

// ValidateIpEndPoint validates an IpEndPoint (TS 29.510 6.1.6.2.5).
// At most one of ipv4 and ipv6 may be set.
func ValidateIpEndPoint(ipv4, ipv6 string, port int32) error {
	if ipv4 != "" && ipv6 != "" {
		return fmt.Errorf("IpEndPoint cannot contain both IPv4 and IPv6")
	}
	if ipv4 != "" {
		if ip := net.ParseIP(ipv4); ip == nil || ip.To4() == nil {
			return fmt.Errorf("invalid IpEndPoint IPv4: %q", ipv4)
		}
	}
	if ipv6 != "" {
		if ip := net.ParseIP(ipv6); ip == nil || ip.To4() != nil {
			return fmt.Errorf("invalid IpEndPoint IPv6: %q", ipv6)
		}
	}
	if !IsValidPort(port) {
		return fmt.Errorf("invalid IpEndPoint port: %d", port)
	}
	return nil
}

// ValidateApiPrefix validates NFService.apiPrefix (TS 29.510 6.1.6.2.3):
// empty, or an escaped absolute path such as "/operator/edge" (not a URL).
func ValidateApiPrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	u, err := url.Parse(prefix)
	if err != nil || !strings.HasPrefix(prefix, "/") || u.Host != "" || u.EscapedPath() != prefix {
		return fmt.Errorf("invalid NFService.apiPrefix %q", prefix)
	}
	return nil
}
