package memqueue

import (
	"github.com/lazygophers/utils/cryptox"
	"time"
)

// 优先级
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
)

// Message 是一个通用的消息结构体，用于在内存队列中存储和传递数据。
// 它包含了消息的基本信息、重试机制、过期时间、执行时间和执行超时等属性。
type Message[T any] struct {
	// CreatedAt 表示消息的创建时间，用于记录消息的生成时间点。
	CreatedAt time.Time

	// Id 是消息的唯一标识符，用于区分不同的消息。
	Id string
	// Data 是消息携带的数据，类型为泛型 T，可以存储任意类型的数据。
	Data T

	// Priority 表示消息的优先级，用于确定消息的调度顺序。
	Priority Priority

	// RetryCount 表示消息当前已重试的次数，用于跟踪消息的重试状态。
	RetryCount int
	// ExecAt 表示消息的执行时间，表示该消息应该被执行的时间点。
	ExecAt time.Time
}

func (p *Message[T]) init(c *Config[T]) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	if p.Id == "" {
		p.Id = cryptox.UUID()
	}

	if p.ExecAt.IsZero() && c.Delay > 0 {
		p.ExecAt = time.Now().Add(c.Delay)
	}

	p.RetryCount = 0
}
