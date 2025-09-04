package memqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDeferredQueueHeap_Push 测试向堆中推送消息的逻辑。
func TestDeferredQueueHeap_Push(t *testing.T) {
	// 初始化一个空的延迟队列堆
	q := NewDeferredQueueHeap[string]()

	// 创建两个消息，并设置不同的执行时间
	msg1 := &Message[string]{Data: "later", ExecAt: time.Now().Add(2 * time.Second)}
	msg2 := &Message[string]{Data: "sooner", ExecAt: time.Now().Add(1 * time.Second)}

	// 推送消息到堆中
	q.Push(msg1)
	q.Push(msg2)

	// 验证堆的大小是否正确
	assert.Equal(t, 2, q.Depth(), "堆中应有2个消息")

	// 验证堆顶的元素是否是执行时间最早的那个。
	// 由于 Pop 只会弹出*已过期*的消息，而我们刚推入的消息都未到期，
	// 所以这里 Pop 应该返回一个空切片。
	popped := q.Pop()
	assert.Empty(t, popped, "对于未到期的消息，Pop 应该返回空")

	// 让时间流逝，使其中一个消息过期
	time.Sleep(1100 * time.Millisecond) // 等待 1.1 秒，确保 msg2 过期

	// 再次尝试弹出
	popped = q.Pop()
	assert.Equal(t, 1, len(popped), "在等待后，应该弹出1个已过期的消息")
	assert.Equal(t, "sooner", popped[0].Data, "弹出的应该是执行时间较早的消息")
	assert.Equal(t, 1, q.Depth(), "弹出一个后，堆中应只剩一个消息")
}

// TestDeferredQueueHeap_Pop 测试从堆中弹出到期消息的逻辑。
func TestDeferredQueueHeap_Pop(t *testing.T) {
	q := NewDeferredQueueHeap[string]()

	// 创建一个已过期的消息和一个未到期的消息
	expiredMsg := &Message[string]{Data: "expired", ExecAt: time.Now().Add(-1 * time.Second)}
	futureMsg := &Message[string]{Data: "future", ExecAt: time.Now().Add(1 * time.Hour)}

	// 推送消息
	q.Push(expiredMsg)
	q.Push(futureMsg)

	// 弹出消息
	popped := q.Pop()

	// 验证是否只弹出了已过期的消息
	assert.Equal(t, 1, len(popped), "应仅弹出1个已过期的消息")
	assert.Equal(t, expiredMsg.Data, popped[0].Data, "弹出的消息内容应正确")
	assert.Equal(t, 1, q.Depth(), "堆中应还剩下1个未到期的消息")
}

// TestDeferredQueueHeap_Pop_Empty 测试从空堆中弹出的情况。
func TestDeferredQueueHeap_Pop_Empty(t *testing.T) {
	q := NewDeferredQueueHeap[string]()
	popped := q.Pop()
	assert.Empty(t, popped, "从空堆中应弹出空切片")
}

// TestDeferredQueueHeap_Pop_AllExpired 测试所有消息都已到期的情况。
func TestDeferredQueueHeap_Pop_AllExpired(t *testing.T) {
	q := NewDeferredQueueHeap[string]()
	msg1 := &Message[string]{Data: "expired1", ExecAt: time.Now().Add(-2 * time.Second)}
	msg2 := &Message[string]{Data: "expired2", ExecAt: time.Now().Add(-1 * time.Second)}

	q.Push(msg1)
	q.Push(msg2)

	popped := q.Pop()
	assert.Equal(t, 2, len(popped), "应弹出所有已过期的消息")
	assert.Equal(t, 0, q.Depth(), "弹出所有消息后，堆应为空")
}

// BenchmarkDeferredQueueHeap_Push 是对基于堆的延迟队列 Push 操作的基准测试。
func BenchmarkDeferredQueueHeap_Push(b *testing.B) {
	q := NewDeferredQueueHeap[int]()
	msg := &Message[int]{Data: 1, ExecAt: time.Now()}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(msg)
	}
}

// BenchmarkDeferredQueueHeap_Pop 是对基于堆的延迟队列 Pop 操作的基准测试。
func BenchmarkDeferredQueueHeap_Pop(b *testing.B) {
	q := NewDeferredQueueHeap[int]()
	// 预填充大量数据
	for i := 0; i < b.N; i++ {
		q.Push(&Message[int]{Data: i, ExecAt: time.Now().Add(time.Duration(i) * time.Nanosecond)})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
