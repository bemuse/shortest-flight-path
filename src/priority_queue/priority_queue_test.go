package priority_queue

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


func TestIntPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue(IntLessThan)
	pq.Push(5)
	pq.Push(2)
	pq.Push(3)

	vals := []int{2, 3, 5}
	for _, val := range vals {
		i := pq.Pop().(int)
		if i != val {
			t.Error(fmt.Sprintf("popped %d instead of %d", i, val))
		}
	}
}

func TestStringPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue(StringLessThan)
	pq.Push("dog")
	pq.Push("cat")
	pq.Push("mouse")

	vals := []string{"cat", "dog", "mouse"}
	for _, val := range vals {
		s := pq.Pop().(string)
		if s != val {
			t.Error(fmt.Sprintf("popped %s instead of %s", s, val))
		}
	}
}