package graph

type Vertex struct {
	ID int64
	edges []Edge
}

func (v Vertex) AddEdge(e Edge) {
	v.edges = append(v.edges, e)
}

