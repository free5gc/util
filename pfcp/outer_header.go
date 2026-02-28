package pfcp

import (
	"errors"
	"fmt"
	"net"
)

// OuterHeaderCreationFields represents parsed Outer Header Creation IE fields.
// This local implementation explicitly handles 3-byte C-TAG/S-TAG fields
// as defined for the Outer Header Creation IE in 3GPP TS 29.244, Figure 8.2.56-1.
//
// The go-pfcp library (github.com/wmnsk/go-pfcp) incorrectly uses Uint32 (4 bytes)
// to read 3-byte C-TAG/S-TAG fields, causing payload misalignment.
type OuterHeaderCreationFields struct {
	OuterHeaderCreationDescription uint16
	TEID                           uint32
	IPv4Address                    net.IP
	IPv6Address                    net.IP
	PortNumber                     uint16
	CTag                           uint32
	STag                           uint32
}

// Bit flags in Octet 5 (payload[0]) determine which fields are present:
//
//	Bit 1 (0x01): GTP-U/UDP/IPv4  → TEID + IPv4
//	Bit 2 (0x02): GTP-U/UDP/IPv6  → TEID + IPv6
//	Bit 3 (0x04): UDP/IPv4        → IPv4 + Port
//	Bit 4 (0x08): UDP/IPv6        → IPv6 + Port
//	Bit 5 (0x10): IPv4            → IPv4
//	Bit 6 (0x20): IPv6            → IPv6
//	Bit 7 (0x40): C-TAG           → 3-octet C-TAG
//	Bit 8 (0x80): S-TAG           → 3-octet S-TAG
func ParseOuterHeaderCreation(payload []byte) (*OuterHeaderCreationFields, error) {
	l := len(payload)
	if l < 2 {
		return nil, errors.New("OuterHeaderCreation payload too short: need at least 2 bytes")
	}

	f := &OuterHeaderCreationFields{}
	f.OuterHeaderCreationDescription = uint16(payload[0])<<8 | uint16(payload[1])
	offset := 2

	// oct5 is the first byte containing bit flags
	oct5 := payload[0]

	// TEID: present if bit 1 (0x01) or bit 2 (0x02) is set
	if (oct5&0x01) != 0 || (oct5&0x02) != 0 {
		if l < offset+4 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for TEID at offset %d", offset)
		}
		f.TEID = uint32(payload[offset])<<24 | uint32(payload[offset+1])<<16 |
			uint32(payload[offset+2])<<8 | uint32(payload[offset+3])
		offset += 4
	}

	// IPv4: present if bit 1 (0x01), bit 3 (0x04), or bit 5 (0x10) is set
	if (oct5&0x01) != 0 || (oct5&0x04) != 0 || (oct5&0x10) != 0 {
		if l < offset+4 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for IPv4 at offset %d", offset)
		}
		f.IPv4Address = make(net.IP, 4)
		copy(f.IPv4Address, payload[offset:offset+4])
		offset += 4
	}

	// IPv6: present if bit 2 (0x02), bit 4 (0x08), or bit 6 (0x20) is set
	if (oct5&0x02) != 0 || (oct5&0x08) != 0 || (oct5&0x20) != 0 {
		if l < offset+16 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for IPv6 at offset %d", offset)
		}
		f.IPv6Address = make(net.IP, 16)
		copy(f.IPv6Address, payload[offset:offset+16])
		offset += 16
	}

	// Port Number: present if bit 3 (0x04) or bit 4 (0x08) is set
	if (oct5&0x04) != 0 || (oct5&0x08) != 0 {
		if l < offset+2 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for Port at offset %d", offset)
		}
		f.PortNumber = uint16(payload[offset])<<8 | uint16(payload[offset+1])
		offset += 2
	}

	// C-TAG: present if bit 7 (0x40) is set
	if (oct5 & 0x40) != 0 {
		if l < offset+3 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for C-TAG at offset %d", offset)
		}
		f.CTag = uint32(payload[offset])<<16 | uint32(payload[offset+1])<<8 | uint32(payload[offset+2])
		offset += 3
	}

	// S-TAG: present if bit 8 (0x80) is set
	if (oct5 & 0x80) != 0 {
		if l < offset+3 {
			return nil, fmt.Errorf("OuterHeaderCreation: insufficient bytes for S-TAG at offset %d", offset)
		}
		f.STag = uint32(payload[offset])<<16 | uint32(payload[offset+1])<<8 | uint32(payload[offset+2])
	}

	return f, nil
}

// HasTEID checks if OuterHeaderCreation has TEID field
func (f *OuterHeaderCreationFields) HasTEID() bool {
	oct5 := uint8((f.OuterHeaderCreationDescription & 0xff00) >> 8)
	return (oct5&0x01) != 0 || (oct5&0x02) != 0
}

// HasIPv4 checks if OuterHeaderCreation has IPv4 field
func (f *OuterHeaderCreationFields) HasIPv4() bool {
	oct5 := uint8((f.OuterHeaderCreationDescription & 0xff00) >> 8)
	return (oct5&0x01) != 0 || (oct5&0x04) != 0 || (oct5&0x10) != 0
}
