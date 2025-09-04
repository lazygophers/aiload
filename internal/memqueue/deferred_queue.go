package memqueue

// DeferredQueue 是对高效的、基于堆的延迟队列的封装。
// 它现在内部使用 DeferredQueueHeap 来管理所有延迟消息，
// 从而获得了 O(log n) 的插入和 O(1) 的峰值检查性能。
type DeferredQueue[T any] struct {
	*DeferredQueueHeap[T]
}

// NewDeferredQueue 创建一个新的延迟队列实例。
// 它通过初始化内部的 DeferredQueueHeap 来提供一个高性能的延迟消息解决方案。
func NewDeferredQueue[T any]() *DeferredQueue[T] {
	return &DeferredQueue[T]{
		DeferredQueueHeap: NewDeferredQueueHeap[T](),
	}
}
