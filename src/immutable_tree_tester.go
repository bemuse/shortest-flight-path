package main

import (
	it "immutable_tree"
)

func intValueTester() {
	tree := it.NewTree()
	var oldTree *it.Tree
	vs := [...]int{5, 8, 1, 7, 2, 6, 12, 9, 4, 3, 10, 0, 99}
	for i, v := range vs {
		if len(vs)/2 == i {
			oldTree = tree
		}
		iv := it.IntValue(v)
		tree = tree.AddValue(&iv)
	}

	oldTree.Display()
	tree.Display()
}

func stringValueTester() {
	tree := it.NewTree()
	var oldTree *it.Tree
	vs := [...]string{"cat", "mouse", "dog", "elephant", "anaconda", "zebra", "fox", "yak", "goose", "tiger", "lion", "bobcat", "chicken"}
	for i, v := range vs {
		if len(vs)/2 == i {
			oldTree = tree
		}
		sv := it.StringValue(v)
		tree = tree.AddValue(&sv)
	}
	oldTree.Display()
	tree.Display()
}

func main() {
	intValueTester()
	stringValueTester()
}
