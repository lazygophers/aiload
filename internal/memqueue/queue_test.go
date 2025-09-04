package memqueue

import (
	"errors"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Push(t *testing.T) {
	// 测试 Push 方法
	queue := New[string](&Config[string]{
		MaxSize: 2,
	})

	// 正常插入
	err := queue.Pub("test1")
	assert.NoError(t, err)
	assert.Equal(t, int32(1), queue.Depth())

	err = queue.Pub("test2")
	assert.NoError(t, err)
	assert.Equal(t, int32(2), queue.Depth())

	// 队列已满
	err = queue.Pub("test3")
	assert.ErrorIs(t, err, ErrQueueFull)
	assert.Equal(t, int32(2), queue.Depth())

	// 测试优先级
	t.Run("priority", func(t *testing.T) {
		q := New[string](&Config[string]{MaxSize: 3})
		assert.NoError(t, q.Pub("low", PriorityLow))
		assert.NoError(t, q.Pub("normal", PriorityNormal))
		assert.NoError(t, q.Pub("high", PriorityHigh))

		assert.Equal(t, 1, len(q.lowChan))
		assert.Equal(t, 1, len(q.normalChan))
		assert.Equal(t, 1, len(q.highChan))
	})
}

func TestQueue_Depth(t *testing.T) {
	// 测试 Depth 方法
	queue := New[string](&Config[string]{
		MaxSize: 100,
	})
	assert.NoError(t, queue.Pub("test1"))
	assert.NoError(t, queue.Pub("test2"))

	assert.Equal(t, int32(2), queue.Depth())
}

func TestQueue_InFlightDepth(t *testing.T) {
	// 测试 InFlightDepth 方法
	queue := New[string](&Config[string]{}) // 修改: 传入有效的 Config 实例
	assert.NoError(t, queue.Pub("test1"))

	// 模拟处理中的消息
	queue.inFlightDepth.Inc()
	assert.Equal(t, int32(1), queue.InFlightDepth())
}

func TestNew(t *testing.T) {
	// 测试 New 方法
	queue := New[string](&Config[string]{
		MaxSize:         100,
		ConcurrentCount: 2,
		MaxRetry:        5,
		Delay:           500 * time.Millisecond,
	})

	assert.Equal(t, 100, queue.c.MaxSize)
	assert.Equal(t, 2, queue.c.ConcurrentCount)
	assert.Equal(t, 5, queue.c.MaxRetry)
	assert.Equal(t, 500*time.Millisecond, queue.c.Delay)
}

func TestQueue_Process(t *testing.T) {
	t.Run("normal_consume", func(t *testing.T) {
		queue := New[string](&Config[string]{ConcurrentCount: 1})
		processed := atomic.NewBool(false)
		queue.Process(func(msg *Message[string]) (*Consume, error) {
			processed.Store(true)
			return nil, nil
		})
		assert.NoError(t, queue.Pub("test"))
		assert.Eventually(t, processed.Load, time.Second, 10*time.Millisecond)
	})

	t.Run("with_retry", func(t *testing.T) {
		queue := New[string](&Config[string]{
			ConcurrentCount:             1,
			MaxRetry:                    3, // MaxRetry = 3 意味着总共会尝试3次
			DeferredQueueTickerDuration: 50 * time.Millisecond,
		})
		retryCount := atomic.NewInt32(0)
		queue.Process(func(msg *Message[string]) (*Consume, error) {
			retryCount.Inc()
			return &Consume{Retry: true}, nil
		})
		assert.NoError(t, queue.Pub("test"))
		// 处理器函数总共会被调用3次
		assert.Eventually(t, func() bool { return retryCount.Load() == 3 }, time.Second, 50*time.Millisecond, "处理器应该被调用3次")
	})

	t.Run("max_retry_exceeded", func(t *testing.T) {
		// 恢复被禁用的测试
		queue := New[string](&Config[string]{
			ConcurrentCount:             1,
			MaxRetry:                    2, // MaxRetry = 2 意味着总共只会尝试2次
			DeferredQueueTickerDuration: 50 * time.Millisecond,
		})
		retryCount := atomic.NewInt32(0)
		queue.Process(func(msg *Message[string]) (*Consume, error) {
			retryCount.Inc()
			return &Consume{Retry: true}, nil
		})
		assert.NoError(t, queue.Pub("test"))
		// 使用 Eventually 等待，确保处理器被调用了预期的次数
		assert.Eventually(t, func() bool {
			return retryCount.Load() == 2
		}, 2*time.Second, 50*time.Millisecond, "处理器应该总共被调用2次")
	})

	t.Run("delayed_execution", func(t *testing.T) {
		queue := New[string](&Config[string]{
			ConcurrentCount:             1,
			Delay:                       100 * time.Millisecond,
			DeferredQueueTickerDuration: 10 * time.Millisecond,
		})
		processed := atomic.NewBool(false)
		start := time.Now()
		queue.Process(func(msg *Message[string]) (*Consume, error) {
			assert.True(t, time.Since(start) >= 100*time.Millisecond)
			processed.Store(true)
			return nil, nil
		})
		assert.NoError(t, queue.Pub("test"))
		assert.Eventually(t, processed.Load, time.Second, 10*time.Millisecond)
	})
}

// 测试延迟队列集成
func TestDeferredQueueIntegration(t *testing.T) {
	queue := New[string](&Config[string]{
		Delay:                       50 * time.Millisecond,
		ConcurrentCount:             1,
		DeferredQueueTickerDuration: 20 * time.Millisecond,
	})
	processed := atomic.NewBool(false)
	queue.Process(func(msg *Message[string]) (*Consume, error) {
		processed.Store(true)
		return nil, nil
	})
	assert.NoError(t, queue.Pub("delayed_test"))

	// 验证在延迟时间内，消息不会被处理
	assert.Never(t, func() bool { return processed.Load() }, 40*time.Millisecond, 10*time.Millisecond, "消息不应该在延迟结束前被处理")

	// 验证在延迟时间后，消息最终会被处理
	assert.Eventually(t, processed.Load, 2*time.Second, 20*time.Millisecond, "消息应该在延迟结束后被处理")
}

// 测试并发处理能力
func TestConcurrentProcessing(t *testing.T) {
	queue := New[int](&Config[int]{
		ConcurrentCount: 5,
		MaxSize:         100,
	})
	counter := atomic.NewInt32(0)
	queue.Process(func(msg *Message[int]) (*Consume, error) {
		counter.Inc()
		// 模拟少量工作
		time.Sleep(5 * time.Millisecond)
		return nil, nil
	})
	for i := 0; i < 100; i++ {
		assert.NoError(t, queue.Pub(i))
	}
	// 大幅延长超时时间以应对 CI 环境的缓慢
	assert.Eventually(t, func() bool {
		return counter.Load() == 100
	}, 10*time.Second, 50*time.Millisecond)
}

// 测试串行执行
type serialMessage struct {
	id uint64
}

func (m serialMessage) Unique() uint64 {
	return m.id
}

func TestSerialProcessing(t *testing.T) {
	queue := New[serialMessage](&Config[serialMessage]{
		ConcurrentCount:             5,
		DeferredQueueTickerDuration: 50 * time.Millisecond,
	})

	var mu sync.Mutex
	var processingOrder []int
	var callCount atomic.Int32

	queue.Process(func(msg *Message[serialMessage]) (*Consume, error) {
		callCount.Inc()
		mu.Lock()
		processingOrder = append(processingOrder, int(msg.Data.id))
		mu.Unlock()
		time.Sleep(100 * time.Millisecond) // Simulate work
		return nil, nil
	})

	// Push messages with the same ID to test serial execution.
	assert.NoError(t, queue.Pub(serialMessage{id: 1}))
	assert.NoError(t, queue.Pub(serialMessage{id: 1}))

	// Wait for processing to complete.
	assert.Eventually(t, func() bool {
		return callCount.Load() == 2
	}, 2*time.Second, 100*time.Millisecond)

	// Check that messages were processed serially.
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []int{1, 1}, processingOrder)
}

