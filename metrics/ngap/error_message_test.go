package ngap

import (
	"strings"
	"testing"

	"github.com/free5gc/ngap/aper"
	ngapie "github.com/free5gc/ngap/ie"
)

// TestCauseStrCoverage fails when github.com/free5gc/ngap gains a cause value
// that the maps of this package do not carry, which would otherwise silently
// start reporting the new value as UNKNOWN_CAUSE_ERR.
func TestCauseStrCoverage(t *testing.T) {
	testCases := []struct {
		group    string
		upstream map[aper.Enumerated]string
		local    map[aper.Enumerated]string
	}{
		{"RadioNetwork", ngapie.CauseRadioNetworkStrings, causeRadioNetworkStr},
		{"Transport", ngapie.CauseTransportStrings, causeTransportStr},
		{"Nas", ngapie.CauseNasStrings, causeNasStr},
		{"Protocol", ngapie.CauseProtocolStrings, causeProtocolStr},
		{"Misc", ngapie.CauseMiscStrings, causeMiscStr},
	}

	for _, testCase := range testCases {
		t.Run(testCase.group, func(t *testing.T) {
			for value := range testCase.upstream {
				if _, ok := testCase.local[value]; !ok {
					t.Errorf("cause value %d of group %s has no entry in the local map",
						value, testCase.group)
				}
			}
		})
	}
}

// TestGetCauseErrorStr pins the label values that dashboards and alerting rules
// depend on, and checks that an unknown value cannot leak into the label.
func TestGetCauseErrorStr(t *testing.T) {
	testCases := []struct {
		name     string
		cause    *ngapie.Cause
		expected string
	}{
		{
			name:     "nil cause",
			cause:    nil,
			expected: UNKNOWN_NGAP_TYPE_CAUSE_ERR,
		},
		{
			name:     "nil choice",
			cause:    &ngapie.Cause{},
			expected: UNKNOWN_NGAP_TYPE_CAUSE_ERR,
		},
		{
			name: "radio network",
			cause: &ngapie.Cause{Choice: &ngapie.CauseRadioNetwork{
				Value: ngapie.CauseRadioNetworkPresentRadioConnectionWithUeLost,
			}},
			expected: "RadioNetwork : RadioConnectionWithUeLost",
		},
		{
			name: "transport",
			cause: &ngapie.Cause{Choice: &ngapie.CauseTransport{
				Value: ngapie.CauseTransportPresentTransportResourceUnavailable,
			}},
			expected: "Transport : TransportResourceUnavailable",
		},
		{
			name: "nas",
			cause: &ngapie.Cause{Choice: &ngapie.CauseNas{
				Value: ngapie.CauseNasPresentDeregister,
			}},
			expected: "Nas : Deregister",
		},
		{
			name: "protocol",
			cause: &ngapie.Cause{Choice: &ngapie.CauseProtocol{
				Value: ngapie.CauseProtocolPresentSemanticError,
			}},
			expected: "Protocol : SemanticError",
		},
		{
			name: "misc",
			cause: &ngapie.Cause{Choice: &ngapie.CauseMisc{
				Value: ngapie.CauseMiscPresentUnknownPLMNOrSNPN,
			}},
			expected: "Misc : UnknownPLMN",
		},
		{
			name:     "choice extensions",
			cause:    &ngapie.Cause{Choice: &ngapie.ProtocolIESingleContainerCauseExtIEs{}},
			expected: CHOICE_EXTENSIONS_UNKNOWN_ERR,
		},
		{
			name:     "value outside of the enumeration",
			cause:    &ngapie.Cause{Choice: &ngapie.CauseRadioNetwork{Value: aper.Enumerated(9999)}},
			expected: UNKNOWN_CAUSE_ERR,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if str := GetCauseErrorStr(testCase.cause); str != testCase.expected {
				t.Errorf("expected cause: %q, output cause: %q", testCase.expected, str)
			}
		})
	}
}

// TestCauseLabelCardinalityIsBounded walks the whole value space of every cause
// group: no input may produce a label value carrying the numeric cause value,
// the way ngapie.Cause.String() does for the values it does not know.
func TestCauseLabelCardinalityIsBounded(t *testing.T) {
	choices := []func(aper.Enumerated) ngapie.CauseAlt{
		func(v aper.Enumerated) ngapie.CauseAlt { return &ngapie.CauseRadioNetwork{Value: v} },
		func(v aper.Enumerated) ngapie.CauseAlt { return &ngapie.CauseTransport{Value: v} },
		func(v aper.Enumerated) ngapie.CauseAlt { return &ngapie.CauseNas{Value: v} },
		func(v aper.Enumerated) ngapie.CauseAlt { return &ngapie.CauseProtocol{Value: v} },
		func(v aper.Enumerated) ngapie.CauseAlt { return &ngapie.CauseMisc{Value: v} },
	}

	labels := make(map[string]struct{})

	for _, choice := range choices {
		for value := range 256 {
			str := GetCauseErrorStr(&ngapie.Cause{Choice: choice(aper.Enumerated(value))})

			if strings.Contains(str, "Unregconized") {
				t.Errorf("cause value %d produced the upstream fallback %q", value, str)
			}

			labels[str] = struct{}{}
		}
	}

	// 57 + 2 + 5 + 7 + 6 known values, plus UNKNOWN_CAUSE_ERR.
	if expected := 78; len(labels) != expected {
		t.Errorf("expected %d distinct label values, output %d", expected, len(labels))
	}
}
