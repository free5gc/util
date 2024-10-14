package safe_channel

import "sync"

type SafeCh[T any] struct {
	mu     sync.Mutex
	closed bool
	ch     chan T
}

func NewSafeCh[T any](size int) *SafeCh[T] {
	return &SafeCh[T]{
		ch: make(chan T, size),
	}
}

func (c *SafeCh[T]) Send(e T) {
	c.mu.Lock()
	defer c.mu.Unlock()

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
