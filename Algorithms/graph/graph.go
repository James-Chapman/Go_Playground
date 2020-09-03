package graph

import (
	//"fmt"
	//"container/heap"
	//"math"
)


type Graph struct {
	vertices []Vertex
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) AddVertex(v Vertex) {
	g.vertices = append(g.vertices, v)
}

func (g *Graph) AddEdge(e Edge) {
	g.edges = append(g.edges, e)
}






// // A PriorityQueue implements heap.Interface and holds Items.
// type PriorityQueue []*Edge

// func (pq PriorityQueue) Len() int64 { 
// 	return len(pq) 
// }

// func (pq PriorityQueue) Less(i, j int) bool {
// 	// We want Pop to give us the lowest cost
// 	return pq[i].c < pq[j].c
// }

// func (pq PriorityQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// }

// func (pq *PriorityQueue) Push(x interface{}) {
// 	edge := x.(*Edge)
// 	*pq = append(*pq, edge)
// }

// func (pq *PriorityQueue) Pop() interface{} {
// 	old := *pq
// 	n := len(old)
// 	edge := old[n-1]
// 	old[n-1] = nil  // avoid memory leak
// 	*pq = old[0 : n-1]
// 	return edge
// }


