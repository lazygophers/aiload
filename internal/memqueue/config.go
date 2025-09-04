package memqueue

import "time"

type Config[T any] struct {
	MaxSize int

	OnBeforePub func(msg *Message[T]) bool

	ConcurrentCount int
	MaxRetry        int
	Delay           time.Duration

	DeferredQueueTickerDuration time.Duration
}

func (p *Config[T]) apply() {
	if p.MaxSize <= 0 {
		p.MaxSize = 500
	}

	if p.ConcurrentCount <= 0 {
		p.ConcurrentCount = 1
	}

	if p.MaxRetry <= 0 {
		p.MaxRetry = 3
	}

	if p.DeferredQueueTickerDuration <= 0 {
		p.DeferredQueueTickerDuration = 10 * time.Second
	}
}
