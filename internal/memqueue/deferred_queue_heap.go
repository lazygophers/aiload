package memqueue

import (
	"container/heap"
	"sync"
	"time"
)

// item 是延迟队列中存储的元素。
type item[T any] struct {
	*Message[T]
	// index 是该项在堆中的索引。
	// 这是 heap.Interface 实现所必需的。
	index int
}

// priorityQueue 实现了 heap.Interface，并按消息的执行时间（ExecAt）排序。
type priorityQueue[T any] []*item[T]

// Len 返回队列的长度。
func (pq priorityQueue[T]) Len() int { return len(pq) }

// Less 比较两个元素的优先级。
// 这里我们希望执行时间（ExecAt）最早的元素拥有最高的优先级，因此这是一个最小堆。
func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].ExecAt.Before(pq[j].ExecAt)
}

// Swap 交换队列中的两个元素。
func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push 向队列中添加一个元素。
// heap.Push 会调用此方法。
func (pq *priorityQueue[T]) Push(x any) {
	n := len(*pq)
	_item := x.(*item[T])
	_item.index = n
	*pq = append(*pq, _item)
}

// Pop 从队列中移除并返回优先级最高的元素（执行时间最早的）。
// heap.Pop 会调用此方法。
func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	_item := old[n-1]
	old[n-1] = nil   // 避免内存泄漏
	_item.index = -1 // 表示已移除
	*pq = old[0 : n-1]
	return _item
}

// DeferredQueueHeap 是一个基于最小堆的高效延迟队列。
// 它通过使用 `container/heap` 实现了对延迟消息的快速插入和提取。
type DeferredQueueHeap[T any] struct {
	pq       *priorityQueue[T]
	mux      sync.RWMutex
	itemPool sync.Pool // **性能优化**: 使用 sync.Pool 来复用 item 对象，减少 GC 压力。
}

// NewDeferredQueueHeap 创建一个新的基于堆的延迟队列实例。
func NewDeferredQueueHeap[T any]() *DeferredQueueHeap[T] {
	pq := make(priorityQueue[T], 0)
	// 初始化堆
	heap.Init(&pq)
	return &DeferredQueueHeap[T]{
		pq: &pq,
		itemPool: sync.Pool{
			New: func() any {
				return new(item[T])
			},
		},
	}
}

// Push 将消息推入延迟队列。
// 时间复杂度为 O(log n)。
// **性能优化**: 从 sync.Pool 获取 item 对象，避免频繁的内存分配。
func (p *DeferredQueueHeap[T]) Push(msg *Message[T]) {
	p.mux.Lock()
	defer p.mux.Unlock()

	// 从池中获取一个 item
	poolItem := p.itemPool.Get().(*item[T])
	poolItem.Message = msg
	// index 会在 heap.Push 中设置
	heap.Push(p.pq, poolItem)
}

// Pop 弹出所有已到期的消息。
// 它会检查堆顶的元素，如果其执行时间已到，则将其弹出，
// 并继续检查下一个，直到所有到期的消息都被处理。
// **性能优化**: 将弹出的 item 放回 sync.Pool 以便复用。
func (p *DeferredQueueHeap[T]) Pop() []*Message[T] {
	p.mux.Lock()
	defer p.mux.Unlock()

	var expiredMessages []*Message[T]
	for p.pq.Len() > 0 {
		// 查看堆顶元素而不移除它
		top := (*p.pq)[0]
		if top.ExecAt.After(time.Now()) {
			// 如果堆顶的元素还未到期，那么其他所有元素也肯定未到期
			break
		}

		// 移除已到期的元素
		poppedItem := heap.Pop(p.pq).(*item[T])
		expiredMessages = append(expiredMessages, poppedItem.Message)

		// **性能优化**: 重置并放回池中
		poppedItem.Message = nil
		poppedItem.index = 0
		p.itemPool.Put(poppedItem)
	}

	return expiredMessages
}

// Depth 返回当前延迟队列中的消息数量。
func (p *DeferredQueueHeap[T]) Depth() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.pq.Len()
}
