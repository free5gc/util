package safe_channel

import "sync"

type SafeCh[T any] struct {
	mu     sync.RWMutex
	closed bool
	ch     chan T
}

func NewSafeCh[T any](size int) *SafeCh[T] {
	return &SafeCh[T]{
		ch: make(chan T, size),
	}
}

func (c *SafeCh[T]) Send(e T) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.closed {
		c.ch <- e
	}
}

func (c *SafeCh[T]) GetRcvChan() <-chan T {
	return c.ch
}

func (c *SafeCh[T]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		close(c.ch)
		c.closed = true
	}
}

func (c *SafeCh[T]) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}
