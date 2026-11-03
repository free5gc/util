package serviceuri

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/free5gc/openapi/models"
)

// ErrServiceNotFound indicates that no REGISTERED NFService with the requested
// service name exists.
var ErrServiceNotFound = errors.New("NFService not found")

// Random sources for NF selection; replaced in tests.
var (
	randInt63n = rand.Int63n
	shuffle    = rand.Shuffle
)

type candidate struct {
	profile  *models.Nrf_NFDisc_NFProfile
	service  models.Nrf_NFDisc_NFService
	priority int32
	capacity int32
}

// ResolveNFServiceURI returns the API root of service name in one NF profile.
func ResolveNFServiceURI(
	profile *models.Nrf_NFDisc_NFProfile,
	name models.Nrf_NFMgmt_ServiceName,
	opts Options,
) (Result, error) {
	if profile == nil {
		return Result{}, errors.New("NF profile is nil")
	}
	_, result, err := resolve(matchingCandidates(profile, name), name, opts)
	return result, err
}

// ResolveFromSearchResult selects an NF instance from an NRF discovery result
// and returns it with the API root of service name.
func ResolveFromSearchResult(
	instances []models.Nrf_NFDisc_NFProfile,
	name models.Nrf_NFMgmt_ServiceName,
	opts Options,
) (*models.Nrf_NFDisc_NFProfile, Result, error) {
	var cands []candidate
	for i := range instances {
		cands = append(cands, matchingCandidates(&instances[i], name)...)
	}
	return resolve(cands, name, opts)
}

// resolve tries the candidates in selection order and returns the first one
// that yields a valid API root.
func resolve(
	cands []candidate, name models.Nrf_NFMgmt_ServiceName, opts Options,
) (*models.Nrf_NFDisc_NFProfile, Result, error) {
	if len(cands) == 0 {
		return nil, Result{}, fmt.Errorf("%w: %q", ErrServiceNotFound, name)
	}
	var lastErr error
	for _, c := range orderCandidates(cands) {
		result, err := Build(toProfile(c.profile), toService(c.service), opts)
		if err == nil {
			return c.profile, result, nil
		}
		lastErr = fmt.Errorf("NFService %q of NF %s: %w", name, c.profile.NfInstanceId, err)
	}
	return nil, Result{}, lastErr
}

// matchingCandidates returns the REGISTERED NFServices named name.
//
// TS 29.510 6.1.6.2.2: the NRF returns either nfServiceList or the deprecated
// nfServices (NOTE 15), and NFService priority / capacity have precedence over
// the NFProfile values. The models cannot tell an absent value from 0, so 0
// falls back to the NFProfile value.
func matchingCandidates(profile *models.Nrf_NFDisc_NFProfile, name models.Nrf_NFMgmt_ServiceName) []candidate {
	services := profile.NfServices
	if len(profile.NfServiceList) != 0 {
		services = make([]models.Nrf_NFDisc_NFService, 0, len(profile.NfServiceList))
		for _, s := range profile.NfServiceList {
			services = append(services, s)
		}
	}

	var out []candidate
	for _, s := range services {
		if s.ServiceName != name || s.NfServiceStatus != models.Nrf_NFMgmt_NFServiceStatus_REGISTERED {
			continue
		}
		c := candidate{profile: profile, service: s, priority: s.Priority, capacity: s.Capacity}
		if c.priority == 0 {
			c.priority = profile.Priority
		}
		if c.capacity == 0 {
			c.capacity = profile.Capacity
		}
		out = append(out, c)
	}
	return out
}

// orderCandidates applies TS 29.510 6.1.6.2.2 NOTE 4: priority and capacity are
// used like SRV priority and weight in IETF RFC 2782. Lower priority comes
// first; within a priority, the RFC 2782 weighted random order is used.
func orderCandidates(cands []candidate) []candidate {
	sort.SliceStable(cands, func(i, j int) bool { return cands[i].priority < cands[j].priority })

	ordered := make([]candidate, 0, len(cands))
	for start := 0; start < len(cands); {
		end := start
		for end < len(cands) && cands[end].priority == cands[start].priority {
			end++
		}
		ordered = append(ordered, orderByWeight(cands[start:end])...)
		start = end
	}
	return ordered
}

// orderByWeight implements the RFC 2782 weight algorithm: weight 0 entries
// first, pick a uniform random number in [0, sum], select the first entry whose
// running sum is >= that number, and repeat for the remaining entries.
func orderByWeight(group []candidate) []candidate {
	remaining := append([]candidate(nil), group...)
	shuffle(len(remaining), func(i, j int) { remaining[i], remaining[j] = remaining[j], remaining[i] })
	sort.SliceStable(remaining, func(i, j int) bool {
		return remaining[i].capacity == 0 && remaining[j].capacity != 0
	})

	ordered := make([]candidate, 0, len(group))
	for len(remaining) > 0 {
		var sum int64
		for _, c := range remaining {
			sum += int64(c.capacity)
		}
		r := randInt63n(sum + 1)
		pick := len(remaining) - 1
		var running int64
		for i, c := range remaining {
			running += int64(c.capacity)
			if running >= r {
				pick = i
				break
			}
		}
		ordered = append(ordered, remaining[pick])
		remaining = append(remaining[:pick], remaining[pick+1:]...)
	}
	return ordered
}

func toProfile(p *models.Nrf_NFDisc_NFProfile) Profile {
	return Profile{
		FQDN:          p.Fqdn,
		InterPlmnFQDN: p.InterPlmnFqdn,
		IPv4Addresses: p.Ipv4Addresses,
		IPv6Addresses: p.Ipv6Addresses,
	}
}

func toService(s models.Nrf_NFDisc_NFService) Service {
	endpoints := make([]Endpoint, 0, len(s.IpEndPoints))
	for _, ep := range s.IpEndPoints {
		endpoints = append(endpoints, Endpoint{IPv4Address: ep.Ipv4Address, IPv6Address: ep.Ipv6Address, Port: ep.Port})
	}
	return Service{
		Scheme:        string(s.Scheme),
		FQDN:          s.Fqdn,
		InterPlmnFQDN: s.InterPlmnFqdn,
		APIPrefix:     s.ApiPrefix,
		Endpoints:     endpoints,
	}
}
