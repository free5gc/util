// Package serviceuri constructs a 3GPP NF service API root from separately
// advertised scheme, service/profile addresses, endpoints, and apiPrefix.
// It deliberately avoids a dependency on free5gc/openapi so any NF can use it.
package serviceuri

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ErrAddressUnavailable indicates that another matching NFService may be tried.
var ErrAddressUnavailable = errors.New("NFService has no usable FQDN or IP address")

type Endpoint struct {
	IPv4Address string
	IPv6Address string
	Port        int32
}

type Profile struct {
	FQDN          string
	IPv4Addresses []string
	IPv6Addresses []string
}

type Service struct {
	Scheme    string
	FQDN      string
	APIPrefix string
	Endpoints []Endpoint
}

type Options struct {
	InterPLMN   bool
	IndirectSCP bool
	// Legacy compatibility only; 3GPP apiPrefix must be a path, not a URL.
	AllowLegacyAbsoluteAPIPrefix bool
}

type Result struct {
	APIRoot                     string
	UsedLegacyAbsoluteAPIPrefix bool
}

func Build(profile Profile, service Service, opts Options) (Result, error) {
	scheme := strings.ToLower(service.Scheme)
	if scheme != "http" && scheme != "https" {
		return Result{}, fmt.Errorf("invalid NFService.scheme %q", service.Scheme)
	}
	prefix := service.APIPrefix
	if strings.HasPrefix(prefix, "http://") || strings.HasPrefix(prefix, "https://") {
		if !opts.AllowLegacyAbsoluteAPIPrefix {
			return Result{}, fmt.Errorf("non-standard absolute NFService.apiPrefix %q", prefix)
		}
		u, err := url.Parse(prefix)
		if err != nil || u.Scheme != scheme || u.Hostname() == "" || u.User != nil ||
			u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
			return Result{}, fmt.Errorf("invalid legacy apiPrefix URL %q", prefix)
		}
		return Result{APIRoot: strings.TrimSuffix(u.String(), "/"), UsedLegacyAbsoluteAPIPrefix: true}, nil
	}
	if err := ValidateAPIPrefix(prefix); err != nil {
		return Result{}, err
	}

	requireFQDN := scheme == "https" || opts.InterPLMN || opts.IndirectSCP
	if requireFQDN || service.FQDN != "" {
		host := service.FQDN
		if host == "" {
			host = profile.FQDN
		}
		if host == "" {
			return Result{}, fmt.Errorf("FQDN is required for scheme %q and selected communication mode", scheme)
		}
		port, err := PortForFQDN(service.Endpoints)
		if err != nil {
			return Result{}, err
		}
		root, err := BuildAPIRoot(scheme, host, port, prefix)
		return Result{APIRoot: root}, err
	}

	candidates := make([]Endpoint, 0, len(service.Endpoints))
	for _, ep := range service.Endpoints {
		if ep.IPv4Address != "" || ep.IPv6Address != "" {
			if err := ValidateIPPoint(ep); err != nil {
				return Result{}, err
			}
			candidates = append(candidates, ep)
		}
	}
	if len(candidates) != 0 {
		ep := candidates[rand.Intn(len(candidates))]
		host := ep.IPv4Address
		if host == "" {
			host = ep.IPv6Address
		}
		root, err := BuildAPIRoot(scheme, host, ep.Port, prefix)
		return Result{APIRoot: root}, err
	}
	port, err := PortForFQDN(service.Endpoints)
	if err != nil {
		return Result{}, err
	}
	if profile.FQDN != "" {
		root, err := BuildAPIRoot(scheme, profile.FQDN, port, prefix)
		return Result{APIRoot: root}, err
	}
	addresses := append([]string{}, profile.IPv4Addresses...)
	addresses = append(addresses, profile.IPv6Addresses...)
	if len(addresses) == 0 {
		return Result{}, ErrAddressUnavailable
	}
	root, err := BuildAPIRoot(scheme, addresses[rand.Intn(len(addresses))], port, prefix)
	return Result{APIRoot: root}, err
}

func ValidateIPPoint(p Endpoint) error {
	if p.IPv4Address != "" && p.IPv6Address != "" {
		return fmt.Errorf("IpEndPoint cannot contain both IPv4 and IPv6")
	}
	if p.IPv4Address != "" {
		ip := net.ParseIP(p.IPv4Address)
		if ip == nil || ip.To4() == nil {
			return fmt.Errorf("invalid IpEndPoint IPv4: %q", p.IPv4Address)
		}
	}
	if p.IPv6Address != "" {
		ip := net.ParseIP(p.IPv6Address)
		if ip == nil || ip.To4() != nil {
			return fmt.Errorf("invalid IpEndPoint IPv6: %q", p.IPv6Address)
		}
	}
	if p.Port < 0 || p.Port > 65535 {
		return fmt.Errorf("invalid IpEndPoint port: %d", p.Port)
	}
	return nil
}

func PortForFQDN(points []Endpoint) (int32, error) {
	var result int32
	for _, ep := range points {
		if ep.Port < 0 || ep.Port > 65535 {
			return 0, fmt.Errorf("invalid IpEndPoint port: %d", ep.Port)
		}
		if ep.Port == 0 {
			continue
		}
		if result != 0 && result != ep.Port {
			return 0, fmt.Errorf("ambiguous ports for FQDN: %d and %d", result, ep.Port)
		}
		result = ep.Port
	}
	return result, nil
}

func ValidateAPIPrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if !strings.HasPrefix(prefix, "/") || strings.HasPrefix(prefix, "//") || strings.ContainsAny(prefix, "?#\\") {
		return fmt.Errorf("invalid NFService.apiPrefix path %q", prefix)
	}
	for _, r := range prefix {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("control character in NFService.apiPrefix")
		}
	}
	u, err := url.Parse(prefix)
	if err != nil || u.IsAbs() || u.Host != "" || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("invalid NFService.apiPrefix %q", prefix)
	}
	return nil
}

func BuildAPIRoot(scheme, host string, port int32, prefix string) (string, error) {
	if port < 0 || port > 65535 {
		return "", fmt.Errorf("invalid port %d", port)
	}
	if port == 0 {
		port = 80
		if scheme == "https" {
			port = 443
		}
	}
	if host == "" || strings.ContainsAny(host, "/?#@\\ \t\n\r") {
		return "", fmt.Errorf("invalid target host %q", host)
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	}
	if strings.Contains(host, ":") && net.ParseIP(host) == nil {
		return "", fmt.Errorf("invalid FQDN/IP host %q", host)
	}
	if strings.Contains(host, "%") {
		return "", fmt.Errorf("zone identifiers are not valid here: %q", host)
	}
	return scheme + "://" + net.JoinHostPort(host, strconv.Itoa(int(port))) + strings.TrimSuffix(prefix, "/"), nil
}