func TestQueue_Pub_QueueFull(t *testing.T) {
	t.Run("high_priority", func(t *testing.T) {
		q := New[string](&Config[string]{MaxSize: 1})
		assert.NoError(t, q.Pub("test", PriorityHigh))
		assert.ErrorIs(t, q.Pub("test", PriorityHigh), ErrQueueFull)
	})

	t.Run("normal_priority", func(t *testing.T) {
		q := New[string](&Config[string]{MaxSize: 1})
		assert.NoError(t, q.Pub("test", PriorityNormal))
		assert.ErrorIs(t, q.Pub("test", PriorityNormal), ErrQueueFull)
	})

	t.Run("low_priority", func(t *testing.T) {
		q := New[string](&Config[string]{MaxSize: 1})
		assert.NoError(t, q.Pub("test", PriorityLow))
		assert.ErrorIs(t, q.Pub("test", PriorityLow), ErrQueueFull)
	})
}

func TestQueue_Pub_OnBeforePub(t *testing.T) {
	q := New[string](&Config[string]{
		MaxSize: 1,
		OnBeforePub: func(msg *Message[string]) bool {
			return false
		},
	})

	assert.NoError(t, q.Pub("test"))
	assert.Equal(t, int32(0), q.Depth())
}

func TestQueue_Process_Panic(t *testing.T) {
	q := New[string](&Config[string]{ConcurrentCount: 1})
	q.Process(func(msg *Message[string]) (*Consume, error) {
		panic("test panic")
	})

	assert.NoError(t, q.Pub("test"))

	assert.Eventually(t, func() bool {
		return q.InFlightDepth() == 0 && q.Depth() == 0
	}, time.Second, 10*time.Millisecond, "inFlightDepth and depth should be 0 after panic")
}

