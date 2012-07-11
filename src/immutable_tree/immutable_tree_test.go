package immutable_tree

import (
	"testing"
)

func TestInts(t *testing.T) {
	tree := NewTree()
	var oldTree *Tree
	vs := [...]int{5, 8, 1, 7, 2, 6, 12, 9, 4, 3, 10, 0, 99}
	for i, v := range vs {
		if len(vs)/2 == i {
			oldTree = tree
		}
		iv := IntValue(v)
		tree = tree.AddValue(&iv)
	}

	for i, v := range vs {
		iv := IntValue(v)
		if !tree.HasValue(&iv) {
			t.Error("missing value")
		}
		if i < len(vs)/2 {
			if !oldTree.HasValue(&iv) {
				t.Error("missing value")
			}
		} else {
			if oldTree.HasValue(&iv) {
				t.Error("extra value")
			}
		}
	}
}

func TestStrings(t *testing.T) {
	tree := NewTree()
	var oldTree *Tree
	vs := [...]string{"cat", "mouse", "dog", "elephant", "anaconda", "zebra", "fox", "yak", "goose", "tiger", "lion", "bobcat", "chicken"}
	for i, v := range vs {
		if len(vs)/2 == i {
			oldTree = tree
		}
		iv := StringValue(v)
		tree = tree.AddValue(&iv)
	}

	for i, v := range vs {
		iv := StringValue(v)
		if !tree.HasValue(&iv) {
			t.Error("missing value")
		}
		if i < len(vs)/2 {
			if !oldTree.HasValue(&iv) {
				t.Error("missing value")
			}
		} else {
			if oldTree.HasValue(&iv) {
				t.Error("extra value")
			}
		}
	}
}
