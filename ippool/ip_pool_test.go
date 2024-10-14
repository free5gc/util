package ippool

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIPPool(t *testing.T) {
	var ipPool *IPPool
	var err error

	// invalid CIDR
	_, err = NewIPPool("10.10.0.1")
	require.Error(t, err)

	// invalid MaskLen
	_, err = NewIPPool("10.10.0.1/32")
	require.Error(t, err)

	// valid CIDR
	ipPool, err = NewIPPool("10.10.0.0/24")
	require.NoError(t, err)

	var allocIP net.IP

	// make allowed ip pools
	var ipPoolList []net.IP
	for i := 1; i <= 254; i += 1 {
		ipStr := fmt.Sprintf("10.10.0.%d", i)
		ipPoolList = append(ipPoolList, net.ParseIP(ipStr).To4())
	}

	// allocate
	for i := 1; i <= 254; i += 1 {
		allocIP, err = ipPool.Allocate(nil)
		require.NoError(t, err)
		require.Contains(t, ipPoolList, allocIP)
	}

	// Re-allocate the last one allocated IP
	reAllocIP, used := ipPool.Reallocate(allocIP)
	require.True(t, used)
	require.Equal(t, allocIP, reAllocIP)

	// ip pool is empty
	allocIP, err = ipPool.Allocate(nil)
	require.Error(t, err)
	require.Nil(t, allocIP)

	// reserve some IP randomly pick from list
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec
	var ipResPoolList []net.IP
	for n, i := range rng.Perm(254) {
		fmt.Println("Reserved:", ipPoolList[i].String())
		ipResPoolList = append(ipResPoolList, ipPoolList[i])
		if n > 100 {
			break
		}
	}
	fmt.Println("Reserve total:", len(ipResPoolList))

	// release all ip
	fmt.Println("Total:", len(ipPoolList))
	for i := 0; i < len(ipPoolList); i += 1 {
		err = ipPool.Release(ipPoolList[i])
		require.NoError(t, err)
	}

	fmt.Println("Allocate specified IP in IP pool ...")
	// allocate specify ip in reserved list
	for _, ip := range ipResPoolList {
		allocIP, err = ipPool.Allocate(ip)
		require.NoError(t, err)
		require.Equal(t, ip, allocIP)
	}
	fmt.Println(ipPool)

	fmt.Println("Allocate IP in free IP pool ...")

	// allocate ip in free pool
	for i := 0; i < (254 - len(ipResPoolList)); i += 1 {
		allocIP, err = ipPool.Allocate(nil)
		require.NoError(t, err)
		require.NotContains(t, ipResPoolList, allocIP)
	}

	// ip pool is empty
	allocIP, err = ipPool.Allocate(nil)
	require.Error(t, err)
	require.Nil(t, allocIP)
}

func TestIpPool_ExcludeRange(t *testing.T) {
	// from 0 to 255
	ipPool, err := NewIPPool("10.10.0.0/24")
	require.NoError(t, err)

	require.Equal(t, 0x0a0a0000, ipPool.Pool.Min())
	require.Equal(t, 0x0a0a00FF, ipPool.Pool.Max())
	// Not contains 0 (network id) & 255 (broadcast ip)
	require.Equal(t, 254, ipPool.Pool.Remain())

	// from 0 to 15
	excludeIPPool, err := NewIPPool("10.10.0.0/28")
	require.NoError(t, err)

	require.Equal(t, 0x0a0a0000, excludeIPPool.Pool.Min())
	require.Equal(t, 0x0a0a000F, excludeIPPool.Pool.Max())

	// Not contains 0 (network id) & 15 (broadcast ip)
	require.Equal(t, 14, excludeIPPool.Pool.Remain())

	err = ipPool.Exclude(excludeIPPool)
	require.NoError(t, err)

	// from 16 to 254
	require.Equal(t, 239, ipPool.Pool.Remain())

	for i := 16; i <= 254; i++ {
		allocate, err := ipPool.Allocate(nil)
		require.NoError(t, err)
		require.Equal(t, net.ParseIP(fmt.Sprintf("10.10.0.%d", i)).To4(), allocate)

		err = ipPool.Release(allocate)
		require.NoError(t, err)
	}
}