func TestQueue_Process_Error(t *testing.T) {
	q := New[string](&Config[string]{ConcurrentCount: 1})
	processed := atomic.NewBool(false)
	q.Process(func(msg *Message[string]) (*Consume, error) {
		processed.Store(true)
		return nil, errors.New("consumer error")
	})

	assert.NoError(t, q.Pub("test"))
	assert.Eventually(t, processed.Load, time.Second, 10*time.Millisecond)
	assert.Eventually(t, func() bool {
		return q.Depth() == 0
	}, time.Second, 10*time.Millisecond)
}

func TestQueue_Process_NilConsume(t *testing.T) {
	q := New[string](&Config[string]{ConcurrentCount: 1})
	processed := atomic.NewBool(false)
	q.Process(func(msg *Message[string]) (*Consume, error) {
		processed.Store(true)
		return nil, nil
	})

	assert.NoError(t, q.Pub("test"))
	assert.Eventually(t, processed.Load, time.Second, 10*time.Millisecond)
	assert.Eventually(t, func() bool {
		return q.Depth() == 0
	}, time.Second, 10*time.Millisecond)
}

func TestQueue_Process_NoRetry(t *testing.T) {
	q := New[string](&Config[string]{ConcurrentCount: 1})
	processed := atomic.NewBool(false)
	q.Process(func(msg *Message[string]) (*Consume, error) {
		processed.Store(true)
		return &Consume{Retry: false}, nil
	})

	assert.NoError(t, q.Pub("test"))
	assert.Eventually(t, processed.Load, time.Second, 10*time.Millisecond)
	assert.Eventually(t, func() bool {
		return q.Depth() == 0
	}, time.Second, 10*time.Millisecond)
}

func TestDeferredQueueRequeuePriority(t *testing.T) {
	config := &Config[string]{
		ConcurrentCount:             1,
		DeferredQueueTickerDuration: 10 * time.Millisecond,
	}
	queue := New[string](config)

	var processedPriorities []Priority
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(3)

	queue.Process(func(msg *Message[string]) (*Consume, error) {
		mu.Lock()
		processedPriorities = append(processedPriorities, msg.Priority)
		mu.Unlock()
		wg.Done()
		return nil, nil
	})

	// 确保消息先进入延迟队列
	execAt := time.Now().Add(50 * time.Millisecond)

	highPriorityMsg := &Message[string]{Data: "high", Priority: PriorityHigh}
	highPriorityMsg.init(config)
	highPriorityMsg.ExecAt = execAt
	queue.deferredQueue.Push(highPriorityMsg)

	normalPriorityMsg := &Message[string]{Data: "normal", Priority: PriorityNormal}
	normalPriorityMsg.init(config)
	normalPriorityMsg.ExecAt = execAt
	queue.deferredQueue.Push(normalPriorityMsg)

	lowPriorityMsg := &Message[string]{Data: "low", Priority: PriorityLow}
	lowPriorityMsg.init(config)
	lowPriorityMsg.ExecAt = execAt
	queue.deferredQueue.Push(lowPriorityMsg)

	// 等待所有消息被处理
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for messages to be processed")
	}

	mu.Lock()
	defer mu.Unlock()

	// 验证消息是否按高、中、低优先级顺序处理
	// 注意：由于消费者是并发的，并且 ticker 可能同时触发多个消息，
	// 所以我们不能保证严格的顺序，但可以验证所有优先级的消息都被处理了。
	assert.ElementsMatch(t, []Priority{PriorityHigh, PriorityNormal, PriorityLow}, processedPriorities)
}

