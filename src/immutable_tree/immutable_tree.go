package immutable_tree

import (
	"fmt"
)

type Comparable interface {
	CompareTo(other Comparable) int
	String() string
}

type node struct {
	value       Comparable
	left, right *node
}

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{nil}
}

func (thisNode *node) addHelper(newNode *node) *node {
	// compare := (*newNode.value).CompareTo(thisNode.value)
	compare := newNode.value.CompareTo(thisNode.value)
	if compare < 0 {
		if thisNode.left == nil {
			return &node{thisNode.value, newNode, thisNode.right}
		} else {
			newLeft := thisNode.left.addHelper(newNode)
			if newLeft == thisNode.left {
				return thisNode
			} else {
				return &node{thisNode.value, newLeft, thisNode.right}
			}
		}
	} else if compare > 0 {
		if thisNode.right == nil {
			return &node{thisNode.value, thisNode.left, newNode}
		} else {
			newRight := thisNode.right.addHelper(newNode)
			if newRight == thisNode.right {
				return thisNode
			} else {
				return &node{thisNode.value, thisNode.left, newRight}
			}
		}
	}
	return thisNode
}

func (t *Tree) AddValue(value Comparable) (result *Tree) {
	newNode := &node{value, nil, nil}
	if t.root == nil {
		return &Tree{newNode}
	} else {
		newRoot := t.root.addHelper(newNode)
		if newRoot == t.root {
			return t
		} else {
			return &Tree{newRoot}
		}
	}
	return nil // should never be executed
}

func (t *Tree) HasValue(value Comparable) bool {
	cursor := t.root
	for cursor != nil {
		compare := value.CompareTo(cursor.value)
		if compare == 0 {
			return true
		} else if compare < 0 {
			cursor = cursor.left
		} else {
			cursor = cursor.right
		}
	}
	return false
}

func (n *node) displayHelper() {
	if n.left != nil {
		n.left.displayHelper()
	}
	fmt.Printf("%s (%p) ", n.value.String(), n)
	if n.right != nil {
		n.right.displayHelper()
	}
}

func (t *Tree) Display() {
	if t.root == nil {
		fmt.Println("empty")
	} else {
		t.root.displayHelper()
		fmt.Println()
	}
}
