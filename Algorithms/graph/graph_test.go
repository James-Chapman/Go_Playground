package graph

 import (
 	"fmt"
 	"testing"
 	"math/rand"
// 	//"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	pGraph := NewGraph()
	var i int64
	
	for i = 0; i < 1000; i++ {
		v := Vertex{i,nil}
		pGraph.AddVertex(v)
	}

	for i = 0; i < 1000; i++ {
		e := Edge{rand.Intn(1000), rand.Intn(1000), rand.Intn(110)}
		pGraph.AddEdge(e)
	}

	fmt.Println(*pGraph)
}

// func TestBFS(t *testing.T) {
// 	var v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11,v12 Vertex
// 	var e1,e2,e3,e4,e5,e6,e7,e8,e9,e10,e11,e12,e13,e14,e15,e16,e17,e18 Edge

// 	e1 = edge{1, 2, 1}
// 	e2 = edge{2, 3, 1}
// 	e3 = edge{3, 4, 1}
// 	e4 = edge{4, 5, 1}
// 	e5 = edge{5, 6, 1}
// 	e6 = edge{6, 7, 1}
// 	e7 = edge{7, 8, 1}
// 	e8 = edge{8, 9, 1}
// 	e9 = edge{9, 10, 1}
// 	e10 = edge{10, 11, 1}
// 	e11 = edge{11, 12, 1}
// 	e12 = edge{12, 6, 3}
// 	e13 = edge{12, 8, 5}
// 	e14 = edge{12, 10, 4}
// 	e15 = edge{3, 6, 3}
// 	e16 = edge{4, 8, 4}
// 	e17 = edge{5, 10, 5}
// 	e18 = edge{6, 12, 6}

// 	v1 = vertex{1, []Edge{e1}}
// 	v2 = vertex{2, []Edge{e2}}
// 	v3 = vertex{3, []Edge{e3,e15}}
// 	v4 = vertex{4, []Edge{e4,e16}}
// 	v5 = vertex{5, []Edge{e5,e17}}
// 	v6 = vertex{6, []Edge{e6,e18}}
// 	v7 = vertex{7, []Edge{e7}}
// 	v8 = vertex{8, []Edge{e8}}
// 	v9 = vertex{9, []Edge{e9}}
// 	v10 = vertex{10, []Edge{e10}}
// 	v11 = vertex{11, []Edge{e11}}
// 	v12 = vertex{12, []Edge{e12,e13,e14}}

// 	var graph Graph
// 	graph.vertices = []Vertex{v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11,v12}
// 	graph.edges = []Edge{e1,e2,e3,e4,e5,e6,e7,e8,e9,e10,e11,e12,e13,e14,e15,e16,e17,e18}

// 	for _, v := range graph.vertices {
// 		fmt.Println(v)
// 	}

// 	for _, v := range graph.edges {
// 		fmt.Println(v)
// 	}
// }