package slice_heap

type LessThanFunc func(a, b interface{}) bool

type SliceHeap struct {
	collection []interface{}
	lessThan   LessThanFunc
}

func NewSliceHeap(f LessThanFunc) *SliceHeap {
	return &SliceHeap{make([]interface{}, 0), f}
}

func (sh *SliceHeap) Push(item interface{}) {
	sh.collection = append(sh.collection, item)
}

func (sh *SliceHeap) Pop() (result interface{}) {
	l := len(sh.collection)
	result = sh.collection[l-1]
	sh.collection = sh.collection[0 : l-1]
	return
}

func (sh *SliceHeap) Len() int {
	return len(sh.collection)
}

func (sh *SliceHeap) Less(i, j int) bool {
	return sh.lessThan(sh.collection[i], sh.collection[j])
}

func (sh *SliceHeap) Swap(i, j int) {
	sh.collection[i], sh.collection[j] = sh.collection[j], sh.collection[i]
}

func (sh *SliceHeap) IsEmpty() bool {
	return 0 == len(sh.collection)
}
