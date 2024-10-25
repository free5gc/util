package safe_channel

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Simulate 1 receiver and N(N=2) senders situation
func TestSafeChannel(t *testing.T) {
	sCh := NewSafeCh[int](1)
	wg := sync.WaitGroup{}

	// Two senders
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			if i == 0 {
				// Case: send after sCh closed
				time.Sleep(1 * time.Second)

				require.Equal(t, sCh.IsClosed(), true)
				sCh.Send(1) // No panic
				require.Equal(t, sCh.IsClosed(), true)
			} else {
				// Case: send success
				sCh.Send(1)
				require.Equal(t, sCh.IsClosed(), false)
			}
			wg.Done()
		}(i)
	}

	// One receiver
	<-sCh.GetRcvChan()
	sCh.Close()
	require.Equal(t, sCh.IsClosed(), true)

	wg.Wait()
}