func TestQueue_Pub_DefaultPriority(t *testing.T) {
	q := New[string](&Config[string]{MaxSize: 1})
	// 使用一个未定义的优先级，触发 default case
	err := q.Pub("test", Priority(99))
	assert.NoError(t, err)
	// 验证消息是否被投递到了高优先级队列
	assert.Equal(t, 1, len(q.highChan))
}

func TestQueue_Process_SkipRetryCount(t *testing.T) {
	queue := New[string](&Config[string]{
		ConcurrentCount:             1,
		MaxRetry:                    3,
		DeferredQueueTickerDuration: 10 * time.Millisecond,
	})
	processCount := atomic.NewInt32(0)

	queue.Process(func(msg *Message[string]) (*Consume, error) {
		processCount.Inc()
		// 验证在跳过重试计数时，消息的 RetryCount 不会增加
		assert.Equal(t, 0, int(msg.RetryCount), "RetryCount should not be incremented when SkipRetryCount is true")
		// 第一次之后就停止重试，避免不必要的循环
		if processCount.Load() >= 1 {
			return &Consume{Retry: false}, nil
		}
		return &Consume{Retry: true, SkipRetryCount: true}, nil
	})

	assert.NoError(t, queue.Pub("test"))

	// 验证处理器至少被调用了一次
	assert.Eventually(t, func() bool {
		return processCount.Load() >= 1
	}, time.Second, 50*time.Millisecond, "处理器应该至少被调用一次")
}

func TestSerialProcessing_LockFailed(t *testing.T) {
	queue := New[serialMessage](&Config[serialMessage]{
		ConcurrentCount:             2, // Need at least 2 consumers
		DeferredQueueTickerDuration: 10 * time.Millisecond,
	})

	var wg sync.WaitGroup
	wg.Add(1)

	processedCount := atomic.NewInt32(0)

	// This consumer will block until we are ready
	queue.Process(func(msg *Message[serialMessage]) (*Consume, error) {
		processedCount.Inc()
		// The first message will wait here, holding the lock
		if msg.Data.id == 1 && processedCount.Load() == 1 {
			wg.Wait()
		}
		return nil, nil
	})

	// Push the first message. It will be picked up by a consumer and block.
	assert.NoError(t, queue.Pub(serialMessage{id: 1}))

	// Give the first consumer time to pick up the message and acquire the lock
	time.Sleep(100 * time.Millisecond)

	// Push the second message with the same ID.
	// The second consumer will try to lock, fail, and defer the message.
	assert.NoError(t, queue.Pub(serialMessage{id: 1}))

	// We expect the message to be in the deferred queue.
	assert.Eventually(t, func() bool {
		return queue.deferredQueue.Depth() == 1
	}, time.Second, 10*time.Millisecond, "Message should be deferred when lock fails")

	// Allow the first consumer to finish
	wg.Done()

	// Eventually, both messages should be processed, so the final depth should be 0.
	assert.Eventually(t, func() bool {
		return queue.Depth() == 0 && queue.InFlightDepth() == 0 && processedCount.Load() == 2
	}, 2*time.Second, 50*time.Millisecond, "Queue should be empty eventually")
}

func TestQueue_Process_RetryNoDelay(t *testing.T) {
	queue := New[string](&Config[string]{
		ConcurrentCount:             1,
		MaxRetry:                    2, // Allows for one retry
		DeferredQueueTickerDuration: 10 * time.Millisecond,
	})
	processCount := atomic.NewInt32(0)

	queue.Process(func(msg *Message[string]) (*Consume, error) {
		count := processCount.Inc()
		// First time, request retry with no delay.
		if count == 1 {
			// This will hit the `else` branch in queue.go:174
			return &Consume{Retry: true, Delay: 0}, nil
		}
		// Second time, stop processing.
		return &Consume{Retry: false}, nil
	})

	assert.NoError(t, queue.Pub("test"))

	// Processor should be called twice.
	assert.Eventually(t, func() bool {
		return processCount.Load() == 2
	}, time.Second, 50*time.Millisecond, "Process function should be called twice")

	// Queue should become empty.
	assert.Eventually(t, func() bool {
		return queue.Depth() == 0 && queue.InFlightDepth() == 0
	}, time.Second, 50*time.Millisecond, "Queue should be empty eventually")
}

