package ippool

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type IPPool struct {
	IPSubnet *net.IPNet
	Pool     *LazyReusePool
}

func NewIPPool(cidr string) (*IPPool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, errors.Wrapf(err, "NewIPPool ParseCIDR")
	}

	minAddr, maxAddr, err := calcAddrRange(ipNet)
	if err != nil {
		return nil, errors.Wrapf(err, "NewIPPool calcAddrRange")
	}

	newPool, err := NewLazyReusePool(int(minAddr), int(maxAddr))
	if err != nil {
		return nil, errors.Wrapf(err, "NewIPPool NewLazyReusePool")
	}

	if err := newPool.Reserve(int(minAddr), int(minAddr)); err != nil {
		return nil, errors.Wrapf(err, "Remove network id from pool failed for %s", cidr)
	}
	if err := newPool.Reserve(int(maxAddr), int(maxAddr)); err != nil {
		return nil, errors.Wrapf(err, "Remove broadcasting address from pool failed for %s", cidr)
	}

	return &IPPool{IPSubnet: ipNet, Pool: newPool}, nil
}

func calcAddrRange(ipNet *net.IPNet) (minAddr, maxAddr uint32, err error) {
	maskVal := binary.BigEndian.Uint32(ipNet.Mask)
	baseIPVal := binary.BigEndian.Uint32(ipNet.IP)
	// move removing network and broadcast address later
	minAddr = (baseIPVal & maskVal)
	maxAddr = (baseIPVal | ^maskVal)
	if minAddr >= maxAddr {
		return minAddr, maxAddr, errors.New("mask is invalid")
	}
	return minAddr, maxAddr, nil
}

func (p *IPPool) Reallocate(request net.IP) (net.IP, bool) {
	var allocVal int
	if request == nil {
		return nil, false
	}
	allocVal = int(binary.BigEndian.Uint32(request))
	ok := p.Pool.Use(allocVal)
	inUsed := !ok
	return uint32ToIP(uint32(allocVal)), inUsed // #nosec G115
}

func (p *IPPool) Allocate(request net.IP) (net.IP, error) {
	var allocVal int
	var ok bool
	if request != nil {
		allocVal = int(binary.BigEndian.Uint32(request))
		ok = p.Pool.Use(allocVal)
		if !ok {
			return nil, errors.Errorf("IP[%s] is used in Pool[%+v]", request, p.IPSubnet)
		}
		// if allocated request IP address
		goto RETURNIP
	}

	allocVal, ok = p.Pool.Allocate()
	if !ok {
		return nil, errors.Errorf("Pool is empty: %+v", p.IPSubnet)
	}

RETURNIP:
	retIP := uint32ToIP(uint32(allocVal)) // #nosec G115
	return retIP, nil
}

func (p *IPPool) Exclude(excludePool *IPPool) error {
	excludeMin := excludePool.Pool.Min()
	excludeMax := excludePool.Pool.Max()
	if err := p.Pool.Reserve(excludeMin, excludeMax); err != nil {
		return errors.Errorf("exclude uePool fail: %v", err)
	}
	return nil
}

func uint32ToIP(intval uint32) net.IP {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, intval)
	return buf
}

func (p *IPPool) Release(ip net.IP) error {
	if len(ip) < net.IPv4len {
		return errors.Errorf("failed to release invalid Address: %s", ip)
	}
	addrVal := binary.BigEndian.Uint32(ip)
	res := p.Pool.Free(int(addrVal))
	if !res {
		return errors.Errorf("failed to release UE Address: %s", ip)
	}
	return nil
}

func (p *IPPool) String() string {
	str := "["
	elements := p.Pool.Dump()
	for index, element := range elements {
		var firstAddr net.IP
		var lastAddr net.IP
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[0])) // #nosec G115
		firstAddr = buf
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[1])) // #nosec G115
		lastAddr = buf
		if index > 0 {
			str += ("->")
		}
		str += fmt.Sprintf("{%s - %s}", firstAddr.String(), lastAddr.String())
	}
	str += ("]")
	return str
}
