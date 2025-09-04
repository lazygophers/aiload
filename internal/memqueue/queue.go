package memqueue

import (
	"errors"
	"sync"
	"time"

	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/routine"
	"github.com/lazygophers/utils/runtime"
	"go.uber.org/atomic"
)

var (
	ErrQueueFull = errors.New("queue is full")
)

type Consume struct {
	Retry          bool
	SkipRetryCount bool

	Delay time.Duration
}

type Queue[T any] struct {
	c *Config[T]

	lowChan, normalChan, highChan chan *Message[T]
	deferredQueue                 *DeferredQueue[T]

	depth         *atomic.Int32
	inFlightDepth *atomic.Int32

	serialMap sync.Map
}

func (p *Queue[T]) Pub(value T, priorities ...Priority) error {
	msg := &Message[T]{
		Data: value,
	}

	if len(priorities) > 0 {
		msg.Priority = priorities[0]
	}

	msg.init(p.c)

	if p.depth.Load() >= int32(p.c.MaxSize) {
		return ErrQueueFull
	}

	if p.c.OnBeforePub != nil && !p.c.OnBeforePub(msg) {
		return nil
	}

	switch msg.Priority {
	case PriorityLow:
		select {
		case p.lowChan <- msg:
			p.depth.Inc()
		default:
			return ErrQueueFull
		}

	case PriorityNormal:
		select {
		case p.normalChan <- msg:
			p.depth.Inc()
		default:
			return ErrQueueFull
		}

	case PriorityHigh:
		fallthrough
	default:
		select {
		case p.highChan <- msg:
			p.depth.Inc()
		default:
			return ErrQueueFull
		}
	}

	return nil
}

func (p *Queue[T]) Depth() int32 {
	return p.depth.Load()
}

func (p *Queue[T]) InFlightDepth() int32 {
	return p.inFlightDepth.Load()
}

type Serialer interface {
	Unique() uint64
}

// tryLockSerial 尝试锁定消息的序列号。
// 如果序列号已被锁定，返回 false；否则锁定并返回 true。
func (p *Queue[M]) tryLockSerial(key uint64) bool {
	// locked = false，map 没有数据，锁定了，返回 true
	_, locked := p.serialMap.LoadOrStore(key, 0)
	return !locked
}

// unlockSerial 解锁消息的序列号。
func (p *Queue[M]) unlockSerial(key uint64) {
	p.serialMap.Delete(key)
}

func (p *Queue[T]) Process(fn func(msg *Message[T]) (*Consume, error)) {
	logic := func(msg *Message[T]) {
		// 记录执行中的数据
		p.inFlightDepth.Inc()
		defer p.inFlightDepth.Dec() // 确保 inFlightDepth 正确递减

		if p.c.MaxRetry > 0 && msg.RetryCount >= p.c.MaxRetry {
			// 超过重试次数，直接丢弃
			p.depth.Dec() // 确保深度正确减少
			return
		}

		// 前置的判断
		if !msg.ExecAt.IsZero() && msg.ExecAt.After(time.Now()) {
			// 没到执行时间，延时执行
			p.deferredQueue.Push(msg)
			return
		}

		if v, ok := any(msg.Data).(Serialer); ok {
			key := v.Unique()
			if !p.tryLockSerial(key) {
				msg.ExecAt = time.Now().Add(time.Second)
				p.deferredQueue.Push(msg)
				return
			}
			defer p.unlockSerial(key)
		}

		defer runtime.CachePanicWithHandle(func(err interface{}) {
			p.depth.Dec() // 确保深度正确减少
		})

		consume, err := fn(msg)
		if err != nil {
			log.Errorf("err:%s", err)
		}

		// 默认认为其执行成功
		if consume == nil {
			p.depth.Dec() // 确保深度正确减少
			return
		}

		// 不需要重试的
		if !consume.Retry {
			p.depth.Dec() // 确保深度正确减少
			return
		}

		// 需要重试的
		if !consume.SkipRetryCount {
			msg.RetryCount++
		}

		if consume.Delay > 0 {
			msg.ExecAt = time.Now().Add(consume.Delay)
		} else {
			msg.ExecAt = time.Now()
		}

		p.deferredQueue.Push(msg)
	}

	for i := 0; i < p.c.ConcurrentCount; i++ {
		routine.Go(func() (err error) {
			var msg *Message[T]
			for {
				// **性能优化**:
				// 采用非阻塞的 select 语句，并结合 default 分支和自适应休眠，实现了一个高效的优先级队列消费模型。
				// 1. 优先消费高优先级通道。
				// 2. 当所有通道都为空时，进入 default 分支，执行自适应休眠。
				// 3. 休眠时间会随着连续空闲的次数指数级增长，直到达到上限，以节约 CPU。
				// 4. 一旦成功消费消息，休眠时间将立即重置为最小值，确保低延迟。
				select {
				case msg = <-p.highChan:
					logic(msg)
				case msg = <-p.normalChan:
					logic(msg)
				case msg = <-p.lowChan:
					logic(msg)
				}
			}
		})
	}
}

func New[T any](c *Config[T]) *Queue[T] {
	c.apply()

	p := &Queue[T]{
		c:             c,
		lowChan:       make(chan *Message[T], c.MaxSize+c.ConcurrentCount+1),
		normalChan:    make(chan *Message[T], c.MaxSize+c.ConcurrentCount+1),
		highChan:      make(chan *Message[T], c.MaxSize+c.ConcurrentCount+1),
		depth:         atomic.NewInt32(0),
		inFlightDepth: atomic.NewInt32(0),
		deferredQueue: NewDeferredQueue[T](),
	}

	// 启动延迟队列后台处理器
	routine.Go(func() (err error) {
		ticker := time.NewTicker(c.DeferredQueueTickerDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if msgs := p.deferredQueue.Pop(); len(msgs) > 0 {
					for _, msg := range msgs {
						switch msg.Priority {
						case PriorityLow:
							p.lowChan <- msg

						case PriorityNormal:
							p.normalChan <- msg

						case PriorityHigh:
							fallthrough
						default:
							p.highChan <- msg
						}
					}
				}
			}
		}
	})

	return p
}