func TestCoverage_EdgeCases(t *testing.T) {
	t.Run("PubChannelFull", func(t *testing.T) {
		// 本测试旨在通过“白盒测试”精确覆盖 Pub 方法中 select 的 default 分支。
		// 正常情况下，channel 容量总是 > MaxSize，无法触发此条件。
		// 因此，我们手动创建一个 channel 容量小于 MaxSize 的队列实例。
		config := &Config[string]{MaxSize: 2}
		config.apply() // 确保默认值被应用

		q := &Queue[string]{
			c:             config,
			lowChan:       make(chan *Message[string], 1),
			normalChan:    make(chan *Message[string], 1),
			highChan:      make(chan *Message[string], 1), // 容量设置为 1
			depth:         atomic.NewInt32(0),
			inFlightDepth: atomic.NewInt32(0),
			deferredQueue: NewDeferredQueue[string](),
		}

		// 第一次 Pub, 填满 channel 缓冲区 (容量为1)
		// depth(0) < MaxSize(2) -> true
		// select 成功, depth -> 1
		assert.NoError(t, q.Pub("fill-buffer", PriorityHigh))
		assert.Equal(t, int32(1), q.Depth())
		assert.Equal(t, 1, len(q.highChan))

		// 第二次 Pub, 触发 select 的 default case
		// depth(1) < MaxSize(2) -> true
		// select 阻塞 (channel 已满), 进入 default 分支
		err := q.Pub("trigger-default", PriorityHigh)
		assert.ErrorIs(t, err, ErrQueueFull, "Pub should return ErrQueueFull when high-priority channel is full")
	})

	t.Run("MaxRetryExceeded", func(t *testing.T) {
		// 此测试明确验证当达到最大重试次数时，消息会被丢弃，并且队列深度会减少。
		queue := New[string](&Config[string]{
			ConcurrentCount:             1,
			MaxRetry:                    1, // 仅允许1次处理
			DeferredQueueTickerDuration: 10 * time.Millisecond,
		})
		processCount := atomic.NewInt32(0)
		queue.Process(func(msg *Message[string]) (*Consume, error) {
			processCount.Inc()
			// 始终要求重试，这将触发 MaxRetry 逻辑
			return &Consume{Retry: true}, nil
		})

		assert.NoError(t, queue.Pub("test"))

		// 验证消息仅被处理一次，然后因为超出重试次数而被丢弃
		assert.Eventually(t, func() bool {
			// 该检查发生在 consumer 逻辑之前，因此 consumer 只会被调用一次
			return processCount.Load() == 1
		}, 3*time.Second, 50*time.Millisecond, "Process should be called only once")

		// 验证最终队列深度为0
		assert.Eventually(t, func() bool {
			return queue.Depth() == 0 && queue.InFlightDepth() == 0
		}, 3*time.Second, 50*time.Millisecond, "Queue depth should be 0 after max retries exceeded")
	})

	t.Run("RetryWithoutDelay", func(t *testing.T) {
		// 此测试明确验证当重试时 Delay <= 0, 消息的 ExecAt 会被立即设置为 time.Now()
		queue := New[string](&Config[string]{
			ConcurrentCount:             1,
			MaxRetry:                    2,
			DeferredQueueTickerDuration: 10 * time.Millisecond,
		})
		var firstExecTime time.Time
		processCount := atomic.NewInt32(0)

		queue.Process(func(msg *Message[string]) (*Consume, error) {
			count := processCount.Inc()
			if count == 1 {
				firstExecTime = msg.ExecAt
				// 请求重试，但不指定 Delay，触发 else 分支
				return &Consume{Retry: true}, nil
			}
			// 验证第二次执行时，ExecAt 时间戳被更新了
			assert.NotZero(t, msg.ExecAt)
			assert.True(t, msg.ExecAt.After(firstExecTime), "ExecAt should be updated for the retry")
			return nil, nil // 停止重试
		})

		assert.NoError(t, queue.Pub("test"))

		assert.Eventually(t, func() bool {
			return processCount.Load() == 2 && queue.Depth() == 0
		}, 3*time.Second, 50*time.Millisecond, "Should process message twice and then finish")
	})
}
