package slice_heap

import (
	"container/heap"
	"testing"
	"fmt"
)

func IntLessThan(a, b interface{}) bool {
	return a.(int) < b.(int)
}

func StringLessThan(a, b interface{}) bool {
	return a.(string) < b.(string)
}

func TestIntSliceHeap(t *testing.T) {
	sh := NewSliceHeap(IntLessThan)
	heap.Init(sh)
	heap.Push(sh, 5)
	heap.Push(sh, 2)
	heap.Push(sh, 3)

	vals := []int{2, 3, 5}
	for _, val := range vals {
		i := heap.Pop(sh).(int)
		if i != val {
			t.Error(fmt.Sprintf("popped %d instead of %d", i, val))
		}
	}
}

func TestStringSliceHeap(t *testing.T) {
	sh := NewSliceHeap(StringLessThan)
	heap.Push(sh, "dog")
	heap.Push(sh, "cat")
	heap.Push(sh, "mouse")

	vals := []string{"cat", "dog", "mouse"}
	for _, val := range vals {
		s := heap.Pop(sh).(string)
		if s != val {
			t.Error(fmt.Sprintf("popped %s instead of %s", s, val))
		}
	}
}