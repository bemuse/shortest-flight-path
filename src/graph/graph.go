package graph

import (
	"container/heap"
	"fmt"
	sheap "slice_heap"
)

const (
	BUILD_FLAG       = 1
	TRAVERSE_FLAG    = 2
	PROGRESSIVE_FLAG = 4
	DEBUG            = 0 // | PROGRESSIVE_FLAG // | TRAVERSE_FLAG // | BUILD_FLAG 
)

// VisitedList

type VisitedList struct {
	node *Node
	next *VisitedList
}

func (l *VisitedList) HasVisited(n *Node) bool {
	for ; l != nil; l = l.next {
		if l.node == n {
			return true
		}
	}

	return false
}

func (l *VisitedList) AddNode(n *Node) *VisitedList {
	return &VisitedList{n, l}
}

func (l *VisitedList) MakeSlice() (result []*Node) {
	result = make([]*Node, 0)

	// copy values from list into slice
	for l2 := l; l2 != nil; l2 = l2.next {
		result = append(result, l2.node)
	}

	// swap pairs to reverse slice
	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] =
			result[len(result)-1-i], result[i]
	}

	return
}

func (l *VisitedList) Print() {
	if l.next != nil {
		l.next.Print()
	}
	fmt.Print(l.node.Record, "  ")
}

// State

type PrivateTraverseState interface {
	// LessThan(other *PrivateTraverseState) bool
	TraverseStateHelper(v *Vertex) (nextState PrivateTraverseState, ok bool)
}

type PublicTraverseState struct {
	totalCost    float64
	node         *Node
	visited      *VisitedList
	privateState PrivateTraverseState
}

func PublicStateLessThan(d1, d2 interface{}) bool {
	pubState1 := d1.(*PublicTraverseState)
	pubState2 := d2.(*PublicTraverseState)
	return pubState1.totalCost < pubState2.totalCost
}

// Node and NodeDescription and graph

type NodeRecord interface {
	String() string
}

type Node struct {
	Record   NodeRecord
	vertices []*Vertex
}

type Vertex struct {
	From, To *Node
	Cost     float64
}

// Graph

type Graph struct {
	nodes    []*Node
	vertices []*Vertex
}

func NewGraph() *Graph {
	return &Graph{make([]*Node, 0), make([]*Vertex, 0)}
}

func (g *Graph) NewNode(record NodeRecord) *Node {
	n := &Node{record, make([]*Vertex, 0)}
	g.nodes = append(g.nodes, n)
	return n
}

func (g *Graph) ConnectUni(from, to *Node, cost float64) {
	v := &Vertex{from, to, cost}
	g.vertices = append(g.vertices, v)
	from.vertices = append(from.vertices, v)
	if DEBUG&BUILD_FLAG != 0 {
		fmt.Printf("Connecting %s to %s at cost %f.\n", from.Record.String(), to.Record.String(), cost)
	}
}

func (g *Graph) ConnectBi(n1, n2 *Node, cost float64) {
	g.ConnectUni(n1, n2, cost)
	g.ConnectUni(n2, n1, cost)
}

func (g *Graph) Traverse(privateState PrivateTraverseState, from, to *Node) (path []*Node, totalCost float64, ok bool) {
	seen := make(map[*Node]bool)

	state := &PublicTraverseState{0.0, from, &VisitedList{from, nil}, privateState}

	sh := sheap.NewSliceHeap(PublicStateLessThan)
	heap.Init(sh)
	heap.Push(sh, state)

	for !sh.IsEmpty() {
		state = heap.Pop(sh).(*PublicTraverseState)
		if _, found := seen[state.node]; found {
			continue
		}
		seen[state.node] = true

		if DEBUG&PROGRESSIVE_FLAG != 0 {
			fmt.Printf("%f : %s\n", state.totalCost, state.node.Record)
		}

		if state.node == to {
			return state.visited.MakeSlice(), state.totalCost, true
		}

		for _, vertex := range state.node.vertices {
			if DEBUG&TRAVERSE_FLAG != 0 {
				fmt.Printf("Considering %s to %s ... ", state.node.Record, vertex.To.Record)
			}
			totalCost := state.totalCost + vertex.Cost
			nextNode := vertex.To
			if state.visited.HasVisited(nextNode) {
				if DEBUG&TRAVERSE_FLAG != 0 {
					fmt.Println("already been there")
				}
				continue
			}
			nextPrivateState, ok := (state.privateState).TraverseStateHelper(vertex)
			if !ok {
				if DEBUG&TRAVERSE_FLAG != 0 {
					fmt.Println("rejected by calling code")
				}
				continue
			}
			nextPublicState := &PublicTraverseState{totalCost, nextNode, &VisitedList{nextNode, state.visited}, nextPrivateState}
			heap.Push(sh, nextPublicState)

			if DEBUG&TRAVERSE_FLAG != 0 {
				fmt.Printf("good w/ total cost of %f\n", totalCost)
			}
		}
	}

	return nil, 0.0, false
}

func (g *Graph) Display() {
	for _, n := range g.nodes {
		fmt.Printf("%s:\n", n.Record)
		for _, v := range n.vertices {
			fmt.Printf("    %s @ %f\n", v.To.Record, v.Cost)
			if v.From != n {
				panic("non-matching node/vertex")
			}
		}
	}
}
