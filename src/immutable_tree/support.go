package immutable_tree

import (
	"fmt"
)

type IntValue int

func (v1 *IntValue) CompareTo(c2 Comparable) int {
	v2, ok := c2.(*IntValue)
	if !ok {
		panic("IntValue CompareTo passed non-intValue")
	}
	return int(*v1) - int(*v2)
}

func (v *IntValue) String() string {
	return fmt.Sprintf("%d", int(*v))
}

type StringValue string

func (v1 *StringValue) CompareTo(c2 Comparable) int {
	v2, ok := c2.(*StringValue)
	if !ok {
		panic("StringValue CompareTo passed non-stringValue")
	}
	if string(*v1) < string(*v2) {
		return -1
	} else if string(*v1) > string(*v2) {
		return 1
	}
	return 0
}

func (v *StringValue) String() string {
	return string(*v)
}
