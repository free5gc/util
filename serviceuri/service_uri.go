// Package serviceuri resolves the {apiRoot} of an NF service (TS 29.501 4.4.1)
// from NRF discovery results (TS 29.510).
package serviceuri

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/free5gc/util/validator"
)

// ErrAddressUnavailable indicates that an NFService has no usable FQDN or IP address.
var ErrAddressUnavailable = errors.New("NFService has no usable FQDN or IP address")

type endpoint struct {
	IPv4Address string
	IPv6Address string
	Port        int32
}

type profileInfo struct {
	FQDN          string
	InterPlmnFQDN string
	IPv4Addresses []string
	IPv6Addresses []string
}

type serviceInfo struct {
	Scheme        string
	FQDN          string
	InterPlmnFQDN string
	APIPrefix     string
	Endpoints     []endpoint
}

type Options struct {
	// InterPLMN: the target NF is in a different PLMN.
	InterPLMN bool
	// IndirectSCP: Indirect Communication via an SCP.
	IndirectSCP bool
	// AllowLegacyAbsoluteAPIPrefix accepts free5GC's legacy apiPrefix that holds
	// a full URL. Not 3GPP compliant; for backward compatibility only.
	AllowLegacyAbsoluteAPIPrefix bool
}

// resolveService returns the API root of one NFService.
//
// TS 29.510 6.1.6.2.3 NOTE 1 / NOTE 8: an FQDN is used with https (TLS), for
// inter-PLMN signalling and for Indirect Communication; with http and Direct
// Communication, ipEndPoints take precedence over the NFProfile addresses.
func resolveService(profile profileInfo, service serviceInfo, opts Options) (string, error) {
	scheme := strings.ToLower(service.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("invalid NFService.scheme %q", service.Scheme)
	}
	prefix := service.APIPrefix
	if strings.HasPrefix(prefix, "http://") || strings.HasPrefix(prefix, "https://") {
		return buildLegacy(scheme, prefix, opts)
	}
	if err := validator.ValidateApiPrefix(prefix); err != nil {
		return "", err
	}

	if scheme == "https" || opts.InterPLMN || opts.IndirectSCP || service.FQDN != "" {
		host := selectFQDN(profile, service, opts.InterPLMN)
		if host == "" {
			return "", fmt.Errorf("FQDN is required for scheme %q and selected communication mode", scheme)
		}
		port, err := portForFQDN(service.Endpoints)
		if err != nil {
			return "", err
		}
		return buildAPIRoot(scheme, host, port, prefix)
	}

	var endpoints []endpoint
	for _, ep := range service.Endpoints {
		if ep.IPv4Address == "" && ep.IPv6Address == "" {
			continue
		}
		if err := validator.ValidateIpEndPoint(ep.IPv4Address, ep.IPv6Address, ep.Port); err != nil {
			return "", err
		}
		endpoints = append(endpoints, ep)
	}
	if len(endpoints) != 0 {
		// TS 29.510 6.1.6.2.3 NOTE 7: select one address randomly.
		ep := endpoints[rand.Intn(len(endpoints))]
		return buildAPIRoot(scheme, ep.IPv4Address+ep.IPv6Address, ep.Port, prefix)
	}

	port, err := portForFQDN(service.Endpoints)
	if err != nil {
		return "", err
	}
	if profile.FQDN != "" {
		return buildAPIRoot(scheme, profile.FQDN, port, prefix)
	}
	addresses, err := profileAddresses(profile)
	if err != nil {
		return "", err
	}
	if len(addresses) == 0 {
		return "", ErrAddressUnavailable
	}
	return buildAPIRoot(scheme, addresses[rand.Intn(len(addresses))], port, prefix)
}

func buildLegacy(scheme, prefix string, opts Options) (string, error) {
	if !opts.AllowLegacyAbsoluteAPIPrefix {
		return "", fmt.Errorf("non-standard absolute NFService.apiPrefix %q", prefix)
	}
	u, err := url.Parse(prefix)
	if err != nil || u.Scheme != scheme || u.Hostname() == "" || u.User != nil ||
		u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", fmt.Errorf("invalid legacy apiPrefix URL %q", prefix)
	}
	return strings.TrimSuffix(u.String(), "/"), nil
}

// selectFQDN applies TS 29.510 6.1.6.2.3 NOTE 1: the NFService value has
// precedence over the NFProfile value, and inter-PLMN signalling uses
// interPlmnFqdn. Without interPlmnFqdn it falls back to fqdn, which is where
// the NRF exposes interPlmnFqdn towards a different PLMN.
func selectFQDN(profile profileInfo, service serviceInfo, interPLMN bool) string {
	if interPLMN {
		if service.InterPlmnFQDN != "" {
			return service.InterPlmnFQDN
		}
		if profile.InterPlmnFQDN != "" {
			return profile.InterPlmnFQDN
		}
	}
	if service.FQDN != "" {
		return service.FQDN
	}
	return profile.FQDN
}

func profileAddresses(profile profileInfo) ([]string, error) {
	addresses := make([]string, 0, len(profile.IPv4Addresses)+len(profile.IPv6Addresses))
	for _, addr := range profile.IPv4Addresses {
		if err := validator.ValidateIpEndPoint(addr, "", 0); err != nil {
			return nil, fmt.Errorf("invalid NFProfile IPv4 address: %w", err)
		}
		addresses = append(addresses, addr)
	}
	for _, addr := range profile.IPv6Addresses {
		if err := validator.ValidateIpEndPoint("", addr, 0); err != nil {
			return nil, fmt.Errorf("invalid NFProfile IPv6 address: %w", err)
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

// portForFQDN returns the single port registered in ipEndPoints (0 if none).
func portForFQDN(points []endpoint) (int32, error) {
	var result int32
	for _, ep := range points {
		if !validator.IsValidPort(ep.Port) {
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

// buildAPIRoot formats {scheme}://{host}:{port}{apiPrefix}. Port 0 means the
// default port of the scheme. host and port must already be validated.
func buildAPIRoot(scheme, host string, port int32, prefix string) (string, error) {
	if host == "" || strings.ContainsAny(host, "/?#@%[]\\ \t\n\r") ||
		(strings.Contains(host, ":") && net.ParseIP(host) == nil) {
		return "", fmt.Errorf("invalid target host %q", host)
	}
	if port == 0 {
		port = 80
		if scheme == "https" {
			port = 443
		}
	}
	return scheme + "://" + net.JoinHostPort(host, strconv.Itoa(int(port))) + strings.TrimSuffix(prefix, "/"), nil
}
