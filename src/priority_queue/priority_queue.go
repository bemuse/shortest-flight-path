package priority_queue

import (
	"container/heap"
	sheap "slice_heap"
)

type PriorityQueue struct {
	myHeap *sheap.SliceHeap
}

func NewPriorityQueue(f func(a, b interface{}) bool) *PriorityQueue {
	sh := sheap.NewSliceHeap(f)
	heap.Init(sh)
	return &PriorityQueue{sh}
}

func (pq *PriorityQueue) Push(item interface{}) {
	heap.Push(pq.myHeap, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	return heap.Pop(pq.myHeap)
}

func (pq *PriorityQueue) Empty() bool {
	return 0 == pq.myHeap.Len()
}
