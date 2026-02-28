package pfcp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOuterHeaderCreation_GTPUv4(t *testing.T) {
	// Bit 1 (0x01) = GTP-U/UDP/IPv4 -> TEID + IPv4
	payload := []byte{
		0x01, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.True(t, f.HasTEID())
	assert.True(t, f.HasIPv4())
	assert.Equal(t, uint32(1), f.TEID)
	assert.Equal(t, net.IP{10, 0, 0, 1}, f.IPv4Address)
}

func TestParseOuterHeaderCreation_GTPUv6(t *testing.T) {
	// Bit 2 (0x02) = GTP-U/UDP/IPv6 -> TEID + IPv6
	payload := []byte{
		0x02, 0x00,
		0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.True(t, f.HasTEID())
	assert.False(t, f.HasIPv4())
	assert.Equal(t, uint32(2), f.TEID)
	assert.Equal(t, net.ParseIP("::1"), f.IPv6Address)
}

func TestParseOuterHeaderCreation_UDPv4(t *testing.T) {
	// Bit 3 (0x04) = UDP/IPv4 -> IPv4 + Port
	payload := []byte{
		0x04, 0x00,
		0xc0, 0xa8, 0x01, 0x01,
		0x1f, 0x90,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.False(t, f.HasTEID())
	assert.True(t, f.HasIPv4())
	assert.Equal(t, net.IP{192, 168, 1, 1}, f.IPv4Address)
	assert.Equal(t, uint16(8080), f.PortNumber)
}

func TestParseOuterHeaderCreation_UDPv6(t *testing.T) {
	// Bit 4 (0x08) = UDP/IPv6 -> IPv6 + Port
	payload := []byte{
		0x08, 0x00,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x50,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.False(t, f.HasTEID())
	assert.False(t, f.HasIPv4())
	assert.Equal(t, net.ParseIP("2001:db8::1"), f.IPv6Address)
	assert.Equal(t, uint16(80), f.PortNumber)
}

func TestParseOuterHeaderCreation_IPv4Only(t *testing.T) {
	// Bit 5 (0x10) = IPv4 only
	payload := []byte{
		0x10, 0x00,
		0xac, 0x10, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.False(t, f.HasTEID())
	assert.True(t, f.HasIPv4())
	assert.Equal(t, net.IP{172, 16, 0, 1}, f.IPv4Address)
}

func TestParseOuterHeaderCreation_IPv6Only(t *testing.T) {
	// Bit 6 (0x20) = IPv6 only
	payload := []byte{
		0x20, 0x00,
		0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.False(t, f.HasTEID())
	assert.False(t, f.HasIPv4())
	assert.Equal(t, net.ParseIP("fe80::1"), f.IPv6Address)
}

func TestParseOuterHeaderCreation_CTag(t *testing.T) {
	// Bit 1 + Bit 7 (0x41) = GTP-U/UDP/IPv4 + C-TAG
	payload := []byte{
		0x41, 0x00,
		0x00, 0x00, 0x00, 0x05,
		0x0a, 0x00, 0x00, 0x01,
		0x12, 0x34, 0x56,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.Equal(t, uint32(5), f.TEID)
	assert.Equal(t, net.IP{10, 0, 0, 1}, f.IPv4Address)
	assert.Equal(t, uint32(0x123456), f.CTag)
}

func TestParseOuterHeaderCreation_STag(t *testing.T) {
	// Bit 1 + Bit 8 (0x81) = GTP-U/UDP/IPv4 + S-TAG
	payload := []byte{
		0x81, 0x00,
		0x00, 0x00, 0x00, 0x0a,
		0xc0, 0xa8, 0x00, 0x01,
		0xab, 0xcd, 0xef,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.Equal(t, uint32(10), f.TEID)
	assert.Equal(t, net.IP{192, 168, 0, 1}, f.IPv4Address)
	assert.Equal(t, uint32(0xabcdef), f.STag)
}

func TestParseOuterHeaderCreation_CTagAndSTag(t *testing.T) {
	// Bit 1 + Bit 7 + Bit 8 = GTP-U/UDP/IPv4 + C-TAG + S-TAG
	payload := []byte{
		0xc1, 0x00,
		0x00, 0x00, 0x00, 0x0f,
		0x0a, 0x01, 0x02, 0x03,
		0xaa, 0xbb, 0xcc,
		0xdd, 0xee, 0xff,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)
	assert.Equal(t, uint32(15), f.TEID)
	assert.Equal(t, net.IP{10, 1, 2, 3}, f.IPv4Address)
	assert.Equal(t, uint32(0xaabbcc), f.CTag)
	assert.Equal(t, uint32(0xddeeff), f.STag)
}

func TestParseOuterHeaderCreation_TooShort(t *testing.T) {
	_, err := ParseOuterHeaderCreation([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too short")

	_, err = ParseOuterHeaderCreation([]byte{0x01})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too short")
}

func TestParseOuterHeaderCreation_InsufficientTEID(t *testing.T) {
	payload := []byte{0x01, 0x00, 0x00, 0x00}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TEID")
}

func TestParseOuterHeaderCreation_InsufficientIPv4(t *testing.T) {
	payload := []byte{0x10, 0x00, 0x0a, 0x00}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "IPv4")
}

func TestParseOuterHeaderCreation_InsufficientIPv6(t *testing.T) {
	payload := []byte{0x20, 0x00, 0x00, 0x01}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "IPv6")
}

func TestParseOuterHeaderCreation_InsufficientPort(t *testing.T) {
	payload := []byte{
		0x04, 0x00,
		0x0a, 0x00, 0x00, 0x01,
		0x00,
	}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Port")
}

func TestParseOuterHeaderCreation_InsufficientCTag(t *testing.T) {
	payload := []byte{
		0x41, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x01,
		0xaa, 0xbb,
	}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "C-TAG")
}

func TestParseOuterHeaderCreation_InsufficientSTag(t *testing.T) {
	payload := []byte{
		0x81, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x01,
		0xdd, 0xee,
	}
	_, err := ParseOuterHeaderCreation(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "S-TAG")
}

func TestParseOuterHeaderCreation_MemoryIndependence(t *testing.T) {
	payload := []byte{
		0x01, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)

	payload[6] = 0xff
	payload[7] = 0xff

	assert.Equal(t, net.IP{10, 0, 0, 1}, f.IPv4Address)
}

func TestParseOuterHeaderCreation_MemoryIndependenceIPv6(t *testing.T) {
	payload := []byte{
		0x02, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}

	f, err := ParseOuterHeaderCreation(payload)
	require.NoError(t, err)

	expected := net.ParseIP("::1")
	payload[6] = 0xff

	assert.Equal(t, expected, f.IPv6Address)
}

func TestHasTEID(t *testing.T) {
	tests := []struct {
		name string
		desc uint16
		want bool
	}{
		{"bit1 set", 0x0100, true},
		{"bit2 set", 0x0200, true},
		{"bit1+bit2", 0x0300, true},
		{"no teid bits", 0x0400, false},
		{"zero", 0x0000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &OuterHeaderCreationFields{OuterHeaderCreationDescription: tt.desc}
			assert.Equal(t, tt.want, f.HasTEID())
		})
	}
}

func TestHasIPv4(t *testing.T) {
	tests := []struct {
		name string
		desc uint16
		want bool
	}{
		{"bit1 set", 0x0100, true},
		{"bit3 set", 0x0400, true},
		{"bit5 set", 0x1000, true},
		{"bit2 only", 0x0200, false},
		{"zero", 0x0000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &OuterHeaderCreationFields{OuterHeaderCreationDescription: tt.desc}
			assert.Equal(t, tt.want, f.HasIPv4())
		})
	}
}
