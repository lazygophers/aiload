package memqueue

import (
	"fmt"
	"sync"
	"testing"
)

// BenchmarkPub 旨在测量向队列发布消息的性能，
// 特别是在高并发发布场景下的表现。
func BenchmarkPub(b *testing.B) {
	q := New(&Config[int]{
		// 配置足够多的消费者，确保队列能被迅速清空，
		// 从而让基准测试聚焦于发布者（Pub）的性能。
		ConcurrentCount: 64,
		MaxSize:         b.N + 1, // 保证队列在测试期间不会被填满
	})

	// 处理函数尽可能简化，以减小其对测试结果的干扰。
	q.Process(func(msg *Message[int]) (*Consume, error) {
		return nil, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	// b.RunParallel 会启动与 GOMAXPROCS 相等数量的 goroutine 并发执行。
	// 每个 goroutine 会在循环中持续调用 q.Pub。
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if err := q.Pub(i); err != nil {
				// 在足够大的 MaxSize 配置下，理论上不应发生此错误。
				b.Errorf("Pub failed: %v", err)
			}
			i++
		}
	})
}

// benchmarkProcess 是一个辅助函数，用于测试从发布消息到处理完成的
// 端到端吞吐量。
func benchmarkProcess(b *testing.B, concurrentConsumers int) {
	q := New(&Config[int]{
		ConcurrentCount: concurrentConsumers,
		MaxSize:         b.N + 1, // 保证队列有足够的容量
	})

	var wg sync.WaitGroup
	// 我们将等待 b.N 条消息被成功处理。
	wg.Add(b.N)

	// 处理器通过调用 wg.Done() 来发出处理完成的信号。
	q.Process(func(msg *Message[int]) (*Consume, error) {
		wg.Done()
		return nil, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	// 这是基准测试的主循环，Go 的 testing 框架会自动调整 b.N 的值。
	for i := 0; i < b.N; i++ {
		if err := q.Pub(i); err != nil {
			b.Fatalf("Pub failed: %v", err)
		}
	}

	// 等待所有消息被处理完成。
	// 这里消耗的时间是基准测试的一部分，因为它直接反映了处理速度。
	wg.Wait()
}

// BenchmarkProcess 用于测试在不同并发消费者数量下，
// 队列的整体吞吐性能。
func BenchmarkProcess(b *testing.B) {
	scenarios := []int{1, 2, 4, 8, 16, 32}
	for _, consumers := range scenarios {
		// 为每个并发级别创建一个子基准测试。
		b.Run(fmt.Sprintf("Consumers-%d", consumers), func(b *testing.B) {
			benchmarkProcess(b, consumers)
		})
	}
}
